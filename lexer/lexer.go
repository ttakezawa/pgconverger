package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/ttakezawa/pgconverger/token"
)

const eof = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	char         rune

	startPosition int
	startLine     int

	tokens chan token.Token
	state  stateFn
}

func newLexer(input string) *Lexer {
	return &Lexer{
		input:     input,
		tokens:    make(chan token.Token),
		line:      1,
		startLine: 1,
	}
}

func Lex(input string) *Lexer {
	l := newLexer(input)
	l.advance()
	go l.run()
	return l
}

func (l *Lexer) run() {
	for l.state = lexFn; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

func (l *Lexer) advance() {
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

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) word() string {
	i := l.startPosition
	j := l.position
	if i > len(l.input) {
		i = len(l.input)
	}
	if j > len(l.input) {
		j = len(l.input)
	}
	return l.input[i:j]
}

func (l *Lexer) emit(typ token.TokenType) {
	switch typ {
	case token.Space, token.CommentLine, token.CommentBlock:
		// ignore Space and Comment
	default:
		l.tokens <- token.Token{
			Type:    typ,
			Literal: l.word(),
			Line:    l.startLine,
		}
	}
	l.startPosition = l.position
	l.startLine = l.line
}

func (l *Lexer) NextToken() token.Token {
	return <-l.tokens
}

func lexFn(l *Lexer) stateFn {
	switch {
	case isSpace(l.char):
		return lexSpace
	case l.char == '"':
		return lexDoubleQuoteIdentifier
	case isIdentifierStart(l.char):
		return lexIdentifier
	case isNumberStart(l.char):
		return lexNumber
	case l.char == '\'':
		return lexString
	case l.char == ':':
		return lexTypecast
	case l.char == '+':
		l.advance()
		l.emit(token.Plus)
	case l.char == '*':
		l.advance()
		l.emit(token.Asterisk)
	case l.char == '/':
		l.advance()
		l.emit(token.Slash)
	case l.char == ';':
		l.advance()
		l.emit(token.Semicolon)
	case l.char == ',':
		l.advance()
		l.emit(token.Comma)
	case l.char == '(':
		l.advance()
		l.emit(token.LParen)
	case l.char == ')':
		l.advance()
		l.emit(token.RParen)
	case l.char == '-' && l.peekChar() == '-':
		return lexCommentLine
	case l.char == '/' && l.peekChar() == '*':
		return lexCommentBlock
	case l.char == eof:
		return lexEOF
	default:
		return lexIllegal
	}

	return lexFn
}

func lexEOF(l *Lexer) stateFn {
	l.emit(token.EOF)
	return nil
}

func lexIllegal(l *Lexer) stateFn {
	l.advance()
	l.emit(token.Illegal)
	return lexFn
}

func lexSpace(l *Lexer) stateFn {
	for isSpace(l.char) {
		l.advance()
	}
	l.emit(token.Space)
	return lexFn
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// "foobar"
func lexDoubleQuoteIdentifier(l *Lexer) stateFn {
	l.advance()
Loop:
	for {
		switch l.char {
		case eof:
			return lexIllegal
		case '"':
			l.advance()
			break Loop

		}
		l.advance()
	}
	l.emit(token.Identifier)
	return lexFn
}

// ident_start		[A-Za-z\200-\377_]
// ident_cont		[A-Za-z\200-\377_0-9\$]
// identifier		{ident_start}{ident_cont}*
func lexIdentifier(l *Lexer) stateFn {
	l.advance()
	for isIdentifierCont(l.char) {
		l.advance()
	}
	l.emit(token.LookupIdent(l.word()))
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
func lexString(l *Lexer) stateFn {
	l.advance()
Loop:
	for {
		switch l.char {
		case eof:
			return lexIllegal
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
	l.emit(token.String)
	return lexFn
}

// 42
// 3.5
// 4.
// .001
// Not implemented: 5e2
// Not implemented: 1.925e-3
func lexNumber(l *Lexer) stateFn {
Loop:
	for {
		l.advance()
		switch {
		case unicode.IsDigit(l.char):
			continue
		case l.char == '.':
			return lexMantissa
		default:
			break Loop
		}
	}
	l.emit(token.Number)
	return lexFn
}

func lexMantissa(l *Lexer) stateFn {
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
	l.emit(token.Number)
	return lexFn
}

func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}

// -- This is a standard SQL comment
func lexCommentLine(l *Lexer) stateFn {
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
	l.emit(token.CommentLine)
	return lexFn
}

// /* multiline comment
//  * with nesting: /* nested block comment */
//  */
func lexCommentBlock(l *Lexer) stateFn {
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
	l.emit(token.CommentBlock)
	return lexFn
}

func lexTypecast(l *Lexer) stateFn {
	l.advance()
	if l.char != ':' {
		l.advance()
		return lexIllegal
	}
	l.advance()
	l.emit(token.Typecast)
	return lexFn
}
