package common

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/node"
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
		return nil

	case context.StageList:
		fmt.Println(n.Line.String())
		ctx.AddAddress(len(n.Line.Data()))

	case context.StageAssemble:
		n.Line.Address = ctx.GetAddress()
		if err := node.CallChildren(n, ctx); err != nil {
			return err
		}
		ctx.AddAddress(len(n.Line.Data()))
		ctx.GetCurrentBlock().Write(n.Line.Data())
		return nil
	}

	return node.CallChildren(n, ctx)
}
