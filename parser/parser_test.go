package parser

import (
	"strings"
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
  "id" bigint NOT NULL,
  "id2" bigint NOT NULL DEFAULT 'nextval(''users_id_seq''::regclass)',
  name character varying(50) DEFAULT '-' NOT NULL,
  created_at timestamp with time zone
);`,
		},
	}

	for _, tt := range tests {
		p := New(lexer.Lex(tt.input))
		dataDefinition := p.ParseDataDefinition()
		var builder strings.Builder
		dataDefinition.Source(&builder)
		t.Logf("\n%s", builder.String())
		checkParserErrors(t, p)
	}
}
