package token

import (
	"strings"
)

type TokenType string

const (
	Illegal    TokenType = "Illegal"
	EOF        TokenType = "EOF"
	Space      TokenType = "Space"
	Comment    TokenType = "Comment"
	Identifier TokenType = "Identifier"
	String     TokenType = "String"
	Number     TokenType = "Number"
	Semicolon  TokenType = "Semicolon"
	Comma      TokenType = "Comma"
	LParen     TokenType = "LParen"
	RParen     TokenType = "RParen"
	Typecast   TokenType = "Typecast"

	Add        TokenType = "Add"
	Alter      TokenType = "Alter"
	By         TokenType = "By"
	Cache      TokenType = "Cache"
	Column     TokenType = "Column"
	Constraint TokenType = "Constraint"
	Create     TokenType = "Create"
	Default    TokenType = "Default"
	Grant      TokenType = "Grant"
	Increment  TokenType = "Increment"
	Index      TokenType = "Index"
	Insert     TokenType = "Insert"
	Maxvalue   TokenType = "Maxvalue"
	Minvalue   TokenType = "Minvalue"
	No         TokenType = "No"
	Not        TokenType = "Not"
	Null       TokenType = "Null"
	On         TokenType = "On"
	Only       TokenType = "Only"
	Owner      TokenType = "Owner"
	Select     TokenType = "Select"
	Sequence   TokenType = "Sequence"
	Set        TokenType = "Set"
	Start      TokenType = "Start"
	Table      TokenType = "Table"
	To         TokenType = "To"
	Update     TokenType = "Update"
	Using      TokenType = "Using"
	Varying    TokenType = "Varying"
	With       TokenType = "With"

	Bigint    TokenType = "Bigint"
	Bigserial TokenType = "Bigserial"
	Boolean   TokenType = "Boolean"
	Bytea     TokenType = "Bytea"
	Character TokenType = "Character"
	Date      TokenType = "Date"
	Integer   TokenType = "Integer"
	Jsonb     TokenType = "Jsonb"
	Numeric   TokenType = "Numeric"
	Serial    TokenType = "Serial"
	Text      TokenType = "Text"
	Timestamp TokenType = "Timestamp"
	Time      TokenType = "Time"
	Tsvector  TokenType = "Tsvector"
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
