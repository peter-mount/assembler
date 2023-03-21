package memory

// Address a generic memory address
type Address uint

// Address16 a 16 bit memory address, used for 8-bit machines
type Address16 uint16

// Address32 a 32 bit memory address, used for 16-bit and 32-bit machines
type Address32 uint32

// Address64 a 64 bit memory address used for 64-bit machines
type Address64 uint64

type Map struct {
	Blocks []Block
}

type Block struct {
}
