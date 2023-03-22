package assembler

import (
	"assembler/assembler/lexer"
	"assembler/assembler/parser"
	"assembler/memory"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"os"
)

type Assembler struct {
	processorRegistry *parser.ProcessorRegistry `kernel:"inject"`
	memory            *memory.Map
}

func (a *Assembler) Start() error {
	for _, fileName := range flag.Args() {
		if err := a.Assemble(fileName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return nil
		}

	}
	return nil
}

func (a *Assembler) Assemble(fileName string) error {
	lex := lexer.Lexer{}

	err := lex.Parse(fileName)
	if err != nil {
		return err
	}

	parse := parser.Parser{
		ProcessorRegistry: a.processorRegistry,
	}
	_, err = parse.Parse(lex.Lines())
	if err != nil {
		return err
	}

	_ = lex.ForEach(func(l *lexer.Line) error {
		log.Println(l.String())
		return nil
	})

	return err
}
