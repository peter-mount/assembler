package m6502

import (
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/assembler/parser"
	"assembler/assembler/processor/m6502/instruction"
)

// M6502 implements the 6502 processor.
// Subtypes like the 65c02 extends this for the additional instructions
type M6502 struct {
	instructions *node.Map
}

func (p *M6502) PostInit() error {
	NOP := instruction.SimpleInstruction(0xea)
	p.instructions = node.NewMap(
		node.Entry{Name: "BEQ", Handler: NOP},
		node.Entry{Name: "BNE", Handler: NOP},
		node.Entry{Name: "INX", Handler: instruction.SimpleInstruction(0xe8)},
		node.Entry{Name: "INY", Handler: instruction.SimpleInstruction(0xc8)},
		node.Entry{Name: "JSR", Handler: instruction.JSR},
		node.Entry{Name: "LDA", Handler: NOP},
		node.Entry{Name: "LDX", Handler: NOP},
		node.Entry{Name: "LDY", Handler: NOP},
		node.Entry{Name: "NOP", Handler: NOP},
		node.Entry{Name: "RTS", Handler: instruction.SimpleInstruction(0x60)},
	)

	parser.Register(p)
	return nil
}

func (p *M6502) ProcessorName() string { return "6502" }

// Parse parses the current token (always TokenIdent) and returns the current node parsed into an AST tree.

func (p *M6502) Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	if h, resolved := p.instructions.ResolveToken(token); resolved {
		return node.NewWithHandler(token, h), nil
	}

	// Return nil,nil to indicate the instruction was not recognised
	return nil, nil
}
