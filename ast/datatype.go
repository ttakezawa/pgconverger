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

func (*DataTypeInteger) Name() DataTypeName              { return Integer }
func (*DataTypeInteger) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("integer") }

type DataTypeBigint struct {
	Token token.Token
}

func (*DataTypeBigint) Name() DataTypeName              { return Bigint }
func (*DataTypeBigint) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("bigint") }

type DataTypeSmallint struct {
	Token token.Token
}

func (*DataTypeSmallint) Name() DataTypeName              { return Smallint }
func (*DataTypeSmallint) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("smallint") }

type DataTypeBigserial struct{}

func (*DataTypeBigserial) Name() DataTypeName              { return Bigserial }
func (*DataTypeBigserial) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("bigserial") }

type DataTypeBoolean struct{}

func (*DataTypeBoolean) Name() DataTypeName              { return Boolean }
func (*DataTypeBoolean) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("boolean") }

type DataTypeNumeric struct{}

func (*DataTypeNumeric) Name() DataTypeName              { return Numeric }
func (*DataTypeNumeric) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("numeric") }

type DataTypeOptionLength struct {
	token.Token
}

func (dataTypeOptionLength *DataTypeOptionLength) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("(")
	_, _ = w.WriteString(dataTypeOptionLength.Literal)
	_, _ = w.WriteString(")")
}

type DataTypeCharacter struct {
	Varying      bool
	OptionLength *DataTypeOptionLength
}

func (*DataTypeCharacter) Name() DataTypeName { return Character }
func (dataTypeCharacter *DataTypeCharacter) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("character")
	if dataTypeCharacter.Varying {
		_, _ = w.WriteString(" varying")
	}
	if dataTypeCharacter.OptionLength != nil {
		dataTypeCharacter.OptionLength.WriteStringTo(w)
	}
}

type DataTypeText struct{}

func (*DataTypeText) Name() DataTypeName              { return Text }
func (*DataTypeText) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("text") }

type DataTypeUuid struct{}

func (*DataTypeUuid) Name() DataTypeName              { return Uuid }
func (*DataTypeUuid) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("uuid") }

type DataTypeJsonb struct{}

func (*DataTypeJsonb) Name() DataTypeName              { return Jsonb }
func (*DataTypeJsonb) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("jsonb") }

type DataTypeBytea struct{}

func (*DataTypeBytea) Name() DataTypeName              { return Bytea }
func (*DataTypeBytea) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("bytea") }

type DataTypeTsvector struct{}

func (*DataTypeTsvector) Name() DataTypeName              { return Tsvector }
func (*DataTypeTsvector) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("tsvector") }

type DataTypeDate struct{}

func (*DataTypeDate) Name() DataTypeName              { return Date }
func (*DataTypeDate) WriteStringTo(w io.StringWriter) { _, _ = w.WriteString("date") }

type DataTypeTimestamp struct {
	WithTimeZone bool
}

func (*DataTypeTimestamp) Name() DataTypeName { return Timestamp }
func (dataTypeTimestamp *DataTypeTimestamp) WriteStringTo(w io.StringWriter) {
	_, _ = w.WriteString("timestamp")
	if dataTypeTimestamp.WithTimeZone {
		_, _ = w.WriteString(" with time zone")
	} else {
		_, _ = w.WriteString(" without time zone")
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
	Smallint
	// Not implemented: Path
	// Not implemented: PgLsn
	// Not implemented: Point
	// Not implemented: Polygon
	// Not implemented: Real
	// Not implemented: Smallserial
	Serial
	Text
	// Not implemented: Time
	// Not implemented: TimeWithTimeZone
	Timestamp
	// Not implemented: Tsquery
	Tsvector
	// Not implemented: TxidSnapshot
	Uuid
	// Not implemented: Xml
)
