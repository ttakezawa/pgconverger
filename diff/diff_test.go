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

func newReader(name string, data string) fileReader {
	return &reader{bytes.NewBufferString(data), name}
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
				source:  newReader("testsource", ``),
				desired: newReader("testdesired", `CREATE TABLE "x" ( id bigint );`),
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
				source:  newReader("testsource", `CREATE TABLE "x" ( id bigint );`),
				desired: newReader("testdesired", ``),
			},
			want:    `DROP TABLE "public"."x";`,
			wantErr: false,
		},
		{
			name: "add column",
			args: args{
				source:  newReader("testsource", `CREATE TABLE "x" ( id bigint );`),
				desired: newReader("testdesired", `CREATE TABLE "x" ( id bigint, name text );`),
			},
			want:    `ALTER TABLE "public"."x" ADD COLUMN "name" "text";`,
			wantErr: false,
		},
		{
			name: "drop column",
			args: args{
				source:  newReader("testsource", `CREATE TABLE "x" ( id bigint, name text );`),
				desired: newReader("testdesired", `CREATE TABLE "x" ( id bigint );`),
			},
			want:    `ALTER TABLE "public"."x" DROP COLUMN "name";`,
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
