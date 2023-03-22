package parser

import (
	"assembler/assembler/lexer"
	"assembler/assembler/node"
)

type Processor interface {
	// ProcessorName the name of this processor
	ProcessorName() string

	// Parse parses an operation
	Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error)
}
