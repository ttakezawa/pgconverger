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
	Bigint     TokenType = "Bigint"
	By         TokenType = "By"
	Cache      TokenType = "Cache"
	Character  TokenType = "Character"
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
)

var keywords = map[string]TokenType{
	"ADD":        Add,
	"ALTER":      Alter,
	"BIGINT":     Bigint,
	"BY":         By,
	"CACHE":      Cache,
	"CHARACTER":  Character,
	"COLUMN":     Column,
	"CONSTRAINT": Constraint,
	"CREATE":     Create,
	"DEFAULT":    Default,
	"GRANT":      Grant,
	"INCREMENT":  Increment,
	"INDEX":      Index,
	"INSERT":     Insert,
	"MAXVALUE":   Maxvalue,
	"MINVALUE":   Minvalue,
	"NO":         No,
	"NOT":        Not,
	"NULL":       Null,
	"ON":         On,
	"ONLY":       Only,
	"OWNER":      Owner,
	"SELECT":     Select,
	"SEQUENCE":   Sequence,
	"SET":        Set,
	"START":      Start,
	"TABLE":      Table,
	"TO":         To,
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
