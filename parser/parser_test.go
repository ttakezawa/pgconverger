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
		input    string
		expected string
	}{
		{`CREATE TABLE "users" (
    "id" bigint NOT NULL,
    "id2" bigint NOT NULL DEFAULT 'nextval(''users_id_seq''::regclass)',
    bigint bigint,
    name character varying(50) DEFAULT '-' NOT NULL,
    name2 "text",
    data jsonb,
    data2 "jsonb",
    "digit" integer,
    "bytes" bytea,
    is_read boolean,
    "numer" numeric,
    "date" date,
    "vec" "tsvector",
    "x" bigint NULL,
    created_at timestamp with time zone,
    "code" "text" DEFAULT '0001'::"text" NOT NULL,
    expr bigint DEFAULT 1+2*3/4::"text",
    flag1 bigint DEFAULT TRUE,
    flag2 bigint DEFAULT falSe
);`,
			`CREATE TABLE "users" (
    "id" bigint NOT NULL,
    "id2" bigint NOT NULL DEFAULT 'nextval(''users_id_seq''::regclass)',
    "bigint" bigint,
    "name" character varying(50) DEFAULT '-' NOT NULL,
    "name2" text,
    "data" jsonb,
    "data2" jsonb,
    "digit" integer,
    "bytes" bytea,
    "is_read" boolean,
    "numer" numeric,
    "date" date,
    "vec" tsvector,
    "x" bigint,
    "created_at" timestamp with time zone,
    "code" text DEFAULT '0001'::"text" NOT NULL,
    "expr" bigint DEFAULT 1+2*3/4::"text",
    "flag1" bigint DEFAULT TRUE,
    "flag2" bigint DEFAULT FALSE
);
`,
		},
	}

	for i, tt := range tests {
		p := New(lexer.Lex(tt.input))
		dataDefinition := p.ParseDataDefinition()
		var builder strings.Builder
		dataDefinition.WriteStringTo(&builder)
		if builder.String() != tt.expected {
			t.Errorf("case%d:\n\tgot  =      %q,\n\twant =      %q", i+1, builder.String(), tt.expected)
		}
		checkParserErrors(t, p)
	}
}

func TestCreateIndexStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`CREATE INDEX "users_name_key" ON "users" USING "btree" ("name");`,
			`CREATE INDEX "users_name_key" ON "users" USING "btree" ("name");`,
		},
		{
			`CREATE UNIQUE INDEX "users_email_phone_key" ON "users" USING "btree" ("email", "phone");`,
			`CREATE UNIQUE INDEX "users_email_phone_key" ON "users" USING "btree" ("email", "phone");`,
		},
		{
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" (("deleted_at" IS NULL), "name");`,
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" (("deleted_at" IS NULL), "name");`,
		},
		{
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" ((1*2*(3+4)), "name");`,
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" ((1*2*(3+4)), "name");`,
		},
		{
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" ((("deleted_at" IS NULL)), "name");`,
			`CREATE INDEX "users_deleted_at_and_name_key" ON "users" USING "btree" ((("deleted_at" IS NULL)), "name");`,
		},
	}

	for i, tt := range tests {
		p := New(lexer.Lex(tt.input))
		dataDefinition := p.ParseDataDefinition()
		var builder strings.Builder
		dataDefinition.WriteStringTo(&builder)
		if builder.String() != tt.expected {
			t.Errorf("case%d:\n\tgot  =      %q,\n\twant =      %q", i+1, builder.String(), tt.expected)
		}
		checkParserErrors(t, p)
	}
}
