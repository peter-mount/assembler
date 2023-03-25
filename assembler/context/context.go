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
	StageLex                   // Load and lex the sources
	StageParse                 // Initial parsing stage
	StageCompile               // Compile opcodes
	StageOptimise              // Optimise stage to see if an instruction can be reduced in size
	StageBackref               // Resolve Back references
	StageList                  // List compiled listing
	StageSymbols               // List symbols
	StageAssemble              // Assembles each block
	stageCount                 // Must be last entry, not a real stage, used in ForEachStage()
)

type Stage int

type Context interface {
	// GetStage returns the current assembly Stage
	GetStage() Stage
	// ForEachStage calls a function once for each possible Stage
	ForEachStage(StageVisitor) error

	Get(string) (interface{}, error)
	Set(string, interface{}) error

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

	ClearBlocks()
	StartBlock(memory.Address)
	GetCurrentBlock() *Block
	GetAllBlocks() []*Block
}

type StageVisitor func(Stage, Context) error

type context struct {
	vars       map[string]interface{}
	labels     map[string]*lexer.Line
	stage      Stage
	orgAddress memory.Address // Address provided to last ORG statement
	address    memory.Address // Current assembly address
	stack      []interface{}  // Value stack
	curBlock   *Block         // Current block
	blocks     []*Block       // Compiled blocks
}

func New() Context {
	return &context{
		labels: make(map[string]*lexer.Line),
	}
}

func (c *context) GetStage() Stage {
	return c.stage
}

func (c *context) ForEachStage(f StageVisitor) error {
	for stage := StageLex; stage < stageCount; stage++ {
		c.stage = stage
		c.orgAddress = 0
		c.address = 0
		c.ClearStack()
		if err := f(stage, c); err != nil {
			return err
		}
	}
	return nil
}

func (c *context) Get(n string) (interface{}, error) {
	// Labels override variables
	if line := c.GetLabel(n); line != nil {
		return line.Address, nil
	}

	// TODO implement nesting here
	if v, exists := c.vars[n]; exists {
		return v, nil
	}

	// Variable does not exist
	switch c.GetStage() {

	// It does not exist, but may exist after this point so return 0 as a place holder
	case StageCompile:
		return 0, nil

	// Fail as it needs to exist by these stages
	case StageOptimise, StageBackref:
		return nil, errors.IllegalArgument()

	// Default TODO should we fail as this should not happen?
	default:
		return 0, nil
	}
}

func (c *context) Set(n string, v interface{}) error {
	switch c.GetStage() {
	// Set variable for the first time
	case StageCompile:
		if c.GetLabel(n) != nil {
			return errors.IllegalArgument()
		}

		c.vars[n] = v
		return nil

	// Disallow variables to be redefined
	case StageOptimise, StageBackref:
		return errors.IllegalArgument()

	// All other stages just ignore
	default:
		return nil
	}
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
