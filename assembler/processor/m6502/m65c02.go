package m6502

import (
	"assembler/assembler/lexer"
	"assembler/assembler/parser"
	"strings"
)

// M65c02 implements the 65c02 processor.
type M65c02 struct {
	m6502 *M6502 `'kernel:"inject"`
}

func (p *M65c02) PostInit() error {
	parser.Register(p)
	return nil
}

func (p *M65c02) ProcessorName() string { return "65c02" }

// Parse parses the current token (always TokenIdent) and returns the current node parsed into an AST tree.

func (p *M65c02) Parse(curNode *parser.Node, token *lexer.Token, tokens []*lexer.Token) (*parser.Node, error) {
	switch strings.ToLower(token.Text) {

	case "stz":
		// FIXME dummy for testing
		curNode.GetLine().SetData(0x92, 0, 0)

	// Pass on to the 6502 parser for the common ones
	default:
		return p.m6502.Parse(curNode, token, tokens)
	}

	// Return the current node to indicate we have parsed this token
	return curNode, nil
}
