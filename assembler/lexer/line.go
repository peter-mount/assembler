package lexer

import (
	"assembler/memory"
	"fmt"
)

type Line struct {
	Pos     Position       // Position in source file
	Line    string         // Content of line
	Address memory.Address // Address of line
	Length  int            // Opcode length
	Label   string         // Assembly label
	Tokens  []*Token       // Tokenized line
	Comment string         // Comment for the line
}

func (l *Line) String(mem *memory.Map) string {
	if l.Label == "" && l.Line == "" {
		return fmt.Sprintf("%4s %s", "", l.Comment)
	}

	s, h := "", ""
	if l.Address > 0 {
		s = fmt.Sprintf("%04x", l.Address)
		for i := 0; i < l.Length; i++ {
			v, err := mem.ReadByte(l.Address)
			if err != nil {
				h = h + "**"
			} else {
				h = h + fmt.Sprintf("%02x", v)
			}
		}
	} else {
		s = ""
	}
	return fmt.Sprintf("%4s %-8.8s %-8s %-32s %s", s, h, l.Label, l.Line, l.Comment)
}
