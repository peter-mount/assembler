package context

import (
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/memory"
	"sort"
	"strings"
)

const (
	StageInit     Stage = iota // Initialisation of the Assembler
	StageTokenize              // Load and tokenize the sources
	StageParse                 // Initial parsing stage
	StageCompile               // Compile opcodes
	StageOptimise              // Optimise stage to see if an instruction can be reduced in size
	StageBackref               // Resolve Back references
	StageList                  // List compiled listing
	StageSymbols               // List symbols
	StageAssemble              // Assembles each block
	// Must be last entry, not a real stage, used in ForEachStage()
	stageCount
)

type Stage int

type Context interface {
	// GetStage returns the current assembly Stage
	GetStage() Stage
	// ForEachStage calls a function once for each possible Stage
	ForEachStage(StageVisitor) error

	Get(string) (*Value, error)
	Set(string, interface{}) error
	StartScope()
	EndScope()

	// GetLabel returns the Line that contains the given label
	GetLabel(n string) (*lexer.Line, bool)
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
	Pop() (*Value, error)
	Pop2() (*Value, *Value, error)
	Swap() error

	ClearBlocks()
	StartBlock(memory.Address)
	GetCurrentBlock() *Block
	GetAllBlocks() []*Block
}

type StageVisitor func(Stage, Context) error

type context struct {
	stage      Stage                  // Current stage during assembly
	orgAddress memory.Address         // Label provided to last ORG statement
	address    memory.Address         // Current assembly address
	curBlock   *Block                 // Current block
	blocks     []*Block               // Compiled blocks
	labels     map[string]*lexer.Line // Line labels
	stack      []*Value               // Value stack
	vars       []map[string]*Value    // Variables
}

func New() Context {
	return &context{}
}

func (c *context) GetStage() Stage {
	return c.stage
}

func (c *context) ForEachStage(f StageVisitor) error {
	for stage := StageInit; stage < stageCount; stage++ {
		c.stage = stage
		c.orgAddress = 0
		c.address = 0
		c.ClearStack()
		if stage == StageInit {
			c.vars = nil
			c.labels = make(map[string]*lexer.Line)
		} else if err := f(stage, c); err != nil {
			return err
		}
	}
	return nil
}

func (c *context) Get(n string) (*Value, error) {
	// Labels override variables
	if line, defined := c.GetLabel(n); defined {
		return LabelValue(line), nil
	}

	for _, v := range c.vars {
		if val, exists := v[n]; exists {
			return val, nil
		}
	}

	// Variable does not exist
	switch c.GetStage() {

	// It does not exist, but may exist after this point so return 0 as a place holder
	case StageCompile:
		return &zeroValue, nil

	// Fail as it needs to exist by these stages
	case StageOptimise, StageBackref:
		return nil, errors.IllegalArgument()

	// Default TODO should we fail as this should not happen?
	default:
		return &zeroValue, nil
	}
}

func (c *context) Set(n string, v interface{}) error {
	switch c.GetStage() {
	// Set variable for the first time
	case StageCompile:
		if _, defined := c.GetLabel(n); defined {
			return errors.IllegalArgument()
		}

		// Force a scope to start if we don't have none
		if len(c.vars) == 0 {
			c.StartScope()
		}

		// Find existing entry
		for _, v := range c.vars {
			if _, exists := v[n]; exists {
				v[n] = Of(v)
				return nil
			}
		}

		// Set in current scope
		c.vars[0][n] = Of(v)
		return nil

	// Disallow variables to be redefined
	case StageOptimise, StageBackref:
		return errors.IllegalArgument()

	// All other stages just ignore
	default:
		return nil
	}
}

func (c *context) StartScope() {
	c.vars = append([]map[string]*Value{make(map[string]*Value)}, c.vars...)
}

func (c *context) EndScope() {
	if len(c.vars) > 1 {
		c.vars = c.vars[1:]
	}
}

func (c *context) SetLabel(n string, line *lexer.Line) error {
	if _, exists := c.labels[n]; exists {
		return line.Pos.Errorf("label %q already defined", n)
	}
	c.labels[n] = line
	return nil
}
func (c *context) GetLabel(n string) (*lexer.Line, bool) {
	a, e := c.labels[n]
	return a, e
}

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
func (c *context) Push(v interface{}) { c.stack = append(c.stack, Of(v)) }

func (c *context) Pop() (*Value, error) {
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
func (c *context) Pop2() (*Value, *Value, error) {
	// b is first as it's the top value
	b, err := c.Pop()
	if err != nil {
		return nil, nil, err
	}

	// Now a
	a, err := c.Pop()
	if err != nil {
		return nil, nil, err
	}

	// a, b in the order you would expect
	return a, b, nil
}

// Swap swaps the top 2 values on the stack.
// Returns error if the stack doesn't have 2 items to swap.
func (c *context) Swap() error {
	l := len(c.stack)
	if l < 2 {
		return errors.StackEmpty()
	}

	c.stack[l-2], c.stack[l-1] = c.stack[l-1], c.stack[l-2]
	return nil
}

func (c *context) ClearBlocks() {
	c.curBlock = nil
	c.blocks = nil
}
func (c *context) StartBlock(address memory.Address) {
	c.curBlock = &Block{address: address}
	c.blocks = append(c.blocks, c.curBlock)
}
func (c *context) GetCurrentBlock() *Block { return c.curBlock }
func (c *context) GetAllBlocks() []*Block  { return c.blocks }
