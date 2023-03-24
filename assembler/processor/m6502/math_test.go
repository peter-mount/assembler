package m6502

import (
	"assembler/assembler"
	"assembler/assembler/context"
	"assembler/assembler/errors"
	"testing"
)

func TestADC(t *testing.T) {

	asm, err := assembler.NewAssembler(&M6502{}, &M65c02{}, &M65816{})
	if err != nil {
		t.Fatal(err)
	}

	assembler.RunTestScript("6502/ADC", t, asm,
		assembler.TestScript{
			Name:     AMImmediate.String(),
			Src:      []string{" cpu \"6502\"", " org 0x2000", " adc #12"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x69, 12)},
		},
		assembler.TestScript{
			Name:     AMAddress.String(),
			Src:      []string{" cpu \"6502\"", " org 0x2000", " adc 0x1234"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x6d, 0x34, 0x12)},
		},
		assembler.TestScript{
			Name:     AMAddress.String(),
			Src:      []string{" cpu \"6502\"", " org 0x2000", " adc 0x1234"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x6d, 0x34, 0x12)},
		},
		assembler.TestScript{
			Name:     AMZeroPage.String(),
			Src:      []string{" cpu \"6502\"", " org 0x2000", " adc 0x70"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x65, 0x70)},
		},
		assembler.TestScript{
			Name:     AMZeroPage.String(),
			Src:      []string{" cpu \"6502\"", " org 0x2000", " adc 0x70"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x65, 0x70)},
		},
		// 65816 only
		assembler.TestScript{
			Name:  AMAddressLong.String() + "_6502",
			Src:   []string{" cpu \"6502\"", " org 0x2000", " adc 0x123456"},
			Error: errors.IsUnsupportedError,
		},
		assembler.TestScript{
			Name:  AMAddressLong.String() + "_65c02",
			Src:   []string{" cpu \"65c02\"", " org 0x2000", " adc 0x123456"},
			Error: errors.IsUnsupportedError,
		},
		assembler.TestScript{
			Name:     AMAddressLong.String() + "_65816",
			Src:      []string{" cpu \"65816\"", " org 0x2000", " adc 0x123456"},
			Expected: []*context.Block{context.NewBlock(0x2000, 0x6f, 0x56, 0x34, 0x12)},
		},
	)
}
