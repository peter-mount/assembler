package assembler

import (
	"assembler/assembler/lexer"
	"assembler/memory"
	"flag"
	"github.com/peter-mount/go-kernel/v2/log"
)

type Assembler struct {
	memory *memory.Map
	lexer  *lexer.Lexer
}

func (a *Assembler) Start() error {
	for _, fileName := range flag.Args() {
		if err := a.Assemble(fileName); err != nil {
			return err
		}

	}
	return nil
}

func (a *Assembler) Assemble(fileName string) error {
	if a.lexer == nil {
		a.lexer = &lexer.Lexer{}
	}

	err := a.lexer.Parse(fileName)
	if err != nil {
		return err
	}

	_ = a.lexer.ForEach(func(l *lexer.Line) error {
		log.Println(l.String(a.memory))
		return nil
	})

	return err
}
