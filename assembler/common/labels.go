package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/node"
)

func LabelHandler(n *node.Node, ctx context.Context) error {
	if ctx.GetStage() == context.StageCompile {
		l := n.GetLine()
		if l != nil && l.Label != "" {
			ctx.SetLabel(l.Label, l)
		}
	}
	return node.CallChildren(n, ctx)
}
