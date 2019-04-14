package diff

import (
	"fmt"
	"io"
	"strings"
)

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
