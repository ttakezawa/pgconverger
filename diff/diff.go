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
}

func Process(source fileReader, desired fileReader) (*Diff, error) {
	df := &Diff{
		source:  source,
		desired: desired,
	}
	df.sourceErrors, df.sourceDDL = parseOneSide(df.source)
	df.desiredErrors, df.desiredDDL = parseOneSide(df.desired)
	if df.ErrorOrNil() != nil {
		return df, df.ErrorOrNil()
	}

	df.generatePatch()
	return df, nil
}

func parseOneSide(reader fileReader) (errs []error, ddl *ast.DataDefinition) {
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

func (df *Diff) String() string {
	var builder strings.Builder
	df.desiredDDL.WriteStringTo(&builder)
	return builder.String()
}

func (df *Diff) generatePatch() {
	var a analyzer
	a.sourceSchemas, a.sourceTables = processDDL(df.sourceDDL)
	a.desiredSchemas, a.desiredTables = processDDL(df.desiredDDL)

	patch := a.generatePatch()
	log.Printf("%s", patch)
}

type analyzer struct {
	sourceSchemas map[string]struct{}
	sourceTables  map[string]map[string]*Table

	desiredSchemas map[string]struct{}
	desiredTables  map[string]map[string]*Table
}

type (
	Table struct {
		Schema  string
		Name    string
		Columns map[string]*Column
		Indexes map[string]*Index
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

func (a *analyzer) generatePatch() string {
	var builder strings.Builder
	for schema, tables := range a.sourceTables {
		for _, table := range tables {
			desiredTable := findTable(a.desiredTables, schema, table.Name)
			builder.WriteString(diffTable(table, desiredTable))
		}
	}
	return builder.String()
}

func diffTable(sourceTable, desiredTable *Table) string {
	var builder strings.Builder
	for _, desiredColumn := range desiredTable.Columns {
		_, ok := sourceTable.Columns[desiredColumn.Name] //
		if !ok {
			builder.WriteString(
				fmt.Sprintf("ALTER TABLE %s ADD %s %s;\n",
					sourceTable.Name,
					desiredColumn.Name,
					desiredColumn.DataType,
				),
			)
		}
	}
	return builder.String()
}

func findTable(tables map[string]map[string]*Table, schema string, tableName string) *Table {
	return tables[schema][tableName]
}

// processDDL converts to schema and table mappings
func processDDL(ddl *ast.DataDefinition) (schemas map[string]struct{}, tables map[string]map[string]*Table) {
	schemas = make(map[string]struct{})
	tables = make(map[string]map[string]*Table)
	currentSchema := "public"

	for _, statement := range ddl.StatementList {
		switch stmt := statement.(type) {
		case *ast.SetStatement:
			if stmt.Name.Value == "search_path" {
				ident, ok := stmt.Values[0].(*ast.Identifier)
				if ok {
					currentSchema = ident.Value
				}
			}
		case *ast.CreateSchemaStatement:
			schemas[stmt.Name.Value] = struct{}{}
		case *ast.CreateTableStatement:
			tbls, ok := tables[currentSchema]
			if !ok {
				tbls = make(map[string]*Table)
				tables[currentSchema] = tbls
			}
			columns := make(map[string]*Column)
			for _, columnDefinition := range stmt.ColumnDefinitionList {
				col := columnFromAst(columnDefinition)
				columns[col.Name] = col
			}
			tbls[stmt.TableName.Value] = &Table{
				Schema:  currentSchema,
				Name:    stmt.TableName.Value,
				Columns: columns,
			}
		case *ast.CreateIndexStatement:
			// TODO
			// 	tbls := a.sourceTables[a.currentSchema]
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
