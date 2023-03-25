package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/memory"
)

func OrgHandler(n *node.Node, ctx context.Context) error {
	err := node.CallChildren(n, ctx)

	var r interface{}
	var addr int64
	if err == nil {
		r, err = ctx.Pop()
		if err == nil {
			addr, err = ToInt(r)
		}
	}
	if err != nil {
		return n.Token.Pos.Error(err)
	}

	// Set assembly address
	ctx.SetAddress(memory.Address(addr))

	// During StageAssemble start a new block to contain the final output
	if ctx.GetStage() == context.StageAssemble {
		ctx.StartBlock(memory.Address(addr))
	}

	return nil
}
