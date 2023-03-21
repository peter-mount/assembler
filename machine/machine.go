package machine

import (
	"assembler/memory"
	"fmt"
	"strings"
	"sync"
)

var (
	machines = make(map[string]Machine)
	mutex    sync.Mutex
)

type Machine interface {
	// Name of the machine
	Name() string
	// Processor the processor used by this machine
	Processor() string
	// Create a memory.Map relating to this machine
	Create() (*memory.Map, error)
}

func Register(machine Machine) {
	mutex.Lock()
	defer mutex.Unlock()

	n := strings.ToLower(machine.Name())
	if _, exists := machines[n]; exists {
		panic(fmt.Errorf("machine %q already registered", n))
	}
	machines[n] = machine
}

func NewMachine(name, processor string, blocks ...memory.AddressBlock) Machine {
	return &basicMachine{name: name, processor: processor, blocks: blocks}
}

type basicMachine struct {
	name      string
	processor string
	blocks    []memory.AddressBlock
}

func (m *basicMachine) Name() string {
	return m.name
}

func (m *basicMachine) Processor() string {
	return m.processor
}

func (m *basicMachine) Create() (*memory.Map, error) {
	mem := &memory.Map{}
	for _, b := range m.blocks {
		if err := mem.AddBlock(memory.NewMemoryBlock(b)); err != nil {
			return nil, err
		}
	}

	return mem, nil
}
