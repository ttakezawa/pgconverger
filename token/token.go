package token

import (
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	Illegal TokenType = iota
	EOF
	Space
	Comment
	Identifier
	String
	Number
	Semicolon
	Comma
	LParen
	RParen
	Typecast

	Add
	Alter
	By
	Cache
	Column
	Constraint
	Create
	Default
	Grant
	Increment
	Index
	Insert
	Maxvalue
	Minvalue
	No
	Not
	Null
	On
	Only
	Owner
	Select
	Sequence
	Set
	Start
	Table
	To
	Update
	Using
	Varying
	With

	Bigint
	Bigserial
	Boolean
	Bytea
	Character
	Date
	Integer
	Jsonb
	Numeric
	Serial
	Text
	Timestamp
	Time
	Tsvector
)

var keywords = map[string]TokenType{
	"ADD":        Add,
	"ALTER":      Alter,
	"BIGINT":     Bigint,
	"BIGSERIAL":  Bigserial,
	"BOOLEAN":    Boolean,
	"BY":         By,
	"BYTEA":      Bytea,
	"CACHE":      Cache,
	"CHARACTER":  Character,
	"COLUMN":     Column,
	"CONSTRAINT": Constraint,
	"CREATE":     Create,
	"DATE":       Date,
	"DEFAULT":    Default,
	"GRANT":      Grant,
	"INCREMENT":  Increment,
	"INDEX":      Index,
	"INSERT":     Insert,
	"INTEGER":    Integer,
	"JSONB":      Jsonb,
	"MAXVALUE":   Maxvalue,
	"MINVALUE":   Minvalue,
	"NO":         No,
	"NOT":        Not,
	"NULL":       Null,
	"NUMERIC":    Numeric,
	"ON":         On,
	"ONLY":       Only,
	"OWNER":      Owner,
	"SELECT":     Select,
	"SEQUENCE":   Sequence,
	"SERIAL":     Serial,
	"SET":        Set,
	"START":      Start,
	"TABLE":      Table,
	"TEXT":       Text,
	"TIME":       Time,
	"TIMESTAMP":  Timestamp,
	"TO":         To,
	"TSVECTOR":   Tsvector,
	"UPDATE":     Update,
	"USING":      Using,
	"VARYING":    Varying,
	"WITH":       With,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func LookupIdent(ident string) TokenType {
	if keyword, ok := keywords[strings.ToUpper(ident)]; ok {
		return keyword
	}
	return Identifier
}
