package m6502

import (
	"github.com/peter-mount/assembler/processor"
)

func init() {
	processor.Register(M6502(), M65C02(), M65816())
}

// M6502 implements the 6502 processor.
// Subtypes like the 65c02 extends this for the additional instructions
func M6502() processor.Processor {
	return processor.New("6502").
		Include(addBranchOpcodes).
		Include(addRegisterInstructions(ldOpcodes6502)).
		Handle("ADC", adc(AMImmediate, AMAddress, AMZeroPage, AMAbsoluteIndexedX, AMAbsoluteIndexedY, AMZeroPageIndirectX, AMZeroPageIndirectY)).
		Handle("BRK", brk).
		Handle("DEC", dec(AMAddress, AMZeroPage, AMAbsoluteIndexedX, AMZeroPageIndexedX)).
		Handle("INC", inc(AMAddress, AMZeroPage, AMAbsoluteIndexedX, AMZeroPageIndexedX)).
		Simple("INX", 0xe8).
		Simple("INY", 0xc8).
		Handle("JSR", jsr(AMAddress)).
		Simple("NOP", 0xea).
		Simple("RTI", 0x40).
		Simple("RTS", 0x60).
		Handle("SBC", sbc(AMImmediate, AMAddress, AMZeroPage, AMAbsoluteIndexedX, AMAbsoluteIndexedY, AMZeroPageIndirectX, AMZeroPageIndirectY)).
		Build()
}

// M65C02 implements the 65C02 processor which extends the 6502 instruction set
func M65C02() processor.Processor {
	return processor.New("65c02").
		Extends(M6502()).
		Handle("BRA", branchAlways).
		Handle("DEC", dec(AMAccumulator)).
		Handle("INC", inc(AMAccumulator)).
		Handle("SBC", sbc(AMZeroPageIndirect)).
		Build()
}

// M65816 implements the 65816 instruction set which extends the 65C02 with 8-bit & 16-bit instructions
func M65816() processor.Processor {
	return processor.New("65816").
		Extends(M65C02()).
		Handle("ADC", adc(AMAddressLong)).
		Handle("COP", cop).
		Handle("JSR", jsr(AMAddressLong, AMAbsoluteIndexedIndirect)).
		Handle("SBC", sbc(AMAddressLong)).
		Simple("STP", 0xdb).
		Simple("WAI", 0xcb).
		Simple("WDM", 0x42).
		Simple("XCE", 0xfb).
		Build()
}
