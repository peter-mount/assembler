package m6502

import (
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/assembler/parser"
)

// M65816 implements the 16-bit 65816 processor by providing handlers for
// 65816 specific instructions. If an instruction is not defined here
// it then passes it to the 65c02 Processor for the common instructions.
type M65816 struct {
	M65c02       parser.Processor
	instructions *node.Map
}

func (p *M65816) PostInit() error {
	p.instructions = node.NewMap(
		node.Entry{Name: "ADC", Handler: adc(AMAddressLong)},
		node.Entry{Name: "COP", Handler: COP},
		node.Entry{Name: "STP", Handler: SimpleInstruction(0xdb)},
		node.Entry{Name: "WAI", Handler: SimpleInstruction(0xcb)},
		node.Entry{Name: "WDM", Handler: SimpleInstruction(0x42)},
		node.Entry{Name: "XCE", Handler: SimpleInstruction(0xfb)},
	)

	parser.Register(p)
	return nil
}

func (p *M65816) Start() error {
	p.M65c02 = parser.Lookup("65c02")
	return nil
}

func (p *M65816) ProcessorName() string { return "65816" }

// Parse parses the current token (always TokenIdent) and returns the current node parsed into an AST tree.
func (p *M65816) Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	if h, resolved := p.instructions.ResolveToken(token); resolved {
		return node.NewWithHandler(token, h), nil
	}

	// Pass on to the 6502 parser for the common ones
	return p.M65c02.Parse(token, tokens)
}
