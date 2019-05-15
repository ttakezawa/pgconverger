package diff

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
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
	Tables                 map[string]*Table
	Indexes                map[string]*Index
	TableConstraints       map[string]*TableConstraint
	AlterColumnSetDefaults map[string]*AlterColumnSetDefault

	Table struct {
		CreateTableStatement *ast.CreateTableStatement
		string
		Identifier             string
		Columns                map[string]*Column
		Indexes                Indexes
		TableConstraints       TableConstraints
		AlterColumnSetDefaults AlterColumnSetDefaults
	}

	Column struct {
		Name         string
		DataType     string
		NotNull      bool
		Default      string
		SequenceName string
	}

	Index struct {
		CreateIndexStatement *ast.CreateIndexStatement
		Name                 string
	}

	TableConstraint struct {
		Name    string
		Type    ConstraintType
		Columns []string
	}

	AlterColumnSetDefault struct {
		Column  string
		Default string
	}
)

type ConstraintType string

const (
	Unknown    ConstraintType = ""
	Unique     ConstraintType = "UNIQUE"
	PrimaryKey ConstraintType = "PRIMARY KEY"
)

func (tables Tables) SortedKeys() (keys []string) {
	for k := range tables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func (df *Diff) writeTableAnnotation(table *Table) {
	df.stringBuilder.WriteString("-- Table: " + table.Identifier + "\n")
}

func (df *Diff) generatePatch() string {
	df.sourceTables = processDDL(df.sourceDDL)
	df.desiredTables = processDDL(df.desiredDDL)

	for _, identifier := range df.sourceTables.SortedKeys() {
		sourceTable := df.sourceTables[identifier]
		desiredTable := df.desiredTables.FindTable(identifier)
		if desiredTable != nil {
			origBuilder := df.stringBuilder
			tmpBuilder := &strings.Builder{}
			df.stringBuilder = tmpBuilder
			df.diffTable(sourceTable, desiredTable)
			df.stringBuilder = origBuilder
			if tmpBuilder.Len() > 0 {
				df.writeTableAnnotation(sourceTable)
				df.stringBuilder.WriteString(tmpBuilder.String())
				df.stringBuilder.WriteString("\n")
			}
		} else {
			df.writeTableAnnotation(sourceTable)
			df.dropTable(sourceTable)
			df.stringBuilder.WriteString("\n")
		}
	}

	for _, identifier := range df.desiredTables.SortedKeys() {
		desiredTable := df.desiredTables[identifier]
		sourceTable := df.sourceTables.FindTable(identifier)
		if sourceTable != nil {
			// none
		} else {
			df.writeTableAnnotation(desiredTable)
			df.createTable(desiredTable)
			df.stringBuilder.WriteString("\n")
		}
	}

	return df.stringBuilder.String()
}

func (df *Diff) createTable(table *Table) {
	table.CreateTableStatement.WriteStringTo(df.stringBuilder)
	for _, column := range table.Columns {
		if column.SequenceName != "" {
			df.addSequence(table, column)
		}
	}
	for _, index := range table.Indexes {
		index.CreateIndexStatement.WriteStringTo(df.stringBuilder)
		df.stringBuilder.WriteString("\n")
	}
	for _, constraint := range table.TableConstraints {
		df.addTableConstraint(table, constraint)
	}
	for _, alterColumnSetDefault := range table.AlterColumnSetDefaults {
		df.addAlterColumnSetDefault(table, alterColumnSetDefault)
	}
}

func (df *Diff) dropTable(table *Table) {
	df.WriteString(fmt.Sprintf("DROP TABLE %s;\n", table.Identifier))
}

// generate DDL for a table which exists in both.
func (df *Diff) diffTable(sourceTable, desiredTable *Table) {
	for _, sourceColumn := range sourceTable.Columns {
		desiredColumn, ok := desiredTable.Columns[sourceColumn.Name]
		if ok {
			df.alterColumn(sourceTable, sourceColumn, desiredColumn)
		} else {
			df.dropColumn(sourceTable, sourceColumn)
		}
	}

	for _, desiredColumn := range desiredTable.Columns {
		_, ok := sourceTable.Columns[desiredColumn.Name]
		if !ok {
			df.addColumn(sourceTable, desiredColumn)
		}
	}

	for _, sourceIndex := range sourceTable.Indexes {
		desiredIndex, ok := desiredTable.Indexes[sourceIndex.Name]
		if ok {
			// TODO: ALTER INDEX ?
			_ = desiredIndex
		} else {
			df.dropIndex(sourceTable, sourceIndex)
		}
	}

	for _, desiredIndex := range desiredTable.Indexes {
		_, ok := sourceTable.Indexes[desiredIndex.Name]
		if !ok {
			df.createIndex(sourceTable, desiredIndex)
		}
	}

	for _, sourceTableConstraint := range sourceTable.TableConstraints {
		desiredTableConstraint, ok := desiredTable.TableConstraints[sourceTableConstraint.Name]
		if ok {
			// TODO: MODIFY CONSTRAINT ?
			_ = desiredTableConstraint
		} else {
			df.dropTableConstraint(sourceTable, sourceTableConstraint)
		}
	}

	for _, desiredTableConstraint := range desiredTable.TableConstraints {
		_, ok := sourceTable.TableConstraints[desiredTableConstraint.Name]
		if !ok {
			df.addTableConstraint(sourceTable, desiredTableConstraint)
		}
	}

	for _, sourceAlterColumnSetDefault := range sourceTable.AlterColumnSetDefaults {
		desiredAlterColumnSetDefault, ok := desiredTable.AlterColumnSetDefaults[sourceAlterColumnSetDefault.Column]
		if ok {
			// TODO: MODIFY CONSTRAINT ?
			_ = desiredAlterColumnSetDefault
		} else {
			df.dropAlterColumnSetDefault(sourceTable, sourceAlterColumnSetDefault)
		}
	}

	for _, desiredAlterColumnSetDefault := range desiredTable.AlterColumnSetDefaults {
		_, ok := sourceTable.AlterColumnSetDefaults[desiredAlterColumnSetDefault.Column]
		if !ok {
			df.addAlterColumnSetDefault(sourceTable, desiredAlterColumnSetDefault)
		}
	}
}

func (df *Diff) addColumn(table *Table, column *Column) {
	df.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" %s;\n",
		table.Identifier,
		column.Name,
		column.DataType,
	))
}

