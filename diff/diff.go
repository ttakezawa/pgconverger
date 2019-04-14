package diff

import (
	"fmt"
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

func Process(source fileReader, desired fileReader) (error, *Diff) {
	df := &Diff{
		source:  source,
		desired: desired,
	}
	df.sourceErrors, df.sourceDDL = parseOneSide(df.source)
	df.desiredErrors, df.desiredDDL = parseOneSide(df.desired)
	return df.ErrorOrNil(), df
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

type Error struct {
	diff *Diff
}

func (err *Error) Error() string {
	var (
		df     = err.diff
		errors []string
	)
	if df.sourceErrors != nil {
		errors = append(errors, fmt.Sprintf("source has %d errors", len(df.sourceErrors)))
	}
	if df.desiredErrors != nil {
		errors = append(errors, fmt.Sprintf("desired has %d errors", len(df.desiredErrors)))
	}
	return strings.Join(errors, "; ")
}

func (err *Error) Detail() string {
	var (
		df    = err.diff
		lines []string
	)
	if df.sourceErrors != nil {
		lines = append(lines, fmt.Sprintf("%s has %d errors", df.source.Name(), len(df.sourceErrors)))
		for _, e := range df.sourceErrors {
			lines = append(lines, fmt.Sprintf("  %s", e.Error()))
		}
	}
	if df.desiredErrors != nil {
		lines = append(lines, fmt.Sprintf("%s has %d errors", df.desired.Name(), len(df.desiredErrors)))
		for _, e := range df.desiredErrors {
			lines = append(lines, fmt.Sprintf("  %s", e.Error()))
		}
	}
	return strings.Join(lines, "\n")
}

func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, err.Detail())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, err.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", err.Error())
	}
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
