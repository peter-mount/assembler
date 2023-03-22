package context

import (
	"assembler/assembler/errors"
	"assembler/assembler/lexer"
	"assembler/memory"
	"github.com/peter-mount/go-kernel/v2/log"
	"sort"
	"strings"
	"time"
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
	GetStage() Stage
	ForEachStage(func(Stage, Context) error) error
	GetLabel(n string) *lexer.Line
	SetLabel(n string, line *lexer.Line) error
	GetLabels() []string
	GetStartAddress() memory.Address
	GetAddress() memory.Address
	SetAddress(memory.Address)
	AddAddress(int) memory.Address
	ClearStack()
	Push(interface{})
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
	now1 := time.Now()
	for stage := StageLex; stage < stageCount; stage++ {
		now2 := time.Now()
		c.stage = stage
		c.ClearStack()
		if err := f(stage, c); err != nil {
			return err
		}
		log.Printf("Stage %d took %v", stage, time.Now().Sub(now2))
	}
	log.Printf("Assembly took %v", time.Now().Sub(now1))
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
		log.Printf("Pop %v", v)
		return v, nil

	default:
		l := len(c.stack) - 1
		v := c.stack[l]
		c.stack = c.stack[:l]
		log.Printf("Pop %v <= %v", v, c.stack)
		return v, nil
	}
}
