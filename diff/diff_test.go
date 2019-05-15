package diff

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"
)

type reader struct {
	io.Reader
	name string
}

func (r *reader) Name() string {
	return r.name
}

func newReader(data string) fileReader {
	return &reader{bytes.NewBufferString(data), "test"}
}

var reg = regexp.MustCompile(`\s+`)

func canonical(s string) string {
	return reg.ReplaceAllString(strings.TrimSpace(s), " ")
}

func TestProcess(t *testing.T) {
	type args struct {
		source  fileReader
		desired fileReader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "create table with default schema",
			args: args{
				source:  newReader(``),
				desired: newReader(`CREATE TABLE "x" ( id bigint );`),
			},
			want: `
CREATE TABLE "public"."x" (
    "id" bigint
);`,
			wantErr: false,
		},
		{
			name: "create table with explicit schema",
			args: args{
				source:  newReader(``),
				desired: newReader(`CREATE TABLE "myschema"."x" ( id bigint );`),
			},
			want: `
CREATE TABLE "myschema"."x" (
    "id" bigint
);`,
			wantErr: false,
		},
		{
			name: "set search_path and create table",
			args: args{
				source: newReader(``),
				desired: newReader(`
SET search_path = "myschema", pg_catalog;
CREATE TABLE "x" ( id bigint );`),
			},
			want: `
CREATE TABLE "myschema"."x" (
    "id" bigint
);`,
			wantErr: false,
		},
		{
			name: "create table with primary key",
			args: args{
				source: newReader(``),
				desired: newReader(`
CREATE TABLE users (id bigint);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);`),
			},
			want: `
CREATE TABLE "public"."users" (
    "id" bigint
);
ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");`,
			wantErr: false,
		},
		{
			name: "create table with unique",
			args: args{
				source: newReader(``),
				desired: newReader(`
CREATE TABLE users (id bigint);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_pkey UNIQUE (id);`),
			},
			want: `
CREATE TABLE "public"."users" (
    "id" bigint
);
ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_pkey" UNIQUE ("id");`,
			wantErr: false,
		},
		{
			name: "create table with set default",
			args: args{
				source: newReader(``),
				desired: newReader(`
CREATE TABLE users (id bigint);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);`),
			},
			want: `
CREATE TABLE "public"."users" (
    "id" bigint
);
ALTER TABLE ONLY "public"."users" ALTER COLUMN "id" SET DEFAULT "nextval"('public.users_id_seq'::"regclass");`,
			wantErr: false,
		},
		{
			name: "drop table",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint );`),
				desired: newReader(``),
			},
			want:    `DROP TABLE "public"."x";`,
			wantErr: false,
		},
		{
			name: "add column",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint );`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n text );`),
			},
			want:    `ALTER TABLE "public"."x" ADD COLUMN "n" text;`,
			wantErr: false,
		},
		{
			name: "drop column",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n text );`),
				desired: newReader(`CREATE TABLE "x" ( id bigint );`),
			},
			want:    `ALTER TABLE "public"."x" DROP COLUMN "n";`,
			wantErr: false,
		},
		{
			name: "alter column type datatype",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n bigint );`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n text );`),
			},
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" TYPE text;`,
			wantErr: false,
		},
		{
			name: "alter column set not null",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n bigint );`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n bigint NOT NULL);`),
			},
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" SET NOT NULL;`,
			wantErr: false,
		},
		{
			name: "alter column drop not null",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n bigint NOT NULL);`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n bigint);`),
			},
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" DROP NOT NULL;`,
			wantErr: false,
		},
		{
			name: "alter column set default",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n bigint);`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n bigint DEFAULT 42);`),
			},
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" SET DEFAULT 42;`,
			wantErr: false,
		},
		{
			name: "alter column drop default",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n bigint DEFAULT 42);`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n bigint);`),
			},
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" DROP DEFAULT;`,
			wantErr: false,
		},
		{
			name: "alter column type,not null",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n integer);`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n bigint DEFAULT 42);`),
			},
			want: `
ALTER TABLE "public"."x" ALTER COLUMN "n" TYPE bigint;
ALTER TABLE "public"."x" ALTER COLUMN "n" SET DEFAULT 42;`,
			wantErr: false,
		},
		{
			name: "alter column not null,drop default",
			args: args{
				source:  newReader(`CREATE TABLE "x" ( id bigint, n integer default 42);`),
				desired: newReader(`CREATE TABLE "x" ( id bigint, n integer not null);`),
			},
			want: `
ALTER TABLE "public"."x" ALTER COLUMN "n" SET NOT NULL;
ALTER TABLE "public"."x" ALTER COLUMN "n" DROP DEFAULT;`,
			wantErr: false,
		},
		{
			name: "create index",
			args: args{
				source: newReader(`
CREATE TABLE users (name text);`),
				desired: newReader(`
CREATE TABLE users (name text);
CREATE INDEX idx ON users USING btree (name);`),
			},
			want:    `CREATE INDEX "idx" ON "public"."users" USING "btree" ("name");`,
			wantErr: false,
		},
		{
			name: "drop index",
			args: args{
				source: newReader(`
CREATE TABLE users (name text);
CREATE INDEX idx ON users USING btree (name);`),
				desired: newReader(`
CREATE TABLE users (name text);`),
			},
			want:    `DROP INDEX "idx";`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Process(tt.args.source, tt.args.desired)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				type detailError interface{ Detail() string }
				if e, ok := err.(detailError); ok {
					t.Logf("detail: %s", e.Detail())
				}
				return
			}
			if canonical(got) != canonical(tt.want) {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
