package processor

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
)

// SimpleInstruction is used for op codes that have just 1 single instance.
// e.g. RTS, INX or INY on the 6502 family
func SimpleInstruction(opCode ...uint8) node.Handler {
	return func(n *node.Node, ctx context.Context) error {
		switch ctx.GetStage() {
		case context.StageCompile:
			n.GetLine().SetData(opCode...)
		}
		return node.CallChildren(n, ctx)
	}
}
