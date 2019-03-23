package lexer

import (
	"fmt"
	"testing"

	"github.com/ttakezawa/pgconverger/token"
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
		typ     token.TokenType
		literal string
		line    int
	}

	tests := []struct {
		input string
		wants []want
	}{
		{
			input: `abc
x f1 abc$2 "foobar"
'aaaa''bbbb' 'a\'b' 'ab
cd'
1.5 .82
-- aaaa
/* xxx /* aaaa */
yyy */
`,
			wants: []want{
				{token.Identifier, "abc", 1},
				{token.Identifier, "x", 2},
				{token.Identifier, "f1", 2},
				{token.Identifier, "abc$2", 2},
				{token.Identifier, `"foobar"`, 2},
				{token.String, "'aaaa''bbbb'", 3},
				{token.String, "'a\\'b'", 3},
				{token.String, "'ab\ncd'", 3},
				{token.Number, "1.5", 5},
				{token.Number, ".82", 5},
				{token.EOF, "", 9},
			},
		},
		{
			input: `'a`,
			wants: []want{
				{token.Illegal, "'a", 1},
			},
		},
		{
			input: `
CREATE TABLE "users" (
    "id" bigint NOT NULL,
    name character varying(50)
);

ALTER TABLE "users" OWNER TO "api";

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY "users" ALTER COLUMN "id" SET DEFAULT "nextval"('"users_id_seq"'::"regclass");
`,
			wants: []want{
				{token.Create, "CREATE", 2},
				{token.Table, "TABLE", 2},
				{token.Identifier, `"users"`, 2},
				{token.LParen, "(", 2},
				{token.Identifier, `"id"`, 3},
				{token.Bigint, "bigint", 3},
				{token.Not, "NOT", 3},
				{token.Null, "NULL", 3},
				{token.Comma, ",", 3},
				{token.Identifier, "name", 4},
				{token.Character, "character", 4},
				{token.Varying, "varying", 4},
				{token.LParen, "(", 4},
				{token.Number, "50", 4},
				{token.RParen, ")", 4},
				{token.RParen, ")", 5},
				{token.Semicolon, ";", 5},

				{token.Alter, "ALTER", 7},
				{token.Table, "TABLE", 7},
				{token.Identifier, `"users"`, 7},
				{token.Owner, "OWNER", 7},
				{token.To, "TO", 7},
				{token.Identifier, `"api"`, 7},
				{token.Semicolon, ";", 7},

				{token.Create, "CREATE", 9},
				{token.Sequence, "SEQUENCE", 9},
				{token.Identifier, "users_id_seq", 9},
				{token.Start, "START", 10},
				{token.With, "WITH", 10},
				{token.Number, "1", 10},
				{token.Increment, "INCREMENT", 11},
				{token.By, "BY", 11},
				{token.Number, "1", 11},
				{token.No, "NO", 12},
				{token.Minvalue, "MINVALUE", 12},
				{token.No, "NO", 13},
				{token.Maxvalue, "MAXVALUE", 13},
				{token.Cache, "CACHE", 14},
				{token.Number, "1", 14},
				{token.Semicolon, ";", 14},

				{token.Alter, "ALTER", 16},
				{token.Table, "TABLE", 16},
				{token.Only, "ONLY", 16},
				{token.Identifier, `"users"`, 16},
				{token.Alter, "ALTER", 16},
				{token.Column, "COLUMN", 16},
				{token.Identifier, `"id"`, 16},
				{token.Set, "SET", 16},
				{token.Default, "DEFAULT", 16},
				{token.Identifier, `"nextval"`, 16},
				{token.LParen, "(", 16},
				{token.String, "'\"users_id_seq\"'", 16},
				{token.Typecast, "::", 16},
				{token.Identifier, `"regclass"`, 16},
				{token.RParen, ")", 16},
				{token.Semicolon, ";", 16},

				{token.EOF, "", 17},
			},
		},
	}

	for i, tt := range tests {
		l := Lex(tt.input)
		for j, want := range tt.wants {
			got := l.NextToken()
			if got.Type != want.typ {
				t.Errorf("case%d-%d Lexer.NextToken().typ = %v, want %v", i+1, j+i, got.Type, want.typ)
			}
			if got.Literal != want.literal {
				t.Errorf("case%d-%d Lexer.NextToken().literal = %q, want %q", i+1, j+i, got.Literal, want.literal)
			}
			if got.Line != want.line {
				t.Errorf("case%d-%d Lexer.NextToken().line = %v, want %v", i+1, j+i, got.Line, want.line)
			}
		}
	}
}
