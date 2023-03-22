package m6502

import (
	"assembler/assembler/lexer"
	"assembler/assembler/parser"
	"github.com/peter-mount/go-kernel/v2"
	"strings"
)

func init() {
	kernel.Register(&M6502{})
}

// M6502 implements the 6502 processor.
// Subtypes like the 65c02 extends this for the additional instructions
type M6502 struct {
}

func (p *M6502) PostInit() error {
	parser.Register(p)
	return nil
}

func (p *M6502) ProcessorName() string { return "6502" }

// Parse parses the current token (always TokenIdent) and returns the current node parsed into an AST tree.

func (p *M6502) Parse(curNode *parser.Node, token *lexer.Token, tokens []*lexer.Token) (*parser.Node, error) {
	switch strings.ToLower(token.Text) {

	case "iny":
		curNode.GetLine().SetData(0xC8)
	case "rts":
		curNode.GetLine().SetData(0x60)

	// Return nil,nil to indicate the instruction was not recognised
	default:
		return nil, nil
	}

	// Return the current node to indicate we have parsed this token
	return curNode, nil
}
