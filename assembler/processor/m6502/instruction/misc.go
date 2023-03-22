package instruction

import (
	"assembler/assembler/node"
	"context"
)

func NOP(node *node.Node, _ context.Context) error {
	node.GetLine().SetData(0xea)
	return nil
}
