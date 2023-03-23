package m6502

type Register uint8

// The 6502 family registers, not all are used here but here for reference
const (
	RegNone = iota // Holder for when no register has been found
	RegA           // A Accumulator
	RegX           // X Index
	RegY           // Y Index
	RegP           // Processor status flags
	RegS           // Stack pointer
	RegPC          // Program Counter
	RegB           // B 8-bit Accumulator 65816 only hidden but exchangeable with A
	RegC           // C 16-bit Accumulator 65816 alias for A
	RegDBR         // Data Bank Register 65816
	RegPBR         // Program Bank Register 65816
)
