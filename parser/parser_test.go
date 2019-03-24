package parser

import (
	"testing"

	"github.com/ttakezawa/pgconverger/ast"
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
  name character varying(50) DEFAULT '-' NOT NULL
  -- created_at timestamp with time zone
);`,
		},
	}

	for _, tt := range tests {
		p := New(lexer.Lex(tt.input))
		dataDefinition := p.ParseDataDefinition()

		t.Logf("DDL: %+v", dataDefinition)
		t.Logf("DDL: %#v", dataDefinition)

		createTableStmt, _ := dataDefinition.StatementList[0].(*ast.CreateTableStatement)
		t.Logf("CREATE_TABLE: %+v", createTableStmt)
		t.Logf("CREATE_TABLE: %#v", createTableStmt)
		t.Logf("COL1: %#v", createTableStmt.ColumnDefinitionList[0])
		t.Logf("COL2: %#v", createTableStmt.ColumnDefinitionList[1])
		t.Logf("COL3: %#v", createTableStmt.ColumnDefinitionList[2])

		checkParserErrors(t, p)
	}
}
