package m6502

import (
	"github.com/peter-mount/assembler/assembler/node"
)

func adc(addressModes ...AddressMode) node.Handler {
	return instruction(map[AddressMode]byte{
		AMImmediate:            0x69,
		AMAddress:              0x6d,
		AMAddressLong:          0x6f,
		AMZeroPage:             0x65,
		AMZeroPageIndirect:     0x72,
		AMZeroPageIndirectLong: 0x67,
	}, addressModes)
}

func sbc(addressModes ...AddressMode) node.Handler {
	return instruction(map[AddressMode]byte{
		AMImmediate:            0xe9,
		AMAddress:              0xed,
		AMAddressLong:          0xef,
		AMZeroPage:             0xe5,
		AMZeroPageIndirect:     0xe2,
		AMZeroPageIndirectLong: 0xe7,
	}, addressModes)
}

func dec(addressModes ...AddressMode) node.Handler {
	return instruction(map[AddressMode]byte{
		AMAccumulator:      0x3a,
		AMAddress:          0xce,
		AMZeroPage:         0xc6,
		AMAbsoluteIndexedX: 0xde,
		AMZeroPageIndexedX: 0xd6,
	}, addressModes)
}

func inc(addressModes ...AddressMode) node.Handler {
	return instruction(map[AddressMode]byte{
		AMAccumulator:      0x1a,
		AMAddress:          0xee,
		AMZeroPage:         0xe6,
		AMAbsoluteIndexedX: 0xfe,
		AMZeroPageIndexedX: 0xf6,
	}, addressModes)
}
