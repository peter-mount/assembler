package node

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/util"
)

// IntHandler converts its text into an integer and pushes it onto the stack
func IntHandler(n *Node, ctx context.Context) error {
	a, err := util.Atoi(n.Token.Text)
	if err != nil {
		return n.Token.Pos.Error(err)
	}
	ctx.Push(a)

	return CallChildren(n, ctx)
}

// IdentHandler is used for label/variable lookup
func IdentHandler(n *Node, ctx context.Context) error {
	t := n.Token.Text

	if line, exists := ctx.GetLabel(t); exists {
		ctx.Push(line.Address)
	}

	return CallChildren(n, ctx)
}
