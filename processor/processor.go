package processor

import (
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"strings"
)

type Processor interface {
	// ProcessorName the name of this processor
	ProcessorName() string

	// Parse parses an operation
	Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error)
}

type processor struct {
	name         string
	instructions *node.Map
	parent       Processor
	reserved     map[string]interface{}
}

func (p *processor) ProcessorName() string {
	return p.name
}

func (p *processor) Parse(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	if h, resolved := p.instructions.ResolveToken(token); resolved {
		return node.NewWithHandler(token, h), nil
	}

	if p.parent != nil {
		return p.parent.Parse(token, tokens)
	}

	return nil, errors.UnsupportedError(token.Text)
}

func (p *processor) Reserved(s string) bool {
	_, reserved := p.reserved[strings.ToLower(strings.TrimSpace(s))]
	return reserved
}
