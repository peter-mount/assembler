package parser

import "assembler/assembler/lexer"

type Processor interface {
	// ProcessorName the name of this processor
	ProcessorName() string

	// Parse parses an operation
	Parse(curNode *Node, token *lexer.Token, tokens []*lexer.Token) (*Node, error)
}
