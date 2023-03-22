package node

import (
	"assembler/assembler/lexer"
)

type Node struct {
	Parent  *Node        // Parent node, nil=root
	Left    *Node        // Left hand side
	Right   *Node        // Right hand side
	Token   *lexer.Token // Token
	Line    *lexer.Line  // Optional pointer to line
	Handler Handler      // Handler to execute this node
}

func New(token *lexer.Token) *Node {
	return &Node{Token: token}
}

func NewByRune(tokenId rune) *Node {
	return New(&lexer.Token{Token: tokenId})
}

func NewWithHandler(token *lexer.Token, handler Handler) *Node {
	n := New(token)
	n.Handler = handler
	return n
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

func (n *Node) AddLeft(b *Node) {
	if n == nil || b == nil || n == b {
		return
	}

	if n.Left == nil {
		n.Left = b
		//} else {
		//	n.Left.AddRight(b)
	}
}
func (n *Node) AddRight(b *Node) {
	if n == nil || b == nil || n == b {
		return
	}

	if n.Right == nil {
		n.Right = b
		//} else {
		//	n.Right.AddRight(b)
	}
}
