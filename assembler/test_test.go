package assembler

import (
	"assembler/assembler/context"
	"assembler/assembler/processor/m6502"
	"testing"
)

// TestRunTestScript use this not just as a test but a basis of actual tests using this utility
func TestRunTestScript(t *testing.T) {

	// Create an inline assembler
	asm, err := NewAssembler(&m6502.M6502{})
	if err != nil {
		t.Fatal(err)
	}

	// Run a single test script. Here we test the JSR 6502 opcode
	RunTestScript("example", t, asm,
		// Test jsr to OSWRCH on the BBC using all 3 possible ways of writing 0xffee
		TestScript{
			Src:      []string{" CPU \"6502\"", " ORG 0x200", " jsr 0xffee"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x20, 0xee, 0xff)},
		},
	)
}
