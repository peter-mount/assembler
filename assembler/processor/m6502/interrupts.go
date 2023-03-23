package m6502

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
)

// BRK Software Break.
//
// Although the instruction is 1 byte long the Program Counter is incremented
// by 2 to allow for an optional parameter.
//
// To support how older assemblers work (which only write 1 byte) we do
// the following:
//
// BRK			1 byte instruction. It's up to the author to ensure that the following byte(s) is set.
// BRK #0x00	2 bytes with the value included as the second byte.
//
// For example, on the BBC Micro traditionally a break is written as:
//
//	BRK
//	EQUB 0 			; Error code
//	EQUS "Silly"	; Error message
//	EQUB 0			; End of message marker
//
// As an alternate:
//
//	BRK #0 			; Error code
//	EQUS "Silly"	; Error message
//	EQUB 0			; End of message marker
func BRK(n *node.Node, ctx context.Context) error {
	var err error
	switch ctx.GetStage() {

	case context.StageCompile:
		// Presume we are using 2 bytes initially, we will reduce to 1 in StageOptimise
		n.GetLine().SetData(0, 0)

	// StageOptimise to see which form we are using
	case context.StageOptimise, context.StageBackref:
		params, err := GetAddressing(n, ctx, AMImplied, AMZeroPage, AMImmediate)
		if err != nil {
			return n.Token.Pos.Error(err)
		}

		if params.AddressMode == AMImplied {
			n.GetLine().SetData(0x00)
		} else {
			n.GetLine().SetData(0x00, byte(params.Value&0xff))
		}
	}

	return err
}

// COP 65816 coprocessor instruction
func COP(n *node.Node, ctx context.Context) error {
	var err error
	switch ctx.GetStage() {

	case context.StageCompile:
		n.GetLine().SetData(0, 0)

	case context.StageBackref:
		// COP const - but we'll accept COP #value as well
		params, err := GetAddressing(n, ctx, AMImplied, AMZeroPage, AMImmediate)
		if err != nil {
			return n.Token.Pos.Error(err)
		}

		n.GetLine().SetData(0x02, byte(params.Value&0xff))
	}

	return err
}
