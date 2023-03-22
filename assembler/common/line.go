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

	case context.StageCompile:
		n.Line.Address = ctx.GetAddress()

	case context.StageList:
		log.Println(n.Line.String())

	}

	return node.CallChildren(n, ctx)
}
