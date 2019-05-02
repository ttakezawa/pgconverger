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
			name: "create table",
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
			want:    `ALTER TABLE "public"."x" ADD COLUMN "n" "text";`,
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
			want:    `ALTER TABLE "public"."x" ALTER COLUMN "n" TYPE "text";`,
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
ALTER TABLE "public"."x" ALTER COLUMN "n" TYPE "bigint";
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
