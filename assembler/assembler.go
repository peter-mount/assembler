package assembler

import (
	"assembler/assembler/context"
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/assembler/parser"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"os"
	"strings"
	"time"
)

// Assembler handles the actual assembly of one or more source projects
type Assembler struct {
	processorRegistry *parser.ProcessorRegistry `kernel:"inject"`
	showAssembly      *bool                     `kernel:"flag,show-assembly,Show assembly"`
	showSymbols       *bool                     `kernel:"flag,show-symbols,Show symbols"`
	showTimings       *bool                     `kernel:"flag,show-stage-timings,Show stage timings"`
}

// assembler holds the state during a single project's assembly
type assembler struct {
	lexer        *lexer.Lexer
	parser       *parser.Parser
	root         *node.Node
	showAssembly bool
	showSymbols  bool
	showTimings  bool
	fileName     string
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
	asm := &assembler{
		lexer: &lexer.Lexer{},
		parser: &parser.Parser{
			ProcessorRegistry: a.processorRegistry,
		},
		root:         nil,
		fileName:     fileName,
		showAssembly: *a.showAssembly,
		showSymbols:  *a.showSymbols,
		showTimings:  *a.showTimings,
	}

	now1 := time.Now()
	ctx := context.New()

	err := ctx.ForEachStage(asm.processStage)
	if err != nil {
		return err
	}

	log.Printf("Assembly took %v", time.Now().Sub(now1))
	return err
}

func (a *assembler) processStage(stage context.Stage, ctx context.Context) error {
	if a.showTimings {
		now2 := time.Now()
		defer func() {
			log.Printf("Stage %d took %v", stage, time.Now().Sub(now2))
		}()
	}

	switch stage {
	case context.StageLex:
		return a.lexer.Parse(a.fileName)

	case context.StageParse:
		root, err := a.parser.Parse(a.lexer.Lines())
		if err != nil {
			return err
		}
		a.root = root
		return nil

	case context.StageList:
		if a.showAssembly {
			return a.listSources(ctx)
		}

	case context.StageSymbols:
		if a.showSymbols {
			return a.listSymbols(ctx)
		}

	default:
		return a.root.Visit(ctx)
	}

	return nil
}

func (a *assembler) listSources(ctx context.Context) error {
	if err := a.root.Visit(ctx); err != nil {
		return err
	}
	// Write the last address at the end, this is the address after the
	// previous content but handy in a listing
	log.Printf("%04x", ctx.GetAddress())
	return nil
}

func (a *assembler) listSymbols(ctx context.Context) error {
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
