package node

import (
	"assembler/assembler/context"
	"assembler/assembler/lexer"
	"fmt"
	"strings"
)

type Handler func(*Node, context.Context) error

func (a Handler) Then(b Handler) Handler {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(n *Node, ctx context.Context) error {
		if err := a(n, ctx); err != nil {
			return err
		}
		return b(n, ctx)
	}
}

func HandlerAdaptor(n *Node) Handler {
	return func(_ *Node, ctx context.Context) error {
		return n.Handler(n, ctx)
	}
}

type Map map[string]Handler

type Entry struct {
	Name    string
	Handler Handler
}

func NewMap(entries ...Entry) *Map {
	m := &Map{}
	for _, e := range entries {
		n := strings.ToLower(e.Name)
		if _, exists := (*m)[n]; exists {
			panic(fmt.Errorf("NodeHandlerMap already has %q registered", n))
		}
		(*m)[n] = e.Handler
	}
	return m
}

func (m *Map) ResolveToken(token *lexer.Token) (Handler, bool) {
	n := strings.ToLower(token.Text)
	h, exists := (*m)[n]
	if !exists {
		return nil, false
	}
	token.Token = lexer.TokenOpcode
	return h, true
}
