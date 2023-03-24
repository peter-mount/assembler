package errors

import (
	"errors"
	"fmt"
)

var (
	illegalArgument = errors.New("illegal argument")
	outOfBounds     = errors.New("out of bounds")
	stackEmpty      = errors.New("stack empty")
)

func IllegalArgument() error           { return illegalArgument }
func IsIllegalArgument(err error) bool { return err == illegalArgument }

func OutOfBounds() error           { return outOfBounds }
func IsOutOfBounds(err error) bool { return err == outOfBounds }

func StackEmpty() error           { return stackEmpty }
func IsStackEmpty(err error) bool { return err == stackEmpty }

type unsupportedError struct {
	uid uint64
	s   string
}

func (e *unsupportedError) Error() string {
	return e.s
}

func UnsupportedError(f string, a ...interface{}) error {
	return &unsupportedError{s: fmt.Sprintf("Unsupported: "+f, a...)}
}

func IsUnsupportedError(err error) bool {
	_, ok := err.(*unsupportedError)
	return ok
}
