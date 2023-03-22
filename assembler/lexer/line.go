package lexer

import (
	"assembler/memory"
	"fmt"
	"strings"
)

type Line struct {
	Pos     Position       // Position in source file
	Line    string         // Content of line
	Address memory.Address // Address of line
	data    []byte         // Compiled opcode
	Label   string         // Assembly label
	Tokens  []*Token       // Tokenized line
	Comment string         // Comment for the line
}

func (l *Line) Data() []byte {
	if l != nil {
		return l.data
	}
	return nil
}

func (l *Line) SetData(d ...byte) int {
	if l != nil {
		l.data = d
		return len(d)
	}
	return 0
}

func (l *Line) String() string {
	if l.Label == "" && l.Line == "" {
		return fmt.Sprintf("%4s %s", "", l.Comment)
	}

	s, h := "", ""
	if l.Address > 0 {
		s = fmt.Sprintf("%04x", l.Address)
		for _, v := range l.data {
			h = h + fmt.Sprintf("%02x ", v)
		}
	} else {
		s = ""
	}
	a := []string{fmt.Sprintf("%4s %-11.11s %-8s %-32s %s", s, h, l.Label, l.Line, l.Comment)}
	for len(h) > 11 {
		h = h[12:]
		if len(h) > 0 {
			a = append(a, fmt.Sprintf("%4s %-11.11s", "", h))
		}
	}
	return strings.Join(a, "\n")
}
