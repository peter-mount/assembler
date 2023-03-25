package node

import (
	"github.com/peter-mount/assembler/assembler/lexer"
	"strings"
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
	n := &Node{Token: token}

	// Ensure we have a default handler for specific token types
	switch token.Token {
	case lexer.TokenStart:
		n.Handler = CallChildren

	case lexer.TokenInt:
		n.Handler = IntHandler

	case lexer.TokenIdent:
		n.Handler = IdentHandler
	}

	return n
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
		b.Parent = n
	} else {
		n.Left.AddLeft(b)
	}
}

func (n *Node) AddRight(b *Node) {
	if n == nil || b == nil || n == b {
		return
	}

	if n.Right == nil {
		n.Right = b
		b.Parent = n
	} else {
		n.Right.AddRight(b)
	}
}

func (n *Node) AddAllRightTokens(b ...*lexer.Token) {
	for _, t := range b {
		n.AddRight(New(t))
	}
}

func (n *Node) AddAllRight(b ...*Node) {
	for _, b1 := range b {
		n.AddRight(b1)
	}
}

func (n *Node) GetChildren() []*Node {
	var r []*Node
	cn := n
	for cn != nil {
		r = append(r, cn)
		cn = cn.Right
	}
	return r
}

func (n *Node) GetTokenText() []string {
	var r []string
	cn := n
	for cn != nil {
		r = append(r, cn.Token.Text)
		cn = cn.Right
	}
	return r
}

func MatchPattern(src []*Node, pat ...string) bool {
	if len(src) != len(pat) {
		return false
	}

	for i, s := range pat {
		if s != "" && strings.ToLower(s) != strings.ToLower(src[i].Token.Text) {
			return false
		}
	}
	return true
}
