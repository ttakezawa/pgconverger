package ast

import (
	"io"
	"strings"

	"github.com/ttakezawa/pgconverger/token"
)

type Node interface {
	WriteStringTo(w io.StringWriter)
}

func FormatNode(node Node) string {
	var builder strings.Builder
	node.WriteStringTo(&builder)
	return builder.String()
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type TableName struct {
	SchemaIdentifier *Identifier
	TableIdentifier  *Identifier
}

func (tableName *TableName) WriteStringTo(w io.StringWriter) {
	if tableName.SchemaIdentifier != nil {
		tableName.SchemaIdentifier.WriteStringTo(w)
		_, _ = w.WriteString(`.`)
	}
	tableName.TableIdentifier.WriteStringTo(w)
}

func (tableName *TableName) String() string {
	var builder strings.Builder
	tableName.WriteStringTo(&builder)
	return builder.String()
}

func (tableName *TableName) SetSchema(schema string) {
	// TODO: refactor
	tableName.SchemaIdentifier = &Identifier{
		Token: token.Token{
			Type:    token.Identifier,
			Literal: `"` + schema + `"`,
			Line:    0,
		},
		Value: schema,
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString(`"`)
	_, _ = w.WriteString(identifier.Value)
	_, _ = w.WriteString(`"`)
}

type DataDefinition struct {
	StatementList []Statement
}

func (dataDefinition *DataDefinition) WriteStringTo(w io.StringWriter) {
	for i, statement := range dataDefinition.StatementList {
		if i > 0 {
			_, _ = w.WriteString("\n")
		}
		statement.WriteStringTo(w)
	}
}

type CreateSchemaStatement struct {
	Name *Identifier
}

func (*CreateSchemaStatement) statementNode() {}

func (createSchemaStatement *CreateSchemaStatement) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("CREATE SCHEMA ")
	createSchemaStatement.Name.WriteStringTo(w)
	_, _ = w.WriteString(";\n")
}

// CREATE [ [ GLOBAL | LOCAL ] { TEMPORARY | TEMP } | UNLOGGED ] TABLE [ IF NOT EXISTS ] table_name ( [
//   { column_name data_type [ COLLATE collation ] [ column_constraint [ ... ] ]
//     | table_constraint
//     | LIKE source_table [ like_option ... ] }
//     [, ... ]
// ] )
// [ INHERITS ( parent_table [, ... ] ) ]
// [ PARTITION BY { RANGE | LIST } ( { column_name | ( expression ) } [ COLLATE collation ] [ opclass ] [, ... ] ) ]
// [ WITH ( storage_parameter [= value] [, ... ] ) | WITH OIDS | WITHOUT OIDS ]
// [ ON COMMIT { PRESERVE ROWS | DELETE ROWS | DROP } ]
// [ TABLESPACE tablespace_name ]
type CreateTableStatement struct {
	TableName            *TableName
	ColumnDefinitionList []*ColumnDefinition
}

func (*CreateTableStatement) statementNode() {}

func (createTableStatement *CreateTableStatement) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("CREATE TABLE ")
	createTableStatement.TableName.WriteStringTo(w)
	_, _ = w.WriteString(" (\n")

	for i, columnDefinition := range createTableStatement.ColumnDefinitionList {
		_, _ = w.WriteString("    ")
		columnDefinition.WriteStringTo(w)
		if i < len(createTableStatement.ColumnDefinitionList)-1 {
			_, _ = w.WriteString(",")
		}
		_, _ = w.WriteString("\n")
	}
	_, _ = w.WriteString(");\n")
}

type ColumnDefinition struct {
	Name           *Identifier
	Type           DataType
	ConstraintList []ColumnConstraint
}

func (columnDefinition *ColumnDefinition) WriteStringTo(w io.StringWriter) {
	columnDefinition.Name.WriteStringTo(w)
	_, _ = w.WriteString(" ")
	columnDefinition.Type.WriteStringTo(w)
	for _, constraint := range columnDefinition.ConstraintList {
		if _, ok := constraint.(*ColumnConstraintNull); ok {
			// When Constraint Null, skip it.
			continue
		}
		_, _ = w.WriteString(" ")
		constraint.WriteStringTo(w)
	}
}

// { NOT NULL |
//   NULL |
//   DEFAULT expr |
//   UNIQUE index_parameters |
//   PRIMARY KEY index_parameters
// }
type ColumnConstraint interface {
	Node
}

// NOT NULL
type ColumnConstraintNotNull struct {
}

func (*ColumnConstraintNotNull) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("NOT NULL")
}

// NULL
type ColumnConstraintNull struct {
}

func (*ColumnConstraintNull) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("NULL")
}

// DEFAULT expr
type ColumnConstraintDefault struct {
	Expr Expression
}

func (columnConstraintDefault *ColumnConstraintDefault) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("DEFAULT ")
	columnConstraintDefault.Expr.WriteStringTo(w)
}

type IndexTarget struct {
	Node   Node
	IsDesc bool
}

func (indexTarget *IndexTarget) WriteStringTo(w io.StringWriter) {
	indexTarget.Node.WriteStringTo(w)
	if indexTarget.IsDesc {
		_, _ = w.WriteString(" DESC")
	}
}

type InfixExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (infixExpr *InfixExpression) expressionNode() {}
func (infixExpr *InfixExpression) WriteStringTo(w io.StringWriter) {
	infixExpr.Left.WriteStringTo(w)
	if infixExpr.Operator.Type == token.Is {
		_, _ = w.WriteString(" IS ")
	} else {
		_, _ = w.WriteString(infixExpr.Operator.Literal)
	}
	infixExpr.Right.WriteStringTo(w)
}

type StringLiteral struct {
	Token token.Token
}

func (il *StringLiteral) expressionNode() {}
func (il *StringLiteral) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString(il.Token.Literal)
}

type NumberLiteral struct {
	Token token.Token
}

func (il *NumberLiteral) expressionNode() {}
func (il *NumberLiteral) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString(il.Token.Literal)
}

type BooleanLiteral struct {
	Token token.Token
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) WriteStringTo(w io.StringWriter) {
	if bl.IsTrue() {
		_, _ = w.WriteString("TRUE")
	} else {
		_, _ = w.WriteString("FALSE")
	}
}

func (bl *BooleanLiteral) IsTrue() bool {
	return bl.Token.Type == token.True
}

type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString(nl.Token.Literal)
}

type GroupedExpression struct {
	Expression Expression
}

func (groupedExpr *GroupedExpression) expressionNode() {}
func (groupedExpr *GroupedExpression) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("(")
	groupedExpr.Expression.WriteStringTo(w)
	_, _ = w.WriteString(")")
}

// CREATE [ UNIQUE ] INDEX [ CONCURRENTLY ] [ [ IF NOT EXISTS ] name ] ON table_name [ USING method ]
//     ( { column_name | ( expression ) } [ COLLATE collation ] [ opclass ] [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...] )
//     [ WITH ( storage_parameter = value [, ... ] ) ]
//     [ TABLESPACE tablespace_name ]
//     [ WHERE predicate ]
type CreateIndexStatement struct {
	UniqueIndex  bool
	Concurrently bool
	IfNotExists  bool
	Name         *Identifier
	TableName    *TableName
	UsingMethod  *Identifier
	IndexTargets []*IndexTarget // Slice of (Identifier OR Expression)
}

func (*CreateIndexStatement) statementNode() {}

func (createIndexStatement *CreateIndexStatement) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("CREATE ")
	if createIndexStatement.UniqueIndex {
		_, _ = w.WriteString("UNIQUE ")
	}
	_, _ = w.WriteString("INDEX ")
	if createIndexStatement.Concurrently {
		_, _ = w.WriteString("CONCURRENTLY ")
	}
	if createIndexStatement.IfNotExists {
		_, _ = w.WriteString("IF NOT EXISTS ")
	}
	if createIndexStatement.Name != nil {
		createIndexStatement.Name.WriteStringTo(w)
		_, _ = w.WriteString(" ")
	}
	_, _ = w.WriteString("ON ")
	createIndexStatement.TableName.WriteStringTo(w)
	_, _ = w.WriteString(" ")
	if createIndexStatement.UsingMethod != nil {
		_, _ = w.WriteString("USING ")
		createIndexStatement.UsingMethod.WriteStringTo(w)
		_, _ = w.WriteString(" ")
	}
	_, _ = w.WriteString("(")
	for i, indexTarget := range createIndexStatement.IndexTargets {
		if i != 0 {
			_, _ = w.WriteString(", ")
		}
		if indexTarget != nil {
			indexTarget.WriteStringTo(w)
		}
	}
	_, _ = w.WriteString(");")
}

type SetStatement struct {
	Name   *Identifier
	Values []Expression
}

func (*SetStatement) statementNode() {}

func (setStatement *SetStatement) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("SET ")
	_, _ = w.WriteString(setStatement.Name.Value)
	_, _ = w.WriteString(" = ")
	for i, value := range setStatement.Values {
		if i != 0 {
			_, _ = w.WriteString(", ")
		}
		ident, ok := value.(*Identifier)
		if ok && ident.Value == "pg_catalog" {
			_, _ = w.WriteString(ident.Value)
		} else {
			value.WriteStringTo(w)
		}
	}
	_, _ = w.WriteString(";")
}

// CREATE SEQUENCE "users_id_seq"
//     START WITH 1
//     INCREMENT BY 1
//     NO MINVALUE
//     NO MAXVALUE
//     CACHE 1;
type CreateSequenceStatement struct {
	Name        *Identifier
	StartWith   Expression
	IncrementBy Expression
	NoMinvalue  bool
	NoMaxvalue  bool
	Cache       Expression
}

func (*CreateSequenceStatement) statementNode() {}

func (createSequenceStatement *CreateSequenceStatement) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("CREATE SEQUENCE ")
	createSequenceStatement.Name.WriteStringTo(w)
	if createSequenceStatement.StartWith != nil {
		_, _ = w.WriteString("\n    START WITH ")
		createSequenceStatement.StartWith.WriteStringTo(w)
	}
	if createSequenceStatement.IncrementBy != nil {
		_, _ = w.WriteString("\n    INCREMENT BY ")
		createSequenceStatement.IncrementBy.WriteStringTo(w)
	}
	if createSequenceStatement.NoMinvalue {
		_, _ = w.WriteString("\n    NO MINVALUE")
	}
	if createSequenceStatement.NoMaxvalue {
		_, _ = w.WriteString("\n    NO MAXVALUE")
	}
	if createSequenceStatement.Cache != nil {
		_, _ = w.WriteString("\n    CACHE ")
		createSequenceStatement.Cache.WriteStringTo(w)
	}
	_, _ = w.WriteString(";")
}
