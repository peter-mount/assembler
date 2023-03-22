package node

import (
	"assembler/assembler/context"
	"assembler/util"
)

func IntHandler(n *Node, ctx context.Context) error {
	switch ctx.GetStage() {

	case context.StageCompile:
		a, err := util.Atoi(n.Token.Text)
		if err != nil {
			return n.Token.Pos.Error(err)
		}
		ctx.Push(a)
	}

	return CallChildren(n, ctx)
}
