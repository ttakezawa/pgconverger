package ast

import (
	"fmt"
	"strings"

	"github.com/ttakezawa/pgconverger/token"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type DataDefinition struct {
	StatementList []Statement
}

func (dataDefinition *DataDefinition) String() string {
	var b strings.Builder
	for _, statement := range dataDefinition.StatementList {
		b.WriteString(statement.String())
	}
	return b.String()
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
	TableName            string
	ColumnDefinitionList []*ColumnDefinition
}

func (*CreateTableStatement) statementNode() {}
func (createTableStatement *CreateTableStatement) String() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "CREATE TABLE '%s' ();", createTableStatement.TableName)
	return b.String()
}

type ColumnDefinition struct {
	Name           token.Token
	Type           DataType
	ConstraintList []ColumnConstraint
}

type DataType interface {
	Name() DataTypeName
}

type DataTypeBigint struct {
	Token token.Token
}

func (*DataTypeBigint) Name() DataTypeName { return Bigint }

type DataTypeBigserial struct{}

func (*DataTypeBigserial) Name() DataTypeName { return Bigserial }

type DataTypeCharacter struct {
	Varying bool
	Length  *token.Token
}

func (*DataTypeCharacter) Name() DataTypeName { return Character }

type DataTypeTimestamp struct {
	WithTimeZone bool
}

func (*DataTypeTimestamp) Name() DataTypeName { return Timestamp }

//go:generate stringer -type=DataTypeName
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
