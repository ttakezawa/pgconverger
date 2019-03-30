package ast

import (
	"io"

	"github.com/ttakezawa/pgconverger/token"
)

type Node interface {
	Source(w io.StringWriter)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) Source(w io.StringWriter) {
	_, _ = w.WriteString(`"`)
	_, _ = w.WriteString(identifier.Value)
	_, _ = w.WriteString(`"`)
}

type DataDefinition struct {
	StatementList []Statement
}

func (dataDefinition *DataDefinition) Source(w io.StringWriter) {
	for i, statement := range dataDefinition.StatementList {
		if i > 0 {
			_, _ = w.WriteString("\n")
		}
		statement.Source(w)
	}
}

type CreateSchemaStatement struct {
	Name *Identifier
}

func (*CreateSchemaStatement) statementNode() {}

func (createSchemaStatement *CreateSchemaStatement) Source(w io.StringWriter) {
	_, _ = w.WriteString("CREATE SCHEMA ")
	createSchemaStatement.Name.Source(w)
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
	TableName            *Identifier
	ColumnDefinitionList []*ColumnDefinition
}

func (*CreateTableStatement) statementNode() {}

func (createTableStatement *CreateTableStatement) Source(w io.StringWriter) {
	_, _ = w.WriteString("CREATE TABLE ")
	createTableStatement.TableName.Source(w)
	_, _ = w.WriteString(" (\n")

	for i, columnDefinition := range createTableStatement.ColumnDefinitionList {
		_, _ = w.WriteString("    ")
		columnDefinition.Source(w)
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

func (columnDefinition *ColumnDefinition) Source(w io.StringWriter) {
	columnDefinition.Name.Source(w)
	_, _ = w.WriteString(" ")
	columnDefinition.Type.Source(w)
	for _, constraint := range columnDefinition.ConstraintList {
		if _, ok := constraint.(*ColumnConstraintNull); ok {
			// When Constraint Null, skip it.
			continue
		}
		_, _ = w.WriteString(" ")
		constraint.Source(w)
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

func (*ColumnConstraintNotNull) Source(w io.StringWriter) {
	_, _ = w.WriteString("NOT NULL")
}

// NULL
type ColumnConstraintNull struct {
}

func (*ColumnConstraintNull) Source(w io.StringWriter) {
	_, _ = w.WriteString("NULL")
}

// DEFAULT expr
type ColumnConstraintDefault struct {
	Expr Expression
}

func (columnConstraintDefault *ColumnConstraintDefault) Source(w io.StringWriter) {
	_, _ = w.WriteString("DEFAULT ")
	columnConstraintDefault.Expr.Source(w)
}

type InfixExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (infixExpr *InfixExpression) expressionNode() {}
func (infixExpr *InfixExpression) Source(w io.StringWriter) {
	infixExpr.Left.Source(w)
	_, _ = w.WriteString(infixExpr.Operator.Literal)
	infixExpr.Right.Source(w)
}

type StringLiteral struct {
	Token token.Token
}

func (il *StringLiteral) expressionNode() {}
func (il *StringLiteral) Source(w io.StringWriter) {
	_, _ = w.WriteString(il.Token.Literal)
}

type NumberLiteral struct {
	Token token.Token
}

func (il *NumberLiteral) expressionNode() {}
func (il *NumberLiteral) Source(w io.StringWriter) {
	_, _ = w.WriteString(il.Token.Literal)
}
