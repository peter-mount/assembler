package assembler

import (
	"assembler/assembler/context"
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/assembler/parser"
	"assembler/memory"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"os"
	"strings"
)

type Assembler struct {
	processorRegistry *parser.ProcessorRegistry `kernel:"inject"`
	memory            *memory.Map
	lexer             *lexer.Lexer
	root              *node.Node
	fileName          string
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
	a.fileName = fileName

	ctx := context.New()

	err := ctx.ForEachStage(a.processStage)
	if err != nil {
		return err
	}

	return err
}

func (a *Assembler) processStage(stage context.Stage, ctx context.Context) error {
	switch stage {
	case context.StageLex:
		a.lexer = &lexer.Lexer{}
		return a.lexer.Parse(a.fileName)

	case context.StageParse:
		parse := parser.Parser{
			ProcessorRegistry: a.processorRegistry,
		}
		root, err := parse.Parse(a.lexer.Lines())
		if err != nil {
			return err
		}
		a.root = root
		return nil

	case context.StageList:
		if err := a.root.Visit(ctx); err != nil {
			return err
		}
		// Write the last address at the end, this is the address after the
		// previous content but handy in a listing
		log.Printf("%04x", ctx.GetAddress())
		return nil

	case context.StageSymbols:
		labels := ctx.GetLabels()
		ml := 8
		for _, l := range labels {
			le := len(l)
			if le > ml {
				ml = le
			}
		}
		f := fmt.Sprintf("%%-%d.%ds", ml, ml)
		s := fmt.Sprintf(f+" Address", "Label")
		log.Printf("%s\n%s", s, strings.Repeat("=", len(s)))
		f = f + " %x"
		for _, l := range labels {
			log.Printf(f, l, ctx.GetLabel(l).Address)
		}
		log.Println(strings.Repeat("=", len(s)))
		return nil

	default:
		return a.root.Visit(ctx)
	}
}
