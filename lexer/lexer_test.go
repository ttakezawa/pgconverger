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
		{
			input: `
CREATE TABLE users (
    id bigint NOT NULL,
    name character varying(50)
);

ALTER TABLE users OWNER TO api;

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
`,
			wants: []want{
				{Identifier, "CREATE", 2},
				{Identifier, "TABLE", 2},
				{Identifier, "users", 2},
				{LParen, "(", 2},
				{Identifier, "id", 3},
				{Identifier, "bigint", 3},
				{Identifier, "NOT", 3},
				{Identifier, "NULL", 3},
				{Comma, ",", 3},
				{Identifier, "name", 4},
				{Identifier, "character", 4},
				{Identifier, "varying", 4},
				{LParen, "(", 4},
				{Number, "50", 4},
				{RParen, ")", 4},
				{RParen, ")", 5},
				{Semicolon, ";", 5},

				{Identifier, "ALTER", 7},
				{Identifier, "TABLE", 7},
				{Identifier, "users", 7},
				{Identifier, "OWNER", 7},
				{Identifier, "TO", 7},
				{Identifier, "api", 7},
				{Semicolon, ";", 7},

				{Identifier, "CREATE", 9},
				{Identifier, "SEQUENCE", 9},
				{Identifier, "users_id_seq", 9},
				{Identifier, "START", 10},
				{Identifier, "WITH", 10},
				{Number, "1", 10},
				{Identifier, "INCREMENT", 11},
				{Identifier, "BY", 11},
				{Number, "1", 11},
				{Identifier, "NO", 12},
				{Identifier, "MINVALUE", 12},
				{Identifier, "NO", 13},
				{Identifier, "MAXVALUE", 13},
				{Identifier, "CACHE", 14},
				{Number, "1", 14},
				{Semicolon, ";", 14},
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
				t.Errorf("case%d-%d lexer.NextToken().val = %q, want %q", i+1, j+i, got.val, want.val)
			}
			if got.line != want.line {
				t.Errorf("case%d-%d lexer.NextToken().line = %v, want %v", i+1, j+i, got.line, want.line)
			}
		}
	}
}
