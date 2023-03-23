package errors

import "errors"

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
