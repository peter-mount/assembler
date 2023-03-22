package parser

import (
	"assembler/assembler/lexer"
)

type Node struct {
	Parent *Node        // Parent node, nil=root
	Left   *Node        // Left hand side
	Right  *Node        // Right hand side
	Token  *lexer.Token // Token
	Line   *lexer.Line  // Optional pointer to line
}

func NewNode(token *lexer.Token) *Node {
	return &Node{Token: token}
}

func NewNodeByRune(tokenId rune) *Node {
	return NewNode(&lexer.Token{Token: tokenId})
}

func (n *Node) GetLine() *lexer.Line {
	if n.Line != nil {
		return n.Line
	}
	if n.Parent != nil {
		return n.Parent.GetLine()
	}
	return nil
}