func (df *Diff) addSequence(table *Table, column *Column) {
	var (
		schemaName = table.CreateTableStatement.TableName.SchemaIdentifier.Value
		tableName  = table.CreateTableStatement.TableName.TableIdentifier.Value
		columnName = column.Name
		sequence   = column.SequenceName
	)
	df.WriteString(
		fmt.Sprintf(`CREATE SEQUENCE "%s"."%s"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE "%s"."%s" OWNED BY "%s"."%s";`,
			schemaName, sequence,
			schemaName, sequence, tableName, columnName),
	)
	df.WriteString("\n")
}

func (df *Diff) addTableConstraint(table *Table, tableConstraint *TableConstraint) {
	df.WriteString(fmt.Sprintf(
		`ALTER TABLE ONLY %s ADD CONSTRAINT "%s" %s (`,
		table.Identifier,
		tableConstraint.Name,
		tableConstraint.Type,
	))

	for i, c := range tableConstraint.Columns {
		if i != 0 {
			df.WriteString(", ")
		}
		df.WriteString(fmt.Sprintf(`"%s"`, c))
	}
	df.WriteString(");\n")
}

func (df *Diff) addAlterColumnSetDefault(table *Table, alterColumnSetDefault *AlterColumnSetDefault) {
	df.WriteString(fmt.Sprintf(
		"ALTER TABLE ONLY %s ALTER COLUMN \"%s\" SET DEFAULT %s;\n",
		table.Identifier,
		alterColumnSetDefault.Column,
		alterColumnSetDefault.Default,
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
		df.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN \"%s\" TYPE %s;\n",
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

func (df *Diff) createIndex(_ *Table, index *Index) {
	index.CreateIndexStatement.WriteStringTo(df.stringBuilder)
	df.WriteString("\n")
}

func (df *Diff) dropIndex(table *Table, index *Index) {
	df.WriteString(fmt.Sprintf("DROP INDEX \"%s\";\n",
		index.Name,
	))
}

func (df *Diff) dropTableConstraint(table *Table, tableConstraint *TableConstraint) {
	df.WriteString(fmt.Sprintf("ALTER TABLE ONLY %s DROP CONSTRAINT \"%s\";\n",
		table.Identifier,
		tableConstraint.Name,
	))
}

func (df *Diff) dropAlterColumnSetDefault(table *Table, alterColumnSetDefault *AlterColumnSetDefault) {
	df.WriteString(fmt.Sprintf("ALTER TABLE ONLY %s ALTER COLUMN \"%s\" DROP DEFAULT;\n",
		table.Identifier,
		alterColumnSetDefault.Column,
	))
}

func (tables Tables) FindTable(identifier string) *Table {
	return tables[identifier]
}

func (tables Tables) FindTableBy(searchPath, tableName string) *Table {
	return tables[`"`+searchPath+`"."`+tableName+`"`]
}

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
		CreateTableStatement:   createTableStatement,
		Identifier:             identifier,
		Columns:                columns,
		Indexes:                make(Indexes),
		TableConstraints:       make(TableConstraints),
		AlterColumnSetDefaults: make(AlterColumnSetDefaults),
	}
}

func (tables Tables) AddIndex(searchPath string, createIndexStatement *ast.CreateIndexStatement) {
	if createIndexStatement.TableName.SchemaIdentifier == nil {
		createIndexStatement.TableName.SetSchema(searchPath)
	}
	tableName := createIndexStatement.TableName.String()
	table := tables.FindTable(tableName)
	if table == nil {
		log.Printf("irregular create index to unknown table=%s", tableName)
		return
	}
	indexName := createIndexStatement.Name.Value
	table.Indexes[indexName] = &Index{
		CreateIndexStatement: createIndexStatement,
		Name:                 indexName,
	}
}

func (tables Tables) AddSequence(searchPath string, alterSequenceStatement *ast.AlterSequenceStatement) {
	tableName := alterSequenceStatement.OwnedByTable.TableIdentifier.Value
	table := tables.FindTableBy(searchPath, tableName)
	if table == nil {
		log.Printf("irregular alter sequence to unknown table=%s", tableName)
		return
	}
	columnName := alterSequenceStatement.OwnedByColumn.Value
	column := table.Columns[columnName]
	if column == nil {
		log.Printf("irregular alter sequence to unknown column=%s", columnName)
		return
	}
	column.SequenceName = alterSequenceStatement.Name.SequenceIdentifier.Value
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
			tables.AddIndex(searchPath, stmt)
		case *ast.AlterSequenceStatement:
			tables.AddSequence(searchPath, stmt)
		case *ast.AlterTableStatement:
			processAlterTableStatement(searchPath, tables, stmt)
		default:
			log.Printf("skip statement: %v", stmt)
		}
	}
	return tables
}

