package m6502

const (
	NOP = -(iota + 1000) // Start tokens at -1000 so we don't clash with internal ones
	BNE
	JSR
	LDA
	LDX
	LDY
	INX
	INY
	RTS
)
