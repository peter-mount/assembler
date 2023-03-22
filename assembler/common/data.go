package common

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
)

func DataBlock(b ...uint8) node.Handler {
	return func(n *node.Node, ctx context.Context) error {
		switch ctx.GetStage() {

		case context.StageCompile:
			n.GetLine().SetData(b...)
			ctx.AddAddress(len(b))
		}
		return node.CallChildren(n, ctx)
	}
}
