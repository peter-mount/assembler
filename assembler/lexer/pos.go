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
	return &Error{s: fmt.Sprintf(p.String()+" "+s, a...)}
}

func (p Position) Error(err error) error {
	if err1, ok := err.(*Error); ok {
		return err1
	}
	return &Error{s: fmt.Sprintf("%s %s", p.String(), err.Error()), e: err}
}

// Error is a trivial implementation of error.
type Error struct {
	s string
	e error
}

func (e *Error) Error() string {
	return e.s
}

func (e *Error) HasCause() bool {
	return e.e != nil
}

func (e *Error) Cause() error {
	return e.e
}

func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// GetCause gets the first cause that's not a positioned Error
func GetCause(err error) error {
	if err1, ok := err.(*Error); ok {
		return GetCause(err1.Cause())
	}
	return err
}
