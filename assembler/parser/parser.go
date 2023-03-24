package parser

import (
	"assembler/assembler/common"
	"assembler/assembler/lexer"
	"assembler/assembler/node"
	"assembler/machine"
	"assembler/memory"
	processor2 "assembler/processor"
	"assembler/util"
	"fmt"
	"strings"
)

// Parser takes the tokenized lines and forms a series of AST trees
type Parser struct {
	ProcessorRegistry *processor2.Registry
	root              *node.Node
	machine           *machine.Machine
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

	case command == "machine":
		return nil, nil

	case command == "org":
		cn := node.NewWithHandler(token, common.OrgHandler)
		cn.AddAllRightTokens(tokens...)
		return cn, nil

	case command == "equb", command == "equs":
		return p.parseEqub(token, tokens)

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

func (p *Parser) parseEqub(tok *lexer.Token, tokens []*lexer.Token) (*node.Node, error) {
	var b []byte
	for _, token := range tokens {
		switch token.Token {

		// Strings are just appended as-is
		case lexer.TokenString, lexer.TokenRawString:
			b = append(b, token.Text...)

		// Integers
		case lexer.TokenInt:
			a, err := util.Atoi(token.Text)
			if err != nil {
				return nil, err
			}
			// TODO lookahead for ';' separator, if present then a 16-bit integer not 8-bit

			// Handle negative values
			if a < 0 {
				a = a + 255
			}

			// Validate range, TODO handle 16-bit here as well
			if a < 0 || a > 255 {
				return nil, fmt.Errorf("%q is not a byte", token.Text)
			}

			b = append(b, byte(a))

		// TODO If TokenIdent then do a variable/label lookup here

		// Ignore valid value separators
		case ',', ';':

		default:
			return nil, token.Pos.Errorf("unsupported token %q", token.Text)
		}
	}

	tok.Token = lexer.TokenData
	return node.NewWithHandler(tok, common.DataBlock(b...)), nil
}
