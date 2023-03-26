package m6502

import (
	"github.com/peter-mount/assembler/assembler/node"
)

type ldOpcode struct {
	name     string                // Instruction name
	operands map[AddressMode]uint8 // Map of AddressMode -> OpCode
}

var (
	ldaOpcodes = map[AddressMode]byte{
		AMImmediate:                     0xA9,
		AMAddress:                       0xAD,
		AMAddressLong:                   0xAF,
		AMZeroPage:                      0xA5,
		AMZeroPageIndirect:              0xB2,
		AMZeroPageIndirectLong:          0xA7,
		AMAbsoluteIndexedX:              0xBD,
		AMAbsoluteLongIndexedX:          0xBF,
		AMAbsoluteIndexedY:              0xB9,
		AMZeroPageIndexedX:              0xB5,
		AMZeroPageIndirectX:             0xA1,
		AMZeroPageIndirectY:             0xB1,
		AMZeroPageIndirectLongY:         0xB7,
		AMStackRelative:                 0xA3,
		AMStackRelativeIndirectIndexedY: 0xB3,
	}
	ldxOpcodes = map[AddressMode]byte{
		AMImmediate:        0xA2,
		AMAddress:          0xAE,
		AMZeroPage:         0xA6,
		AMAbsoluteIndexedY: 0xBE,
		AMZeroPageIndexedY: 0xB6,
	}
	ldyOpcodes = map[AddressMode]byte{
		AMImmediate:        0xA0,
		AMAddress:          0xAC,
		AMZeroPage:         0xA4,
		AMAbsoluteIndexedX: 0xBC,
		AMZeroPageIndexedX: 0xB4,
	}
)

func lda(addressModes ...AddressMode) node.Handler {
	return instruction(ldaOpcodes, addressModes)
}

func ldx() node.Handler {
	return instruction(ldxOpcodes, []AddressMode{AMImmediate, AMAddress, AMZeroPage, AMAbsoluteIndexedY, AMZeroPageIndexedY})
}

func ldy() node.Handler {
	return instruction(ldyOpcodes, []AddressMode{AMImmediate, AMAddress, AMZeroPage, AMAbsoluteIndexedX, AMZeroPageIndexedX})
}
