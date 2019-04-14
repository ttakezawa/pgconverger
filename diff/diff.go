package diff

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/ttakezawa/pgconverger/ast"
	"github.com/ttakezawa/pgconverger/lexer"
	"github.com/ttakezawa/pgconverger/parser"
)

type fileReader interface {
	io.Reader
	Name() string
}

type Diff struct {
	source       fileReader
	sourceDDL    *ast.DataDefinition
	sourceErrors []error

	desired       fileReader
	desiredDDL    *ast.DataDefinition
	desiredErrors []error
}

func Process(source fileReader, desired fileReader) (*Diff, error) {
	df := &Diff{
		source:  source,
		desired: desired,
	}
	df.sourceErrors, df.sourceDDL = parseOneSide(df.source)
	df.desiredErrors, df.desiredDDL = parseOneSide(df.desired)
	return df, df.ErrorOrNil()
}

func parseOneSide(reader fileReader) (errs []error, ddl *ast.DataDefinition) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		errs = append(errs, err)
		return
	}
	p := parser.New(lexer.Lex(reader.Name(), string(input)))
	ddl = p.ParseDataDefinition()
	errs = append(errs, p.Errors()...)
	return
}

func (df *Diff) ErrorOrNil() error {
	if df.sourceErrors != nil {
		return &Error{df}
	}
	if df.desiredErrors != nil {
		return &Error{df}
	}
	return nil
}

func (df *Diff) String() string {
	var builder strings.Builder
	df.desiredDDL.WriteStringTo(&builder)
	return builder.String()
}
