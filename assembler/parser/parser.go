package parser

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/common"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/memory"
	processor2 "github.com/peter-mount/assembler/processor"
	"strings"
)

// Parser takes the tokenized lines and forms a series of AST trees
type Parser struct {
	ProcessorRegistry *processor2.Registry
	root              *node.Node
	org               memory.Address
	processor         processor2.Processor
}

func (p *Parser) Parse(lines []*lexer.Line) (*node.Node, error) {
	root := node.NewByRune(lexer.TokenStart)

	for lineNo, line := range lines {
		nextNode, err := p.parseLine(line)
		if err != nil {
			return nil, err
		}

		if nextNode == nil {
			panic(fmt.Errorf("nil nextNode line %d", lineNo))
		}

		// Chain the line on the root, don't add it to it's tree
		root.Handler = root.Handler.Then(node.HandlerAdaptor(nextNode))
	}

	return root, nil
}

func (p *Parser) parseLine(line *lexer.Line) (*node.Node, error) {
	lineNode := node.NewByRune(lexer.TokenLine)
	lineNode.Line = line
	lineNode.Handler = common.LineHandler

	for i, token := range line.Tokens {
		tok := token.Token
		switch tok {
		case lexer.TokenLabel:
			lineNode.AddLeft(node.NewWithHandler(token, common.LabelHandler))

		case lexer.TokenComment:
			// Drop comments

		case lexer.TokenIdent:
			n, err := p.parseOperand(token, line.Tokens[i+1:])
			if err != nil {
				return nil, token.Pos.Error(err)
			}
			lineNode.AddRight(n)
			return lineNode, nil

		default:
			return nil, token.Pos.Errorf("Unsupported token %d %c", tok, tok)
		}
	}

	return lineNode, nil
}

func (p *Parser) parseOperand(token *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	command := strings.ToLower(token.Text)
	switch {
	case command == "cpu" && len(tokens) > 0:
		p.processor = processor2.Lookup(tokens[0].Text)
		if p.processor == nil {
			return nil, token.Pos.Errorf("unsupported processor %q", tokens[0].Text)
		}
		return nil, nil

	case command == "org" && len(tokens) > 0:
		cn := node.NewWithHandler(token, common.OrgHandler)
		cn.AddAllRightTokens(tokens...)
		return cn, nil

	case command == "equb", command == "equs":
		cn := node.NewWithHandler(token, common.Equb)
		cn.AddAllRightTokens(tokens...)
		return cn, nil

	case command == "equw":
		cn := node.NewWithHandler(token, common.EquW)
		cn.AddAllRightTokens(tokens...)
		return cn, nil

	case command == "equl":
		cn := node.NewWithHandler(token, common.EquL)
		cn.AddAllRightTokens(tokens...)
		return cn, nil

	default:
		if p.processor == nil {
			return nil, token.Pos.Errorf("No processor set for operand %q", token.Text)
		}

		cn, err := p.processor.Parse(token, tokens)
		if err != nil {
			return nil, err
		}

		if cn != nil {
			cn.AddAllRightTokens(tokens...)
			return cn, nil
		}
	}

	return nil, token.Pos.Errorf("Unsupported operand %q", token.Text)
}
