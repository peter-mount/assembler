package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"testing"
)

func TestLabelHandler(t *testing.T) {
	tests := []struct {
		name  string
		label string
	}{
		{name: "NoLabel"},
		{name: "SetLabel", label: "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			line := &lexer.Line{Label: tt.label}
			lineNode := node.NewWithHandler(&lexer.Token{Token: lexer.TokenIdent}, LineHandler)
			lineNode.Line = line

			lineNode.AddLeft(node.NewWithHandler(&lexer.Token{Token: lexer.TokenLabel, Text: tt.label}, LabelHandler))

			rootNode := node.NewByRune(lexer.TokenStart)
			rootNode.Handler = rootNode.Handler.Then(node.HandlerAdaptor(lineNode))

			ctx := context.New()

			if err := ctx.ForEachStage(DefaultStageVisitor(rootNode)); err != nil {
				t.Error(err)
			}

			labels := ctx.GetLabels()
			// If we requested a label then we should have just 1 and it should match
			if tt.label != "" && len(labels) > 0 {
				if len(labels) != 1 {
					t.Errorf("Got %d labels expected 1", len(labels))
				} else if labels[0] != tt.label {
					t.Errorf("Got label %q expected %q", labels[0], tt.label)
				}
			}
		})
	}
}
