package common

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/memory"
	"github.com/peter-mount/go-kernel/v2/log"
)

func OrgHandler(n *node.Node, ctx context.Context) error {
	if ctx.GetStage() == context.StageCompile {
		err := node.CallChildren(n, ctx)

		var addr int64
		if err == nil {
			r, err := ctx.Pop()
			log.Printf("org %v", r)
			if err == nil {
				addr, err = ToInt(r)
			}
		}

		log.Printf("org %x err %v", addr, err)

		if err != nil {
			return n.Token.Pos.Error(err)
		}
		ctx.SetAddress(memory.Address(addr))
		log.Printf("ORG 0x%x", addr)
	}

	return node.CallChildren(n, ctx)
}
