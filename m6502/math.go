package m6502

import (
	"github.com/peter-mount/assembler/assembler/node"
)

var (
	adcOpcodes = map[AddressMode]byte{
		AMImmediate:                     0x69,
		AMAddress:                       0x6d,
		AMAddressLong:                   0x6f,
		AMZeroPage:                      0x65,
		AMZeroPageIndirect:              0x72,
		AMZeroPageIndirectLong:          0x67,
		AMAbsoluteIndexedX:              0x7d,
		AMAbsoluteLongIndexedX:          0x7f,
		AMAbsoluteIndexedY:              0x79,
		AMZeroPageIndexedX:              0x75,
		AMZeroPageIndirectX:             0x61,
		AMZeroPageIndirectY:             0x71,
		AMZeroPageIndirectLongY:         0x77,
		AMStackRelative:                 0x63,
		AMStackRelativeIndirectIndexedY: 0x73,
	}

	andOpcodes = map[AddressMode]byte{
		AMImmediate:            0x29,
		AMAddress:              0x2d,
		AMAddressLong:          0x2f,
		AMZeroPage:             0x25,
		AMZeroPageIndirect:     0x32,
		AMZeroPageIndirectLong: 0x27,
		AMAbsoluteIndexedX:     0x3d,
		// long indexed x 0x3f,
		AMAbsoluteIndexedY: 0x39,
	}

	decOpcodes = map[AddressMode]byte{
		AMAccumulator:      0x3a,
		AMAddress:          0xce,
		AMZeroPage:         0xc6,
		AMAbsoluteIndexedX: 0xde,
		AMZeroPageIndexedX: 0xd6,
	}

	incOpcodes = map[AddressMode]byte{
		AMAccumulator:      0x1a,
		AMAddress:          0xee,
		AMZeroPage:         0xe6,
		AMAbsoluteIndexedX: 0xfe,
		AMZeroPageIndexedX: 0xf6,
	}

	sbcOpcodes = map[AddressMode]byte{
		AMImmediate:            0xe9,
		AMAddress:              0xed,
		AMAddressLong:          0xef,
		AMZeroPage:             0xe5,
		AMZeroPageIndirect:     0xe2,
		AMZeroPageIndirectLong: 0xe7,
		AMAbsoluteIndexedX:     0xfd,
		AMAbsoluteIndexedY:     0xf9,
		AMZeroPageIndexedX:     0xf5,
		AMZeroPageIndirectX:    0xe1,
		AMZeroPageIndirectY:    0xf1,
	}
)

func adc(addressModes ...AddressMode) node.Handler {
	return instruction(adcOpcodes, addressModes)
}

func and(addressModes ...AddressMode) node.Handler {
	return instruction(andOpcodes, addressModes)
}

func dec(addressModes ...AddressMode) node.Handler {
	return instruction(decOpcodes, addressModes)
}

func inc(addressModes ...AddressMode) node.Handler {
	return instruction(incOpcodes, addressModes)
}

func sbc(addressModes ...AddressMode) node.Handler {
	return instruction(sbcOpcodes, addressModes)
}
