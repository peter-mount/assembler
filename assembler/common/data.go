package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/node"
)

func DataBlock(b ...uint8) node.Handler {
	return func(n *node.Node, ctx context.Context) error {
		switch ctx.GetStage() {

		case context.StageCompile:
			n.GetLine().SetData(b...)
		}
		return node.CallChildren(n, ctx)
	}
}
