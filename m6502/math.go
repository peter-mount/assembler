package m6502

import (
	"github.com/peter-mount/assembler/assembler/node"
)

var (
	adcOpcodes = map[AddressMode]byte{
		AMImmediate:            0x69,
		AMAddress:              0x6d,
		AMAddressLong:          0x6f,
		AMZeroPage:             0x65,
		AMZeroPageIndirect:     0x72,
		AMZeroPageIndirectLong: 0x67,
		AMAbsoluteIndexedX:     0x7d,
		AMAbsoluteIndexedY:     0x79,
		AMZeroPageIndexedX:     0x75,
		AMZeroPageIndirectX:    0x61,
		AMZeroPageIndirectY:    0x71,
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
)

func adc(addressModes ...AddressMode) node.Handler {
	return instruction(adcOpcodes, addressModes)
}

func sbc(addressModes ...AddressMode) node.Handler {
	return instruction(sbcOpcodes, addressModes)
}

func dec(addressModes ...AddressMode) node.Handler {
	return instruction(decOpcodes, addressModes)
}

func inc(addressModes ...AddressMode) node.Handler {
	return instruction(incOpcodes, addressModes)
}
