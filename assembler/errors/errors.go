package errors

import "errors"

var (
	illegalArgument = errors.New("illegal argument")
	stackEmpty      = errors.New("stack empty")
)

func IllegalArgument() error { return illegalArgument }

func IsIllegalArgument(err error) bool { return err == illegalArgument }

func StackEmpty() error           { return stackEmpty }
func IsStackEmpty(err error) bool { return err == stackEmpty }
