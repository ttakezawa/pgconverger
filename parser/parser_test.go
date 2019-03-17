package parser

import (
	"testing"

	"github.com/ttakezawa/pgconverger/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	if len(p.errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range p.errors {
		t.Errorf("  %s", msg)
	}
	t.FailNow()
}

func TestCreateTableStatement(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`
CREATE TABLE "users" (
-- "id" bigint NOT NULL,
-- name character varying(50)
);`,
		},
	}

	for _, tt := range tests {
		p := New(lexer.Lex(tt.input))
		dataDefinition := p.ParseDataDefinition()

		t.Logf("DDL: %+v", dataDefinition)
		checkParserErrors(t, p)
	}
}
