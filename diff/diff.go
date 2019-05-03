package diff

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ttakezawa/pgconverger/ast"
	"github.com/ttakezawa/pgconverger/lexer"
	"github.com/ttakezawa/pgconverger/parser"
)

type fileReader interface {
	io.Reader
	Name() string
}

type Diff struct {
	source       fileReader
	sourceDDL    *ast.DataDefinition
	sourceErrors []error

	desired       fileReader
	desiredDDL    *ast.DataDefinition
	desiredErrors []error

	sourceSchemas map[string]struct{}
	sourceTables  map[string]map[string]*Table

	desiredSchemas map[string]struct{}
	desiredTables  map[string]map[string]*Table

	stringBuilder *strings.Builder
}

func Process(source fileReader, desired fileReader) (string, error) {
	df := &Diff{
		source:        source,
		desired:       desired,
		stringBuilder: &strings.Builder{},
	}
	df.sourceErrors, df.sourceDDL = df.parseOneSide(df.source)
	df.desiredErrors, df.desiredDDL = df.parseOneSide(df.desired)
	if df.ErrorOrNil() != nil {
		return "", df.ErrorOrNil()
	}

	return df.generatePatch(), nil
}

func (df *Diff) WriteString(s string) {
	_, _ = df.stringBuilder.WriteString(s)
}

func (*Diff) parseOneSide(reader fileReader) (errs []error, ddl *ast.DataDefinition) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		errs = append(errs, err)
		return
	}
	p := parser.New(lexer.Lex(reader.Name(), string(input)))
	ddl = p.ParseDataDefinition()
	errs = append(errs, p.Errors()...)
	return
}

func (df *Diff) ErrorOrNil() error {
	if df.sourceErrors != nil {
		return &Error{df}
	}
	if df.desiredErrors != nil {
		return &Error{df}
	}
	return nil
}

type (
	Table struct {
		CreateTableStatement *ast.CreateTableStatement
		Schema               string
		Name                 string
		Columns              map[string]*Column
		Indexes              map[string]*Index
	}

	Column struct {
		Name     string
		DataType string
		NotNull  bool
		Default  string
	}

	Index struct {
		Name    string
		Primary bool
		Unique  bool
		Columns []string
	}
)

func (df *Diff) generatePatch() string {
	df.sourceSchemas, df.sourceTables = processDDL(df.sourceDDL)
	df.desiredSchemas, df.desiredTables = processDDL(df.desiredDDL)

	for searchPath, tables := range df.sourceTables {
		for _, table := range tables {
			desiredTable := findTable(df.desiredTables, searchPath, table.Name)
			if desiredTable != nil {
				df.diffTable(table, desiredTable)
			} else {
				df.dropTable(table)
			}
		}
	}

	for searchPath, tables := range df.desiredTables {
		for _, table := range tables {
			desiredTable := findTable(df.sourceTables, searchPath, table.Name)
			if desiredTable != nil {
				// none
			} else {
				df.createTable(searchPath, table)
			}
		}
	}

	return df.stringBuilder.String()
}

func (df *Diff) createTable(searchPath string, table *Table) {
	if table.CreateTableStatement.SchemaName == nil {
		table.CreateTableStatement.SetSchema(searchPath)
	}
	table.CreateTableStatement.WriteStringTo(df.stringBuilder)
}

func (df *Diff) dropTable(table *Table) {
	df.WriteString(fmt.Sprintf("DROP TABLE \"%s\".\"%s\";\n", table.Schema, table.Name))
}

// generate DDL for a table which exists in both.
func (df *Diff) diffTable(sourceTable, desiredTable *Table) {
	for _, sourceColumn := range sourceTable.Columns {
		desiredColumn, ok := desiredTable.Columns[sourceColumn.Name] //
		if ok {
			df.alterColumn(sourceTable, sourceColumn, desiredColumn)
		} else {
			df.dropColumn(sourceTable, sourceColumn)
		}
	}

	for _, desiredColumn := range desiredTable.Columns {
		_, ok := sourceTable.Columns[desiredColumn.Name] //
		if !ok {
			df.addColumn(sourceTable, desiredColumn)
		}
	}
}

