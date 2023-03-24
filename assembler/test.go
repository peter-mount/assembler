package assembler

import (
	"assembler/assembler/context"
	"assembler/assembler/parser"
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"os"
	"strings"
	"testing"
)

type TestScript struct {
	Name     string           // Optional name of step
	Src      []string         // Text of lines to assemble
	Expected []*context.Block // Expected output from assembler
}

func RunTestScript(test string, t *testing.T, assembler *Assembler, scripts ...TestScript) {
	for sid, script := range scripts {
		n := test
		if script.Name != "" && !strings.HasPrefix(script.Name, "_") {
			n = n + "/" + script.Name
		}
		t.Run(fmt.Sprintf("%s_%d", n, sid), func(t *testing.T) {
			if err := runTestScript(t, assembler, script); err != nil {
				panic(err)
			}
		})
	}
}
func runTestScript(t *testing.T, assembler *Assembler, script TestScript) error {

	f, err := os.CreateTemp("", "test*")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	if _, err = f.Write([]byte(strings.Join(script.Src, "\n"))); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpName)

	if err := assembler.Assemble(tmpName); err != nil {
		return err
	}

	blocks := assembler.Blocks()

	lb, le := len(blocks), len(script.Expected)
	if lb != le {
		t.Errorf("got %d blocks, expected %d", lb, le)
	}

	for i, e := range script.Expected {
		if i < lb {
			bb, ob := blocks[i].Bytes(), e.Bytes()
			if len(bb) != len(ob) {
				t.Errorf("block %d length differs\ngot: %d\nexp: %d", i, len(bb), len(ob))
			}
			if !bytes.Equal(bb, ob) {
				t.Errorf("block %d differs\ngot: %x\nexp: %x", i, blocks[i].Bytes(), e.Bytes())
			}
		}
	}

	return nil
}

func NewAssembler(processors ...kernel.PostInitialisableService) (*Assembler, error) {
	for _, proc := range processors {
		if err := proc.PostInit(); err != nil {
			return nil, err
		}
	}

	assembler := &Assembler{ProcessorRegistry: &parser.ProcessorRegistry{}}
	return assembler, nil
}
