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

type DataDefinition struct {
	StatementList []Statement
}

func (dataDefinition *DataDefinition) Source(w io.StringWriter) {
	for i, statement := range dataDefinition.StatementList {
		if i > 0 {
			_, _ = w.WriteString("\n\n")
		}
		statement.Source(w)
	}
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
	TableName            token.Token
	ColumnDefinitionList []*ColumnDefinition
}

func (*CreateTableStatement) statementNode() {}

func identifierToSource(tok token.Token) string {
	if len(tok.Literal) == 0 {
		return `""`
	}
	if tok.Literal[0] == '"' {
		return tok.Literal
	}
	return `"` + tok.Literal + `"`
}

func (createTableStatement *CreateTableStatement) Source(w io.StringWriter) {
	_, _ = w.WriteString("CREATE TABLE ")
	_, _ = w.WriteString(identifierToSource(createTableStatement.TableName))
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
	Name           token.Token
	Type           DataType
	ConstraintList []ColumnConstraint
}

func (columnDefinition *ColumnDefinition) Source(w io.StringWriter) {
	_, _ = w.WriteString(identifierToSource(columnDefinition.Name))
	_, _ = w.WriteString(" ")
	columnDefinition.Type.Source(w)
}

type DataType interface {
	Node
	Name() DataTypeName
}

type DataTypeBigint struct {
	Token token.Token
}

func (*DataTypeBigint) Name() DataTypeName       { return Bigint }
func (*DataTypeBigint) Source(w io.StringWriter) { _, _ = w.WriteString("bigint") }

type DataTypeBigserial struct{}

func (*DataTypeBigserial) Name() DataTypeName       { return Bigserial }
func (*DataTypeBigserial) Source(w io.StringWriter) { _, _ = w.WriteString("bigserial") }

type DataTypeOptionLength struct {
	token.Token
}

func (dataTypeOptionLength *DataTypeOptionLength) Source(w io.StringWriter) {
	_, _ = w.WriteString("(")
	_, _ = w.WriteString(dataTypeOptionLength.Literal)
	_, _ = w.WriteString(")")
}

type DataTypeCharacter struct {
	Varying      bool
	OptionLength *DataTypeOptionLength
}

func (*DataTypeCharacter) Name() DataTypeName { return Character }
func (dataTypeCharacter *DataTypeCharacter) Source(w io.StringWriter) {
	_, _ = w.WriteString("character")
	if dataTypeCharacter.Varying {
		_, _ = w.WriteString(" varying")
	}
	if dataTypeCharacter.OptionLength != nil {
		dataTypeCharacter.OptionLength.Source(w)
	}
}

type DataTypeTimestamp struct {
	WithTimeZone bool
}

func (*DataTypeTimestamp) Name() DataTypeName { return Timestamp }
func (dataTypeTimestamp *DataTypeTimestamp) Source(w io.StringWriter) {
	_, _ = w.WriteString("timestamp")
	if dataTypeTimestamp.WithTimeZone {
		_, _ = w.WriteString(" with time zone")
	}
}

// //go:generate stringer -type=DataTypeName
type DataTypeName int

const (
	Bigint DataTypeName = iota
	Bigserial
	// Not implemented: Bit
	// Not implemented: BitVarying
	Boolean
	// Not implemented: Box
	Bytea
	Character
	// Not implemented: Cidr
	// Not implemented: Circle
	Date
	// Not implemented: DoublePrecision
	// Not implemented: Inet
	Integer
	// Not implemented: Interval
	// Not implemented: Json
	Jsonb
	// Not implemented: Line
	// Not implemented: Lseg
	// Not implemented: Macaddr
	// Not implemented: Macaddr8
	// Not implemented: Money
	Numeric
	// Not implemented: Path
	// Not implemented: PgLsn
	// Not implemented: Point
	// Not implemented: Polygon
	// Not implemented: Real
	// Not implemented: Smallint
	// Not implemented: Smallserial
	Serial
	Text
	// Not implemented: Time
	// Not implemented: TimeWithTimeZone
	Timestamp
	// Not implemented: Tsquery
	Tsvector
	// Not implemented: TxidSnapshot
	// Not implemented: UUID
	// Not implemented: Xml
)

// { NOT NULL |
//   NULL |
//   DEFAULT default_expr |
//   UNIQUE index_parameters |
//   PRIMARY KEY index_parameters
// }
type ColumnConstraint interface {
}

// NOT NULL
type ColumnConstraintNotNull struct {
}

// DEFAULT default_expr
type ColumnConstraintDefault struct {
	Expr DefaultExpr
}

type DefaultExpr interface{}
