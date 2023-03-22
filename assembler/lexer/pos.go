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
	return fmt.Errorf(p.String()+" "+s, a...)
}

func (p Position) Error(err error) error {
	return p.Errorf("%s", err.Error())
}
