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
	root    *Node
	machine *machine.Machine
	org     memory.Address
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
	switch strings.ToLower(token.Text) {
	case "machine":
	case "org":
		if tokens[0].Token == lexer.TokenInt {
			a, err := util.Atoi(tokens[0].Text)
			if err != nil {
				return nil, token.Pos.Error(err)
			}
			p.org = memory.Address(a)
		}
	case "equb", "equs":
		if err := p.parseEqub(curNode, tokens); err != nil {
			return nil, err
		}

	case "iny":
		p.org = p.org.Add(curNode.GetLine().SetData(0xC8))
	case "rts":
		p.org = p.org.Add(curNode.GetLine().SetData(0x60))

	default:
		// TODO implenment, for now NOP
		p.org = p.org.Add(curNode.GetLine().SetData(0xea))
		//return nil, token.Pos.Errorf("Unsupported operand %q", token.Text)
	}
	return curNode, nil
}

func (p *Parser) parseEqub(curNode *Node, tokens []*lexer.Token) error {
	var b []byte
	for _, token := range tokens {
		switch token.Token {
		case lexer.TokenString:
			b = append(b, token.Text[1:len(token.Text)-2]...)
		case lexer.TokenRawString:
			b = append(b, token.Text...)
		case lexer.TokenInt:
			a, err := util.Atoi(token.Text)
			if err != nil {
				return err
			}
			if a < 0 {
				a = a + 255
			}
			if a < 0 || a > 255 {
				return fmt.Errorf("%q is not a byte", token.Text)
			}
			b = append(b, byte(a))
		}
	}
	line := curNode.Line
	if line != nil {
		line.Data = b
		p.org.Add(len(b))
	}
	return nil
}
