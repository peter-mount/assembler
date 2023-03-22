package common

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/memory"
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
	ctx.SetAddress(memory.Address(addr))

	return nil
}
