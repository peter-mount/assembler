package processor

import (
	"assembler/assembler/node"
	"fmt"
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
	name         string
	instructions *node.Map
	parent       Processor
}

func New(name string) Builder {
	return &builder{
		name:         name,
		instructions: node.NewMap(),
	}
}

func (b *builder) Handle(name string, handler node.Handler) Builder {
	b.instructions.AddEntry(node.Entry{Name: name, Handler: handler})
	return b
}

func (b *builder) Simple(name string, bytes ...byte) Builder {
	return b.Handle(name, SimpleInstruction(bytes...))
}

func (b *builder) AddEntry(entries ...node.Entry) Builder {
	for _, entry := range entries {
		b.instructions.AddEntry(entry)
	}
	return b
}

func (b *builder) Extends(parent Processor) Builder {
	if b.parent != nil {
		panic(fmt.Errorf("cannot extend %q with %q as already extending %q", b.name, parent.ProcessorName(), b.parent.ProcessorName()))
	}
	b.parent = parent
	return b
}

func (b *builder) Include(i BuilderInclude) Builder {
	i(b)
	return b
}

func (b *builder) Build() Processor {
	return &processor{
		name:         b.name,
		instructions: b.instructions,
		parent:       b.parent,
	}
}
