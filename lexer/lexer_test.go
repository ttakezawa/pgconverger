package lexer

import (
	"fmt"
	"testing"
)

func Test_lexer_readChar(t *testing.T) {
	type expect struct {
		char rune
		peek rune
		line int
	}

	tests := []struct {
		input   string
		expects []expect
	}{
		{
			input: "a b\ne",
			expects: []expect{
				{'a', ' ', 1},
				{' ', 'b', 1},
				{'b', '\n', 1},
				{'\n', 'e', 2},
				{'e', eof, 2},
				{eof, eof, 2},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i+1), func(t *testing.T) {
			l := newLexer(tt.input)
			for _, expect := range tt.expects {
				l.advance()
				if expect.char != l.char {
					t.Fatalf("expected char=%q, got=%q", expect.char, l.char)
				}
				peek := l.peekChar()
				if expect.peek != peek {
					t.Fatalf("expected peek=%q, got=%q", expect.peek, peek)
				}
				if expect.line != l.line {
					t.Fatalf("expected line=%d, got=%d", expect.line, l.line)
				}
			}
		})
	}
}

func Test_lexer_NextToken(t *testing.T) {
	type want struct {
		typ  tokenType
		val  string
		line int
	}

	tests := []struct {
		input string
		wants []want
	}{
		{
			input: `abc
x f1 abc$2
'aaaa''bbbb' 'a\'b' 'ab
cd'
1.5 .82
-- aaaa
/* xxx /* aaaa */
yyy */
`,
			wants: []want{
				{Identifier, "abc", 1},
				{Identifier, "x", 2},
				{Identifier, "f1", 2},
				{Identifier, "abc$2", 2},
				{String, "'aaaa''bbbb'", 3},
				{String, "'a\\'b'", 3},
				{String, "'ab\ncd'", 3},
				{Number, "1.5", 5},
				{Number, ".82", 5},
				{EOF, "", 9},
			},
		},
		{
			input: `'a`,
			wants: []want{
				{Illegal, "'a", 1},
			},
		},
	}
	for i, tt := range tests {
		l := Lex(tt.input)
		for j, want := range tt.wants {
			got := l.NextToken()
			if got.typ != want.typ {
				t.Errorf("case%d-%d lexer.NextToken().typ = %v, want %v", i+1, j+i, got.typ, want.typ)
			}
			if got.val != want.val {
				t.Errorf("case%d-%d lexer.NextToken().val = %v, want %v", i+1, j+i, got.val, want.val)
			}
			if got.line != want.line {
				t.Errorf("case%d-%d lexer.NextToken().line = %v, want %v", i+1, j+i, got.line, want.line)
			}
		}
	}
}
