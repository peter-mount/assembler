package common

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
	"github.com/peter-mount/go-kernel/v2/log"
)

// LineHandler handles the processing of lexer.Line's
func LineHandler(n *node.Node, ctx context.Context) error {
	ctx.ClearStack()

	switch ctx.GetStage() {

	case context.StageCompile, context.StageOptimise, context.StageBackref:
		n.Line.Address = ctx.GetAddress()
		if err := node.CallChildren(n, ctx); err != nil {
			return err
		}
		ctx.AddAddress(len(n.Line.Data()))

	case context.StageList:
		log.Println(n.Line.String())

	}

	return node.CallChildren(n, ctx)
}
