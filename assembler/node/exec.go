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
	/*if n == nil {
		log.Printf("Node: nil")
	} else {
		log.Printf("Node: %d %v", n.Token.Token, n.Handler)
	}*/

	if n != nil && n.Handler != nil {
		return n.Handler(n, ctx)
	}
	return nil
}
