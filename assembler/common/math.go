package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/memory"
	"github.com/peter-mount/assembler/util"
)

// ToInt converts an interface{} to an int64.
// This supports v if it's an int, int64, Label or string
func ToInt(v interface{}) (int64, error) {
	if a, ok := v.(*context.Value); ok {
		return a.Int(), nil
	}
	if i, ok := v.(int64); ok {
		return i, nil
	}
	if i, ok := v.(int); ok {
		return int64(i), nil
	}
	if a, ok := v.(memory.Address); ok {
		return int64(a), nil
	}
	if l, ok := v.(*lexer.Line); ok {
		return int64(l.Address), nil
	}
	if s, ok := v.(string); ok {
		return util.Atoi(s)
	}
	return 0, errors.IllegalArgument()
}

func ToAddr(v interface{}) (memory.Address, error) {
	i, err := ToInt(v)
	if err != nil {
		return 0, err
	}

	return memory.Address(i), nil
}

// GetNodeInterface visits a node and returns the top value from the stack.
func GetNodeInterface(n *node.Node, ctx context.Context) (interface{}, error) {
	err := n.Visit(ctx)
	if err != nil {
		return 0, err
	}

	r, err := ctx.Pop()
	if err != nil {
		return 0, err
	}
	return r, nil
}

// GetNodeInt visits a node and returns the top value from the stack
// as an int64.
func GetNodeInt(n *node.Node, ctx context.Context) (int64, error) {
	r, err := GetNodeInterface(n, ctx)
	if err != nil {
		return 0, err
	}

	return ToInt(r)
}

// GetNodeAddress visits a node and returns the top value from the stack
// as an Label.
func GetNodeAddress(n *node.Node, ctx context.Context) (memory.Address, error) {
	r, err := GetNodeInterface(n, ctx)
	if err != nil {
		return 0, err
	}

	if a, ok := r.(memory.Address); ok {
		return a, nil
	}

	return ToAddr(r)
}
