package processor

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/node"
	"strings"
)

type Builder interface {
	Handle(string, node.Handler) Builder
	Simple(string, ...byte) Builder
	Extends(Processor) Builder
	Include(BuilderInclude) Builder
	Build() Processor
}

type BuilderInclude func(Builder)

type builder struct {
	name         string                 // The unique name of the processor flavour
	instructions *node.Map              // Instruction handlers for this processor flavour
	parent       Processor              // Parent processor if this is extending an older flavour
	reserved     map[string]interface{} // Reserved words
}

// New creates a new Processor Builder
func New(name string) Builder {
	return &builder{
		name:         name,
		instructions: node.NewMap(),
		reserved:     make(map[string]interface{}),
	}
}

// Handle adds a Handler that will handle a specific opcode
func (b *builder) Handle(name string, handler node.Handler) Builder {
	b.instructions.AddEntry(node.Entry{Name: name, Handler: handler})
	b.reserved[strings.ToLower(name)] = true
	return b
}

// Simple adds a static byte sequence that will be inserted for a specific opcode.
// This is used for Opcodes that have just one instruction, e.g. RTS on the 6502
func (b *builder) Simple(name string, bytes ...byte) Builder {
	return b.Handle(name, SimpleInstruction(bytes...))
}

// Extends sets the parent Processor flavour that this Processor will extend.
// This will panic if the builder has already been set to extend another processor.
func (b *builder) Extends(parent Processor) Builder {
	if b.parent != nil {
		panic(fmt.Errorf("cannot extend %q with %q as already extending %q", b.name, parent.ProcessorName(), b.parent.ProcessorName()))
	}
	b.parent = parent
	return b
}

// Include will pass a builder to a function implementing BuilderInclude.
// This is used when adding custom instruction sets.
// e.g. on the 6502 the conditional branch instructions are included as they have the same handler just
// different opcodes and makes the code easier to read.
// On other architectures you might have optional instruction extensions - think MMX on the x86 family.
// Here you would use Include to add MMX to the processor being built.
func (b *builder) Include(i BuilderInclude) Builder {
	i(b)
	return b
}

// Reserve reserves certain idents so that they cannot be used as labels or variable names.
// These are in addition to opcodes.
func (b *builder) Reserve(s ...string) Builder {
	for _, n := range s {
		b.reserved[strings.ToLower(n)] = true
	}
	return b
}

// Build the final Processor definition.
func (b *builder) Build() Processor {
	return &processor{
		name:         b.name,
		instructions: b.instructions,
		parent:       b.parent,
	}
}
