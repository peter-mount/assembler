package lexer

import (
	"fmt"
)

type Position struct {
	File *File // Link to source file
	Line int   // Line number in file
	Pos  int   // Character position in file
}

func (p Position) String() string {
	if p.Pos < 0 {
		return fmt.Sprintf("[%s:%d]", p.File.BaseName(), p.Line)
	}
	return fmt.Sprintf("[%s,%d,%d]", p.File.BaseName(), p.Line, p.Pos)
}

func (p Position) Errorf(s string, a ...interface{}) error {
	return &errorString{s: fmt.Sprintf(p.String()+" "+s, a...)}
}

func (p Position) Error(err error) error {
	if err1, ok := err.(*errorString); ok {
		return err1
	}
	return &errorString{s: fmt.Sprintf("%s %s", p.String(), err.Error())}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
