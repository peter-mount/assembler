package parser

import (
	"assembler/assembler/lexer"
	"assembler/machine"
	"assembler/memory"
	"assembler/util"
	"fmt"
	"strings"
)

type Parser struct {
	ProcessorRegistry *ProcessorRegistry
	root              *Node
	machine           *machine.Machine
	org               memory.Address
	processor         Processor
}

func (p *Parser) Parse(lines []*lexer.Line) (*Node, error) {
	root := NewNodeByRune(lexer.TokenStart)
	curNode := root

	for _, line := range lines {
		nextNode, err := p.parseLine(curNode, line)
		if err != nil {
			return nil, err
		}
		curNode = nextNode
	}

	return root, nil
}

func (p *Parser) parseLine(curNode *Node, line *lexer.Line) (*Node, error) {
	lineNode := NewNodeByRune(lexer.TokenLine)
	lineNode.Line = line
	curNode.Right = lineNode
	curNode = lineNode

	// TODO this should be set at a later stage
	line.Address = p.org

	for i, token := range line.Tokens {
		tok := token.Token
		switch tok {
		case lexer.TokenLabel:
			lineNode.Left = NewNode(token)
		case lexer.TokenComment:
		// Drop comments
		case lexer.TokenIdent:
			n, err := p.parseOperand(curNode, token, line.Tokens[i+1:])
			if err != nil {
				return nil, token.Pos.Error(err)
			}
			return n, nil
		default:
			return nil, token.Pos.Errorf("Unsupported token %d %c", tok, tok)
		}
	}

	return curNode, nil
}

func (p *Parser) parseOperand(curNode *Node, token *lexer.Token, tokens []*lexer.Token) (*Node, error) {
	command := strings.ToLower(token.Text)
	switch {
	case command == "cpu" && len(tokens) > 0:
		p.processor = p.ProcessorRegistry.Lookup(tokens[0].Text)
		if p.processor == nil {
			return nil, token.Pos.Errorf("unsupported processor %q", tokens[0].Text)
		}

	case command == "machine":

	case command == "org":
		if tokens[0].Token == lexer.TokenInt {
			a, err := util.Atoi(tokens[0].Text)
			if err != nil {
				return nil, token.Pos.Error(err)
			}
			p.org = memory.Address(a)
		}

	case command == "equb", command == "equs":
		if err := p.parseEqub(curNode, tokens); err != nil {
			return nil, err
		}

	case command == "iny":
		p.org = p.org.Add(curNode.GetLine().SetData(0xC8))
	case command == "rts":
		p.org = p.org.Add(curNode.GetLine().SetData(0x60))

	default:
		if p.processor == nil {
			return nil, token.Pos.Errorf("No processor set for operand %q", token.Text)
		}

		cn, err := p.processor.Parse(curNode, token, tokens)
		if err != nil {
			return nil, err
		}

		if cn == nil {
			// FIXME instruction not recognised so fail here
			curNode.GetLine().SetData(0xea)
			//return nil, token.Pos.Errorf("Unsupported operand %q", token.Text)
		}

		// It's now an opcode if it's still an Ident
		if curNode.Token.Token == lexer.TokenIdent {
			curNode.Token.Token = lexer.TokenOpcode
		}
	}

	// FIXME for now increment this, move it to the main assembler stage
	p.org = p.org.Add(len(curNode.GetLine().Data()))

	return curNode, nil
}

func (p *Parser) parseEqub(curNode *Node, tokens []*lexer.Token) error {
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
				return err
			}
			// TODO lookahead for ';' separator, if present then a 16-bit integer not 8-bit

			// Handle negative values
			if a < 0 {
				a = a + 255
			}

			// Validate range, TODO handle 16-bit here as well
			if a < 0 || a > 255 {
				return fmt.Errorf("%q is not a byte", token.Text)
			}

			b = append(b, byte(a))

		// TODO If TokenIdent then do a variable/label lookup here

		// Ignore valid value separators
		case ',', ';':

		default:
			return token.Pos.Errorf("unsupported token %q", token.Text)
		}
	}

	curNode.Token.Token = lexer.TokenData
	curNode.GetLine().SetData(b...)

	return nil
}