func processAlterTableStatement(searchPath string, tables Tables, alterTableStatement *ast.AlterTableStatement) {
	if alterTableStatement.Name.SchemaIdentifier == nil {
		alterTableStatement.Name.SetSchema(searchPath)
	}
	tableName := alterTableStatement.Name.String()
	table := tables.FindTable(tableName)
	if table == nil {
		log.Printf("irregular alter table to unknown table=%s", tableName)
		return
	}
	for _, action := range alterTableStatement.Actions {
		switch v := action.(type) {
		case *ast.TableConstraint:
			var typ ConstraintType
			if v.PrimaryKey {
				typ = PrimaryKey
			}
			if v.Unique {
				typ = Unique
			}
			if typ == Unknown {
				continue
			}
			var columns []string
			for _, c := range v.ColumnList.ColumnNames {
				columns = append(columns, c.Value)
			}
			table.TableConstraints[v.Name.Value] = &TableConstraint{
				Name:    v.Name.Value,
				Type:    typ,
				Columns: columns,
			}
		case *ast.AlterColumnSetDefault:
			var builder strings.Builder
			v.Expr.WriteStringTo(&builder)
			table.AlterColumnSetDefaults[v.Column.Value] = &AlterColumnSetDefault{
				Column:  v.Column.Value,
				Default: builder.String(),
			}
		default:
			log.Printf("skipped table statement")
			return
		}
	}
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
