package context

import (
	"assembler/assembler/errors"
	"assembler/assembler/lexer"
	"assembler/memory"
	"sort"
	"strings"
)

const (
	StageLex     Stage = iota // Load and lex the sources
	StageParse                // Initial parsing stage
	StageCompile              // Compile opcodes
	StageBackref              // Resolve Back references
	StageList                 // List compiled listing
	StageSymbols              // List symbols
	stageCount                // Must be last stage definition
)

type Stage int

type Context interface {
	// GetStage returns the current assembly Stage
	GetStage() Stage
	// ForEachStage calls a function once for each possible Stage
	ForEachStage(func(Stage, Context) error) error

	// GetLabel returns the Line that contains the given label
	GetLabel(n string) *lexer.Line
	// SetLabel sets the Line a label references
	SetLabel(n string, line *lexer.Line) error
	// GetLabels returns all labels in sorted order
	GetLabels() []string

	// GetStartAddress returns the value from the last ORG instruction
	GetStartAddress() memory.Address
	// GetAddress returns the current address being assembled
	GetAddress() memory.Address
	// SetAddress sets the next assembly address. Used by the ORG instruction
	SetAddress(memory.Address)
	// AddAddress adds a value to the current assembly address
	AddAddress(int) memory.Address

	// ClearStack clears the value stack
	ClearStack()
	// Push a value onto the top of the value stack
	Push(interface{})
	// Pop a value from the value stack
	Pop() (interface{}, error)
}

type context struct {
	labels     map[string]*lexer.Line
	stage      Stage
	orgAddress memory.Address
	address    memory.Address
	stack      []interface{}
}

func New() Context {
	return &context{
		labels: make(map[string]*lexer.Line),
	}
}

func (c *context) GetStage() Stage {
	return c.stage
}

func (c *context) ForEachStage(f func(Stage, Context) error) error {
	for stage := StageLex; stage < stageCount; stage++ {
		c.stage = stage
		c.ClearStack()
		if err := f(stage, c); err != nil {
			return err
		}
	}
	return nil
}

func (c *context) SetLabel(n string, line *lexer.Line) error {
	if _, exists := c.labels[n]; exists {
		return line.Pos.Errorf("label %q already defined", n)
	}
	c.labels[n] = line
	return nil
}
func (c *context) GetLabel(n string) *lexer.Line { return c.labels[n] }
func (c *context) GetLabels() []string {
	var a []string
	for k, _ := range c.labels {
		a = append(a, k)
	}
	sort.SliceStable(a, func(i, j int) bool {
		return strings.ToLower(a[i]) < strings.ToLower(a[j])
	})
	return a
}

func (c *context) GetAddress() memory.Address      { return c.address }
func (c *context) GetStartAddress() memory.Address { return c.orgAddress }
func (c *context) SetAddress(address memory.Address) {
	c.orgAddress = address
	c.address = address
}
func (c *context) AddAddress(delta int) memory.Address {
	c.address = c.address.Add(delta)
	return c.address
}

func (c *context) ClearStack()        { c.stack = nil }
func (c *context) Push(v interface{}) { c.stack = append(c.stack, v) }

func (c *context) Pop() (interface{}, error) {
	switch len(c.stack) {
	case 0:
		return nil, errors.StackEmpty()

	case 1:
		v := c.stack[0]
		c.stack = nil
		return v, nil

	default:
		l := len(c.stack) - 1
		v := c.stack[l]
		c.stack = c.stack[:l]
		return v, nil
	}
}
