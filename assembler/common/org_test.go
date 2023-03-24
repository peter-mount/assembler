package common

import (
	"assembler/assembler/context"
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/memory"
	"fmt"
	"testing"
)

func TestOrgHandler(t *testing.T) {

	tests := []memory.Address{0x0e00, 0x1000, 0x2345, 0x8000}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("0x%x", tt), func(t *testing.T) {

			line := &lexer.Line{}
			lineNode := node.NewWithHandler(&lexer.Token{Token: lexer.TokenIdent}, LineHandler)
			lineNode.Line = line

			cn := node.NewWithHandler(&lexer.Token{Token: lexer.TokenIdent}, OrgHandler)
			cn.AddRight(node.New(&lexer.Token{Token: lexer.TokenInt, Text: fmt.Sprintf("%v", tt)}))
			lineNode.AddRight(cn)

			rootNode := node.NewByRune(lexer.TokenStart)
			rootNode.Handler = rootNode.Handler.Then(node.HandlerAdaptor(lineNode))

			ctx := context.New()
			if err := ctx.ForEachStage(DefaultStageVisitor(rootNode)); err != nil {
				t.Error(err)
			} else {
				blocks := ctx.GetAllBlocks()
				switch len(blocks) {
				case 0:
					t.Error("No blocks returned")
				default:
					if blocks[0].Address() != tt {
						t.Errorf("Got 0x%x expected 0x%x", blocks[0].Address(), tt)
					}
				}
			}
		})
	}
}
