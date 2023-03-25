package common

import (
	"bytes"
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"testing"
)

func TestDataBlock(t *testing.T) {

	tests := [][]byte{
		{},
		{1, 2, 3, 4, 5},
		{0, 2, 4, 8, 16, 32, 64, 128},
	}
	for i, ob := range tests {
		t.Run("TestDataBlock", func(t *testing.T) {

			line := &lexer.Line{}
			lineNode := node.NewWithHandler(&lexer.Token{Token: lexer.TokenIdent}, LineHandler)
			lineNode.Line = line
			lineNode.AddRight(node.NewWithHandler(&lexer.Token{Token: lexer.TokenIdent}, DataBlock(ob...)))

			rootNode := node.NewByRune(lexer.TokenStart)
			rootNode.Handler = rootNode.Handler.Then(node.HandlerAdaptor(lineNode))

			ctx := context.New()
			if err := ctx.ForEachStage(DefaultStageVisitor(rootNode)); err != nil {
				t.Error(err)
			} else {
				bb := line.Data()
				if !bytes.Equal(bb, ob) {
					t.Errorf("block %d differs\ngot: %d [%x]\nexp: %d [%x]", i, len(bb), bb, len(ob), ob)
				}
			}
		})
	}
}
