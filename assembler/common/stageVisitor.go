package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/node"
)

// DefaultStageVisitor returns a simple StageVisitor which visits a specific root node.
// It's mostly used within tests
func DefaultStageVisitor(rootNode *node.Node) context.StageVisitor {
	return func(stage context.Stage, ctx context.Context) error {
		switch stage {
		// normal handlers should never get these two so ignore them
		case context.StageInit, context.StageTokenize, context.StageParse:
			return nil
		// Don't do anything for debug output stages
		case context.StageList, context.StageSymbols:
			return nil
		// Visit for all other stages
		default:
			return rootNode.Visit(ctx)
		}
	}
}
