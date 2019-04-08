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
	Plus
	Minus
	Asterisk
	Slash
	Typecast

	Add
	Alter
	By
	Cache
	Column
	Constraint
	Create
	Default
	Extension
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

type keyword struct {
	Type     TokenType
	Reserved bool
}

// https://www.postgresql.org/docs/10/sql-keywords-appendix.html
var keywords = map[string]keyword{
	"ADD":        {Add, false},
	"ALTER":      {Alter, false},
	"BIGINT":     {Bigint, false},
	"BIGSERIAL":  {Bigserial, false},
	"BOOLEAN":    {Boolean, false},
	"BY":         {By, false},
	"BYTEA":      {Bytea, false},
	"CACHE":      {Cache, false},
	"CHARACTER":  {Character, false},
	"COLUMN":     {Column, true},
	"COMMENT":    {Comment, false},
	"CONSTRAINT": {Constraint, true},
	"CREATE":     {Create, true},
	"DATE":       {Date, false},
	"DEFAULT":    {Default, true},
	"EXTENSION":  {Extension, false},
	"FALSE":      {False, true},
	"GRANT":      {Grant, true},
	"INCREMENT":  {Increment, false},
	"INDEX":      {Index, false},
	"INSERT":     {Insert, false},
	"INTEGER":    {Integer, false},
	"JSONB":      {Jsonb, false},
	"MAXVALUE":   {Maxvalue, false},
	"MINVALUE":   {Minvalue, false},
	"NO":         {No, false},
	"NOT":        {Not, true},
	"NULL":       {Null, true},
	"NUMERIC":    {Numeric, false},
	"ON":         {On, true},
	"ONLY":       {Only, true},
	"OWNER":      {Owner, false},
	"SCHEMA":     {Schema, false},
	"SELECT":     {Select, true},
	"SEQUENCE":   {Sequence, false},
	"SERIAL":     {Serial, false},
	"SET":        {Set, false},
	"START":      {Start, false},
	"TABLE":      {Table, true},
	"TEXT":       {Text, false},
	"TIME":       {Time, false},
	"TIMESTAMP":  {Timestamp, false},
	"TO":         {To, true},
	"TRUE":       {True, true},
	"TSVECTOR":   {Tsvector, false},
	"UPDATE":     {Update, false},
	"USING":      {Using, true},
	"VARYING":    {Varying, false},
	"WITH":       {With, true},
	"ZONE":       {Zone, false},
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func (tok *Token) IsReserved() bool {
	key := strings.ToUpper(tok.Type.String())
	keyword, ok := keywords[key]
	return ok && keyword.Reserved
}

func LookupIdent(ident string) TokenType {
	if keyword, ok := keywords[strings.ToUpper(ident)]; ok {
		return keyword.Type
	}
	return Identifier
}
