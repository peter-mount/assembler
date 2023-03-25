package assembler

import (
	"flag"
	"fmt"
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/assembler/parser"
	"github.com/peter-mount/assembler/processor"
	"github.com/peter-mount/go-kernel/v2/log"
	"os"
	"strings"
	"time"
)

// Service handles the actual assembly of one or more source projects
type Service struct {
	ProcessorRegistry *processor.Registry `kernel:"inject"`
	ShowAssembly      *bool               `kernel:"flag,show-assembly,Show assembly"`
	ShowSymbols       *bool               `kernel:"flag,show-symbols,Show symbols"`
	ShowTimings       *bool               `kernel:"flag,show-stage-timings,Show stage timings"`
}

// Assembler holds the state during a single project's assembly
type Assembler struct {
	ProcessorRegistry *processor.Registry
	lexer             *lexer.Lexer
	parser            *parser.Parser
	root              *node.Node
	ShowAssembly      bool
	ShowSymbols       bool
	ShowTimings       bool
	fileName          string // Filename used in normal operation
	lines             []*lexer.Line
	blocks            []*context.Block // All blocks
}

func (a *Service) Start() error {
	for _, fileName := range flag.Args() {
		asm := a.Assembler()
		asm.fileName = fileName
		if err := asm.assemble(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return nil
		}
	}
	return nil
}

func (a *Service) Assembler() *Assembler {
	return &Assembler{
		ShowAssembly: *a.ShowAssembly,
		ShowSymbols:  *a.ShowSymbols,
		ShowTimings:  *a.ShowTimings,
	}
}

func (a *Assembler) Assemble(fileName string) error {
	a.fileName = fileName
	return a.assemble()
}

func (a *Assembler) Blocks() []*context.Block {
	return a.blocks
}

func (a *Assembler) assemble() error {
	now1 := time.Now()
	ctx := context.New()

	err := ctx.ForEachStage(a.processStage)
	if err != nil {
		return err
	}

	log.Printf("Assembly took %v", time.Now().Sub(now1))
	return err
}

func (a *Assembler) processStage(stage context.Stage, ctx context.Context) error {
	if a.ShowTimings {
		now2 := time.Now()
		defer func() {
			log.Printf("Stage %d took %v", stage, time.Now().Sub(now2))
		}()
	}

	switch stage {
	case context.StageInit:
		a.lexer = nil
		a.lines = nil
		a.parser = nil
		a.root = nil
		a.blocks = nil

	case context.StageTokenize:
		a.lexer = &lexer.Lexer{}
		if err := a.lexer.Parse(a.fileName); err != nil {
			return err
		}
		a.lines = a.lexer.Lines()

	case context.StageParse:
		a.parser = &parser.Parser{ProcessorRegistry: a.ProcessorRegistry}
		if a.lexer != nil {
			root, err := a.parser.Parse(a.lines)
			if err != nil {
				return err
			}
			a.root = root
		}

	case context.StageList:
		if a.ShowAssembly {
			return a.listSources(ctx)
		}

	case context.StageSymbols:
		if a.ShowSymbols {
			return a.listSymbols(ctx)
		}

	case context.StageAssemble:
		ctx.ClearBlocks()
		if err := a.root.Visit(ctx); err != nil {
			return err
		}
		a.blocks = ctx.GetAllBlocks()

	default:
		return a.root.Visit(ctx)
	}

	return nil
}

func (a *Assembler) listSources(ctx context.Context) error {
	if err := a.root.Visit(ctx); err != nil {
		return err
	}
	// Write the last address at the end, this is the address after the
	// previous content but handy in a listing
	log.Printf("%04x", ctx.GetAddress())
	return nil
}

func (a *Assembler) listSymbols(ctx context.Context) error {
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
}
