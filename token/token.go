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
	Dot
	Semicolon
	Comma
	LParen
	RParen
	LBracket
	RBracket
	Equal
	Plus
	Minus
	Asterisk
	Slash
	Typecast

	Add
	Alter
	Asc
	BackslashConnect
	By
	Cache
	Column
	Concurrently
	Constraint
	Create
	Database
	Default
	Desc
	Exists
	Extension
	False
	Function
	Grant
	If
	Increment
	Index
	Insert
	Is
	Key
	Maxvalue
	Minvalue
	No
	Not
	Null
	On
	Only
	Operator
	Owned
	Owner
	Primary
	Revoke
	Role
	Schema
	Select
	Sequence
	Set
	Start
	Table
	To
	Trigger
	True
	Unique
	Update
	Using
	Varying
	View
	With
	Without
	Zone

	Bigint
	Smallint
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
	Uuid
)

type keyword struct {
	Type     TokenType
	Reserved bool
}

// https://www.postgresql.org/docs/10/sql-keywords-appendix.html
var keywords = map[string]keyword{
	"ADD":          {Add, false},
	"ALTER":        {Alter, false},
	"ASC":          {Asc, true},
	"\\CONNECT":    {BackslashConnect, false},
	"BIGINT":       {Bigint, false},
	"BIGSERIAL":    {Bigserial, false},
	"BOOLEAN":      {Boolean, false},
	"BY":           {By, false},
	"BYTEA":        {Bytea, false},
	"CACHE":        {Cache, false},
	"CHARACTER":    {Character, false},
	"COLUMN":       {Column, true},
	"COMMENT":      {Comment, false},
	"CONCURRENTLY": {Concurrently, true}, // reserved (can be function or type)
	"CONSTRAINT":   {Constraint, true},
	"CREATE":       {Create, true},
	"DATABASE":     {Database, false},
	"DATE":         {Date, false},
	"DEFAULT":      {Default, true},
	"DESC":         {Desc, true},
	"EXISTS":       {Exists, false}, // non-reserved (cannot be function or type)
	"EXTENSION":    {Extension, false},
	"FALSE":        {False, true},
	"FUNCTION":     {Function, false},
	"GRANT":        {Grant, true},
	"IF":           {If, false},
	"INCREMENT":    {Increment, false},
	"INDEX":        {Index, false},
	"INSERT":       {Insert, false},
	"INTEGER":      {Integer, false},
	"IS":           {Is, true}, // reserved (can be function or type)
	"JSONB":        {Jsonb, false},
	"KEY":          {Key, false},
	"MAXVALUE":     {Maxvalue, false},
	"MINVALUE":     {Minvalue, false},
	"NO":           {No, false},
	"NOT":          {Not, true},
	"NULL":         {Null, true},
	"NUMERIC":      {Numeric, false},
	"ON":           {On, true},
	"ONLY":         {Only, true},
	"OPERATOR":     {Operator, false},
	"OWNED":        {Owned, false},
	"OWNER":        {Owner, false},
	"PRIMARY":      {Primary, true},
	"REVOKE":       {Revoke, false},
	"ROLE":         {Role, false},
	"SCHEMA":       {Schema, false},
	"SELECT":       {Select, true},
	"SEQUENCE":     {Sequence, false},
	"SERIAL":       {Serial, false},
	"SET":          {Set, false},
	"SMALLINT":     {Smallint, false},
	"START":        {Start, false},
	"TABLE":        {Table, true},
	"TEXT":         {Text, false},
	"TIME":         {Time, false},
	"TIMESTAMP":    {Timestamp, false},
	"TO":           {To, true},
	"TRIGGER":      {Trigger, false},
	"TRUE":         {True, true},
	"TSVECTOR":     {Tsvector, false},
	"UNIQUE":       {Unique, true},
	"UPDATE":       {Update, false},
	"USING":        {Using, true},
	"UUID":         {Uuid, false},
	"VARYING":      {Varying, false},
	"VIEW":         {View, false},
	"WITH":         {With, true},
	"WITHOUT":      {Without, true},
	"ZONE":         {Zone, false},
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func (tok *Token) IsKeyword() bool {
	key := strings.ToUpper(tok.Type.String())
	_, ok := keywords[key]
	return ok
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
