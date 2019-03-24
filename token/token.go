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
	CommentBlock
	CommentLine
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
	False
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
	Schema
	Select
	Sequence
	Set
	Start
	Table
	To
	True
	Update
	Using
	Varying
	With
	Zone

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
	"COMMENT":    Comment,
	"CONSTRAINT": Constraint,
	"CREATE":     Create,
	"DATE":       Date,
	"DEFAULT":    Default,
	"FALSE":      False,
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
	"SCHEMA":     Schema,
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
	"TRUE":       True,
	"TSVECTOR":   Tsvector,
	"UPDATE":     Update,
	"USING":      Using,
	"VARYING":    Varying,
	"WITH":       With,
	"ZONE":       Zone,
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
