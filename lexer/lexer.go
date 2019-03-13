package lexer

import (
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	Illegal tokenType = iota
	EOF
	Create
	Identifier
)

type token struct {
	typ  tokenType
	val  string
	line int
	col  int
}

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	char         rune
	tokens       chan token
	state        stateFn
}

func newLexer(input string) *lexer {
	return &lexer{
		input:  input,
		tokens: make(chan token),
		line:   1,
	}
}

func Lex(input string) *lexer {
	l := newLexer(input)
	l.advance()
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexFn; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

func (l *lexer) advance() {
	if l.readPosition >= len(l.input) {
		l.position++
		l.readPosition++
		l.char = eof
		return
	}
	r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.position = l.readPosition
	l.readPosition += size
	l.char = r
	if r == '\n' {
		l.line++
	}
}

func (l *lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *lexer) emit(t token) {
	l.tokens <- t
}

func (l *lexer) NextToken() token {
	return <-l.tokens
}

func lexFn(l *lexer) stateFn {
	switch {
	case isSpace(l.char):
		return lexSpace
	case isAlphaNumeric(l.char):
		return lexIdentifier
	case l.char == eof:
		return nil
	default:
		return lexIllegal
	}
}

func lexIllegal(l *lexer) stateFn {
	l.emit(token{
		typ: Illegal,
		val: string(l.char),
	})
	return nil
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.char) {
		l.advance()
	}
	return lexFn
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func lexIdentifier(l *lexer) stateFn {
	begin := l.position
	line := l.line
	for isAlphaNumeric(l.char) {
		l.advance()
	}
	l.emit(token{
		typ:  Identifier,
		val:  l.input[begin:l.position],
		line: line,
	})
	return lexFn
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
