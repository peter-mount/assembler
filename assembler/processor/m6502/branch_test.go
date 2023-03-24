package m6502

import (
	"assembler/assembler"
	"assembler/assembler/context"
	"testing"
)

func TestJSR(t *testing.T) {

	asm, err := assembler.NewAssembler(&M6502{})
	if err != nil {
		t.Fatal(err)
	}

	// common result for the jsr tests
	cpu := "    CPU \"6502\""
	org := "    ORG 0x2000"
	expected := []*context.Block{context.NewBlock(0x2000, 0x20, 0xee, 0xff)}

	assembler.RunTestScript("6502/TestJSR", t, asm,
		// Test jsr to OSWRCH on the BBC using all 3 possible ways of writing 0xffee
		assembler.TestScript{
			Src:      []string{cpu, org, "    jsr 0xffee"},
			Expected: expected,
		},
		assembler.TestScript{
			Src:      []string{cpu, org, "    JSR &ffee"},
			Expected: expected,
		},
		assembler.TestScript{
			Src:      []string{cpu, org, "    jsr $ffee"},
			Expected: expected,
		},
	)
}

func TestBranch(t *testing.T) {

	asm, err := assembler.NewAssembler(&M6502{})
	if err != nil {
		t.Fatal(err)
	}

	// common result for the jsr tests

	var script []assembler.TestScript

	for op, opcode := range branchOpcodes {
		script = append(script, assembler.TestScript{
			Name: op,
			Src: []string{
				"    CPU \"6502\"",
				"    ORG 0x2000",
				".start LDY #42",
				".l1 LDA text,Y",
				"    " + op + " l2",
				"    JSR &ffee",
				"    INY",
				"    BNE l1",
				"l2  RTS",
				".text EQUS \"Hello world!\", 13, 10, 0",
			},
			Expected: []*context.Block{context.NewBlock(0x2000,
				0xa0, 0x2a,
				0xb9, 0x0e, 0x20,
				opcode, 0x08,
				0x20, 0xee, 0xff,
				0xc8,
				0xd0, 0xf7,
				0x60,
				'H', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!', 13, 10, 0,
				//12,
			)},
		})
	}
	assembler.RunTestScript("6502/TestBranch", t, asm, script...)
}
