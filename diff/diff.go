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

	sourceTables  Tables
	desiredTables Tables

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
		string
		Identifier string
		Columns    map[string]*Column
		Indexes    map[string]*Index
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
	df.sourceTables = processDDL(df.sourceDDL)
	df.desiredTables = processDDL(df.desiredDDL)

	for identifier, table := range df.sourceTables {
		desiredTable := findTable(df.desiredTables, identifier)
		if desiredTable != nil {
			df.diffTable(table, desiredTable)
		} else {
			df.dropTable(table)
		}
	}

	for identifier, table := range df.desiredTables {
		desiredTable := findTable(df.sourceTables, identifier)
		if desiredTable != nil {
			// none
		} else {
			df.createTable(table)
		}
	}

	return df.stringBuilder.String()
}

func (df *Diff) createTable(table *Table) {
	table.CreateTableStatement.WriteStringTo(df.stringBuilder)
}

func (df *Diff) dropTable(table *Table) {
	df.WriteString(fmt.Sprintf("DROP TABLE %s;\n", table.Identifier))
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
	df.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" \"%s\";\n",
		table.Identifier,
		column.Name,
		column.DataType,
	))
}

func (df *Diff) dropColumn(table *Table, column *Column) {
	df.WriteString(fmt.Sprintf("ALTER TABLE %s DROP COLUMN \"%s\";\n",
		table.Identifier,
		column.Name,
	))
}

func (df *Diff) alterColumn(table *Table, sourceColumn *Column, desiredColumn *Column) {
	if sourceColumn.DataType != desiredColumn.DataType {
		df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" TYPE \"%s\";\n",
			table.Identifier,
			desiredColumn.Name,
			desiredColumn.DataType,
		))
	}

	if sourceColumn.NotNull != desiredColumn.NotNull {
		if desiredColumn.NotNull {
			// SET not null
			df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" SET NOT NULL;\n",
				table.Identifier,
				desiredColumn.Name,
			))
		} else {
			// Drop not null
			df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" DROP NOT NULL;\n",
				table.Identifier,
				desiredColumn.Name,
			))
		}
	}

	if sourceColumn.Default != desiredColumn.Default {
		if desiredColumn.Default != "" {
			// SET DEFAULT
			df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" SET DEFAULT %s;\n",
				table.Identifier,
				desiredColumn.Name,
				desiredColumn.Default,
			))
		} else {
			// DROP DEFAULT
			df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" DROP DEFAULT;\n",
				table.Identifier,
				desiredColumn.Name,
			))
		}
	}
}

func findTable(tables Tables, identifier string) *Table {
	return tables[identifier]
}

type Tables map[string]*Table

func (tables Tables) AddTable(searchPath string, createTableStatement *ast.CreateTableStatement) {
	if createTableStatement.TableName.SchemaIdentifier == nil {
		createTableStatement.TableName.SetSchema(searchPath)
	}
	identifier := createTableStatement.TableName.String()

	columns := make(map[string]*Column)
	for _, columnDefinition := range createTableStatement.ColumnDefinitionList {
		col := columnFromAst(columnDefinition)
		columns[col.Name] = col
	}

	tables[identifier] = &Table{
		CreateTableStatement: createTableStatement,
		Identifier:           identifier,
		Columns:              columns,
		// Indexes: indexes // TODO
	}
}

// processDDL converts to schema and table mappings
func processDDL(ddl *ast.DataDefinition) Tables {
	var tables = make(Tables)
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
			// nop
		case *ast.CreateTableStatement:
			tables.AddTable(searchPath, stmt)
		case *ast.CreateIndexStatement:
			// TODO
			// 	tbls := df.sourceTables[df.searchPath]
			// 	t := tbls[stmt.TableName.Value]
			// 	t.CreateIndexes = append(t.CreateIndexes, stmt)
		default:
			log.Printf("skip statement: %v", stmt)
		}
	}
	return tables
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
