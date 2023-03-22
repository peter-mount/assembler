package m6502

import "github.com/peter-mount/go-kernel/v2"

func init() {
	kernel.Register(
		&M6502{},
		&M65c02{},
	)
}

const (
	_                    = -(iota + 1000) // Start tokens at -1000 so we don't clash with internal ones
	Implied                               // Implied: TAY
	Stack                                 // Stack: PHA TODO do we need this and not Implied?
	Accumulator                           // ASL A
	ZeroPage                              // Zero Page/Direct Page: LDA $12
	Absolute                              // Absolute: LDA $1234
	Immediate                             // Immediate: LDA #$12
	ZeroPageIndirect                      // 65c02: LDA ($12)
	AbsoluteLong                          // 65802/65816: LDA $123456
	ZeroPageIndirectLong                  // 65802/65816: LDA [$12]
	BlockMove                             // 65802/65816: MVN source, dest
)
