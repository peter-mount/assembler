package common

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/memory"
	"assembler/util"
	"github.com/peter-mount/go-kernel/v2/log"
)

func OrgHandler(n *node.Node, ctx context.Context) error {
	if ctx.GetStage() == context.StageCompile {
		a, err := util.Atoi(n.Token.Text)
		if err != nil {
			return n.Token.Pos.Error(err)
		}
		ctx.SetAddress(memory.Address(a))
		log.Printf("ORG 0x%x", a)
	}

	return node.CallChildren(n, ctx)
}