func (df *Diff) addColumn(table *Table, column *Column) {
	df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ADD COLUMN \"%s\" \"%s\";\n",
		table.Schema,
		table.Name,
		column.Name,
		column.DataType,
	))
}

func (df *Diff) dropColumn(table *Table, column *Column) {
	df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" DROP COLUMN \"%s\";\n",
		table.Schema,
		table.Name,
		column.Name,
	))
}

func (df *Diff) alterColumn(table *Table, sourceColumn *Column, desiredColumn *Column) {
	if sourceColumn.DataType != desiredColumn.DataType {
		df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ALTER COLUMN \"%s\" TYPE \"%s\";\n",
			table.Schema,
			table.Name,
			desiredColumn.Name,
			desiredColumn.DataType,
		))
	}

	if sourceColumn.NotNull != desiredColumn.NotNull {
		if desiredColumn.NotNull {
			// SET not null
			df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ALTER COLUMN \"%s\" SET NOT NULL;\n",
				table.Schema,
				table.Name,
				desiredColumn.Name,
			))
		} else {
			// Drop not null
			df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ALTER COLUMN \"%s\" DROP NOT NULL;\n",
				table.Schema,
				table.Name,
				desiredColumn.Name,
			))
		}
	}

	if sourceColumn.Default != desiredColumn.Default {
		if desiredColumn.Default != "" {
			// SET DEFAULT
			df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ALTER COLUMN \"%s\" SET DEFAULT %s;\n",
				table.Schema,
				table.Name,
				desiredColumn.Name,
				desiredColumn.Default,
			))
		} else {
			// DROP DEFAULT
			df.WriteString(fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" ALTER COLUMN \"%s\" DROP DEFAULT;\n",
				table.Schema,
				table.Name,
				desiredColumn.Name,
			))
		}
	}
}

func findTable(tables map[string]map[string]*Table, schema string, tableName string) *Table {
	return tables[schema][tableName]
}

// processDDL converts to schema and table mappings
func processDDL(ddl *ast.DataDefinition) (schemas map[string]struct{}, tables map[string]map[string]*Table) {
	schemas = make(map[string]struct{})
	tables = make(map[string]map[string]*Table)
	searchPath := "public"

	for _, statement := range ddl.StatementList {
		switch stmt := statement.(type) {
		case *ast.SetStatement:
			if stmt.Name.Value == "search_path" {
				ident, ok := stmt.Values[0].(*ast.Identifier)
				if ok {
					searchPath = ident.Value
				}
			}
		case *ast.CreateSchemaStatement:
			schemas[stmt.Name.Value] = struct{}{}
		case *ast.CreateTableStatement:
			tbls, ok := tables[searchPath]
			if !ok {
				tbls = make(map[string]*Table)
				tables[searchPath] = tbls
			}
			columns := make(map[string]*Column)
			for _, columnDefinition := range stmt.ColumnDefinitionList {
				col := columnFromAst(columnDefinition)
				columns[col.Name] = col
			}
			tbls[stmt.TableName.Value] = &Table{
				CreateTableStatement: stmt,
				Schema:               searchPath,
				Name:                 stmt.TableName.Value,
				Columns:              columns,
			}
		case *ast.CreateIndexStatement:
			// TODO
			// 	tbls := df.sourceTables[df.searchPath]
			// 	t := tbls[stmt.TableName.Value]
			// 	t.CreateIndexes = append(t.CreateIndexes, stmt)
		default:
			log.Printf("skip statement: %v", stmt)
		}
	}
	return
}

func columnFromAst(columnDefinition *ast.ColumnDefinition) *Column {
	column := &Column{
		Name:     columnDefinition.Name.Value,
		DataType: ast.FormatNode(columnDefinition.Type),
	}
	for _, constraint := range columnDefinition.ConstraintList {
		switch v := constraint.(type) {
		case *ast.ColumnConstraintNotNull:
			column.NotNull = true
		case *ast.ColumnConstraintDefault:
			var buf bytes.Buffer
			v.Expr.WriteStringTo(&buf)
			column.Default = buf.String()
		}
	}
	return column
}
