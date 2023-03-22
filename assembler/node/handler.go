package node

import (
	"assembler/assembler/lexer"
	"context"
	"fmt"
	"strings"
)

type Handler func(*Node, context.Context) error

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
