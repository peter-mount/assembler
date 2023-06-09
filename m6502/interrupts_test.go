package m6502

import (
	"github.com/peter-mount/assembler/assembler"
	"github.com/peter-mount/assembler/assembler/context"
	"testing"
)

func TestBRK(t *testing.T) {

	asm, err := assembler.NewAssembler(M6502())
	if err != nil {
		t.Fatal(err)
	}

	// common result for the jsr tests
	cpu := "    CPU \"6502\""
	org := "    ORG 0x2000"

	assembler.RunTestScript("6502/brk", t, asm,
		// Standard brk. Most assemblers create just 1 byte with the next byte being added
		// by the programmer
		assembler.TestScript{
			Src:      []string{cpu, org, "    brk"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x00)},
		},
		// We support 2 extensions which assemble as 2 bytes as the PC is incremented by 2 when brk is executed
		assembler.TestScript{
			Src:      []string{cpu, org, "    brk 42"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x00, 42)},
		},
		assembler.TestScript{
			Src:      []string{cpu, org, "    brk #81"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x00, 81)},
		},
	)
}

func TestCOP(t *testing.T) {

	asm, err := assembler.NewAssembler(M65816())
	if err != nil {
		t.Fatal(err)
	}

	// common result for the jsr tests
	cpu := "    CPU \"65816\""
	org := "    ORG 0x2000"

	assembler.RunTestScript("6502/cop", t, asm,
		assembler.TestScript{
			Src:      []string{cpu, org, "    cop 42"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x02, 42)},
		},
		assembler.TestScript{
			Src:      []string{cpu, org, "    cop #81"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x02, 81)},
		},
	)
}
