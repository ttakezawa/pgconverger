package ast

import (
	"io"

	"github.com/ttakezawa/pgconverger/token"
)

type DataType interface {
	Node
	Name() DataTypeName
}

type DataTypeInteger struct {
	Token token.Token
}

func (*DataTypeInteger) Name() DataTypeName       { return Integer }
func (*DataTypeInteger) Source(w io.StringWriter) { _, _ = w.WriteString("integer") }

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

type DataTypeText struct{}

func (*DataTypeText) Name() DataTypeName       { return Text }
func (*DataTypeText) Source(w io.StringWriter) { _, _ = w.WriteString("text") }

type DataTypeJsonb struct{}

func (*DataTypeJsonb) Name() DataTypeName       { return Jsonb }
func (*DataTypeJsonb) Source(w io.StringWriter) { _, _ = w.WriteString("jsonb") }

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
