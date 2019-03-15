package lexer

import (
	"unicode"
	"unicode/utf8"
)

type tokenType string

const (
	Illegal    tokenType = "Illegal"
	EOF        tokenType = "EOF"
	Space      tokenType = "Space"
	Comment    tokenType = "Comment"
	Identifier tokenType = "Identifier"
	String     tokenType = "String"
	Number     tokenType = "Number"
)

type token struct {
	typ  tokenType
	val  string
	line int
}

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	char         rune

	startPosition int
	startLine     int

	tokens chan token
	state  stateFn
}

func newLexer(input string) *lexer {
	return &lexer{
		input:     input,
		tokens:    make(chan token),
		line:      1,
		startLine: 1,
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

func (l *lexer) emit(typ tokenType) {
	switch typ {
	case Space, Comment:
		// ignore Space and Comment
	default:
		l.tokens <- token{
			typ:  typ,
			val:  l.input[l.startPosition:l.position],
			line: l.startLine,
		}
	}
	l.startPosition = l.position
	l.startLine = l.line
}

func (l *lexer) NextToken() token {
	return <-l.tokens
}

func lexFn(l *lexer) stateFn {
	switch {
	case isSpace(l.char):
		return lexSpace
	case isIdentifierStart(l.char):
		return lexIdentifier
	case isNumberStart(l.char):
		return lexNumber
	case l.char == '\'':
		return lexString
	case l.char == '-' && l.peekChar() == '-':
		return lexCommentLine
	case l.char == '/' && l.peekChar() == '*':
		return lexCommentBlock
	case l.char == eof:
		return lexEOF
	default:
		return lexIllegal
	}
}

func lexEOF(l *lexer) stateFn {
	l.emit(EOF)
	return nil
}

func lexIllegal(l *lexer) stateFn {
	l.emit(Illegal)
	return nil
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.char) {
		l.advance()
	}
	l.emit(Space)
	return lexFn
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// ident_start		[A-Za-z\200-\377_]
// ident_cont		[A-Za-z\200-\377_0-9\$]
// identifier		{ident_start}{ident_cont}*
func lexIdentifier(l *lexer) stateFn {
	l.advance()
	for isIdentifierCont(l.char) {
		l.advance()
	}
	l.emit(Identifier)
	return lexFn
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierCont(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) || r == '$'
}

// 'Dianne''s horse' => "Dianne's horse"
// Not implemented: 'foo'\n'bar' => 'foobar'
// Not implemented: $$Dianne's horse$$
// Not implemented: $SomeTag$Dianne's horse$SomeTag$
// Not implemented: E'foo' (String Constants With C-Style Escapes)
// Not implemented: U&'d\0061t\+000061' (String Constants With Unicode Escapes)
func lexString(l *lexer) stateFn {
	l.advance()
Loop:
	for {
		switch l.char {
		case eof:
			l.emit(Illegal)
			return nil
		case '\'':
			if l.peekChar() == '\'' {
				// doubled quote.
				l.advance()
				l.advance()
				continue Loop
			}
			l.advance()
			break Loop
		case '\\':
			if l.peekChar() == '\'' {
				// escapted quote.
				l.advance()
			}
			l.advance()
			continue Loop
		}
		l.advance()
	}
	l.emit(String)
	return lexFn
}

// 42
// 3.5
// 4.
// .001
// Not implemented: 5e2
// Not implemented: 1.925e-3
func lexNumber(l *lexer) stateFn {
Loop:
	for {
		l.advance()
		switch {
		case unicode.IsDigit(l.char):
			continue
		case l.char == '.':
			return lexNumberFraction
		default:
			break Loop
		}
	}
	l.emit(Number)
	return lexFn
}

func lexNumberFraction(l *lexer) stateFn {
Loop:
	for {
		l.advance()
		switch {
		case unicode.IsDigit(l.char):
			continue
		default:
			break Loop
		}
	}
	l.emit(Number)
	return lexFn
}

func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}

// -- This is a standard SQL comment
func lexCommentLine(l *lexer) stateFn {
	l.advance()
	l.advance()
Loop:
	for l.char != eof {
		if l.char == '\n' {
			l.advance()
			break Loop
		}
		l.advance()
	}
	l.emit(Comment)
	return lexFn
}

// /* multiline comment
//  * with nesting: /* nested block comment */
//  */
func lexCommentBlock(l *lexer) stateFn {
	l.advance()
	l.advance()
	nestCount := 0
Loop:
	for {
		switch {
		case l.char == eof:
			return lexIllegal
		case l.char == '/' && l.peekChar() == '*':
			nestCount++
			l.advance()
		case l.char == '*' && l.peekChar() == '/':
			if nestCount == 0 {
				l.advance()
				l.advance()
				break Loop
			}
			nestCount--
			l.advance()
			l.advance()
		default:
			l.advance()
		}
	}
	l.emit(Comment)
	return lexFn
}
