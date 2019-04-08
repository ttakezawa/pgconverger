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
	if infixExpr.Operator.Type == token.Is {
		_, _ = w.WriteString(" IS ")
	} else {
		_, _ = w.WriteString(infixExpr.Operator.Literal)
	}
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

type BooleanLiteral struct {
	Token token.Token
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) Source(w io.StringWriter) {
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
func (nl *NullLiteral) Source(w io.StringWriter) {
	_, _ = w.WriteString(nl.Token.Literal)
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
	TableName    *Identifier
	UsingMethod  *Identifier
	IndexTargets []Node // Slice of (Identifier OR Expression)
}

func (*CreateIndexStatement) statementNode() {}

func (createIndexStatement *CreateIndexStatement) Source(w io.StringWriter) {
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
		createIndexStatement.Name.Source(w)
		_, _ = w.WriteString(" ")
	}
	_, _ = w.WriteString("ON ")
	createIndexStatement.TableName.Source(w)
	_, _ = w.WriteString(" ")
	if createIndexStatement.UsingMethod != nil {
		_, _ = w.WriteString("USING ")
		createIndexStatement.UsingMethod.Source(w)
		_, _ = w.WriteString(" ")
	}
	_, _ = w.WriteString("(")
	for i, indexTarget := range createIndexStatement.IndexTargets {
		if i != 0 {
			_, _ = w.WriteString(", ")
		}
		switch t := indexTarget.(type) {
		case *Identifier:
			t.Source(w)
		case Expression:
			_, _ = w.WriteString("(")
			t.Source(w)
			_, _ = w.WriteString(")")
		case nil:
			// nop
		default:
			t.Source(w)
		}
	}
	_, _ = w.WriteString(");")
}
