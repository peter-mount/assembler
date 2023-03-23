package node

import (
	"assembler/assembler/context"
)

func CallChildren(n *Node, ctx context.Context) error {
	if err := n.Left.Visit(ctx); err != nil {
		return err
	}
	return n.Right.Visit(ctx)
}

func (n *Node) Visit(ctx context.Context) error {
	if n != nil && n.Handler != nil {
		err := n.Handler(n, ctx)
		if err != nil {
			return n.Token.Pos.Error(err)
		}
	}
	return nil
}
