package lexer

import (
	"fmt"
	"reflect"
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
	input := `abc
x f1 abc$2
`
	l := Lex(input)

	tests := []struct {
		wantTyp  tokenType
		wantVal  string
		wantLine int
	}{
		{Identifier, "abc", 1},
		{Identifier, "x", 2},
		{Identifier, "f1", 2},
		{Identifier, "abc$2", 2},
	}
	for i, tt := range tests {
		got := l.NextToken()
		if got.typ != tt.wantTyp {
			t.Errorf("case%d lexer.NextToken().typ = %v, want %v", i+1, got.typ, tt.wantTyp)
		}
		if got.val != tt.wantVal {
			t.Errorf("case%d lexer.NextToken().val = %v, want %v", i+1, got.val, tt.wantVal)
		}
		if got.line != tt.wantLine {
			t.Errorf("case%d lexer.NextToken().line = %v, want %v", i+1, got.line, tt.wantLine)
		}
	}
	var (
		got  = l.NextToken()
		want = token{}
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("lexer.NextToken() = %#v, want %#v", got, want)
	}
}
