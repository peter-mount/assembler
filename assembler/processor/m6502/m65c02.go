package m6502

import (
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/assembler/parser"
)

// M65c02 implements the 65c02 processor by providing handlers for
// 65c02 specific instructions. If an instruction is not defined here
// it then passes it to the 6502 Processor for the common instructions.
type M65c02 struct {
	m6502        *M6502 `'kernel:"inject"`
	instructions *node.Map
}

func (p *M65c02) PostInit() error {
	p.instructions = node.NewMap(
		node.Entry{Name: "BRA", Handler: BranchAlways},
	)

	parser.Register(p)
	return nil
}

func (p *M65c02) ProcessorName() string { return "65c02" }

// Parse parses the current token (always TokenIdent) and returns the current node parsed into an AST tree.
func (p *M65c02) Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	if h, resolved := p.instructions.ResolveToken(token); resolved {
		return node.NewWithHandler(token, h), nil
	}

	// Pass on to the 6502 parser for the common ones
	return p.m6502.Parse(token, tokens)
}
