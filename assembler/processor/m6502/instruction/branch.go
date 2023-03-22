package instruction

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
)

func JSR(n *node.Node, ctx context.Context) error {
	switch ctx.GetStage() {

	case context.StageCompile:
		n.GetLine().SetData(0x20, 0, 0)
		ctx.AddAddress(3)
	}

	return node.CallChildren(n, ctx)
}
