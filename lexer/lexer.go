package lexer

import (
	"unicode"
	"unicode/utf8"
)

type tokenType string

const (
	tokenError tokenType = "error"
	tokenEOF   tokenType = "eof"
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
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

func (l *lexer) advance() {
	if l.readPosition >= len(l.input) {
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

func lexText(l *lexer) stateFn {
	peek := l.peekChar()
	switch {
	case isAlphaNumeric(peek):
		return lexIdentifier
	default:
		// TODO: IMPLEMENT
		return nil
	}
}

func lexIdentifier(l *lexer) stateFn {
	return nil // TODO: IMPLEMENT
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
