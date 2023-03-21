package machine

import (
	"assembler/memory"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Machines struct {
	Machines []*Machine `yaml:"machines"`
}

type Machine struct {
	Name        string                `yaml:"name"`
	Processor   string                `yaml:"processor"`
	Description string                `yaml:"description"`
	Notes       string                `yaml:"notes"`
	Memory      []memory.AddressBlock `yaml:"memory"`
}

func (m *Machine) create() (*memory.Map, error) {
	mem := &memory.Map{}
	for _, b := range m.Memory {
		if err := mem.AddBlock(memory.NewMemoryBlock(b)); err != nil {
			return nil, err
		}
	}

	return mem, nil
}

type Service struct {
	ShowMachines *bool `kernel:"flag,machines,Show available machine names"`
	machines     map[string]*Machine
}

func (s *Service) Start() error {
	s.machines = make(map[string]*Machine)

	// etc directory
	etc := filepath.Join(filepath.Dir(os.Args[0]), "../etc/machines")

	if err := walk.NewPathWalker().
		Then(s.loadMachineDefs).
		IsFile().
		PathHasSuffix(".yml").
		Walk(etc); err != nil {
		return err
	}

	if *s.ShowMachines {
		return fmt.Errorf("Available machines: %s\n",
			strings.Join(s.AvailableMachines(), ", "))
	}

	return nil
}

func (s *Service) CreateMachine(n string) (*memory.Map, error) {
	m, exists := s.machines[strings.ToLower(n)]
	if !exists {
		return nil, fmt.Errorf("machine %q is unknown", n)
	}
	return m.create()
}

func (s *Service) AvailableMachines() []string {
	var a []string
	for k, _ := range s.machines {
		a = append(a, k)
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	return a
}

func (s *Service) loadMachineDefs(path string, _ os.FileInfo) error {
	var machine Machines
	b, err := os.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(b, &machine)
	}
	if err == nil {
		for _, m := range machine.Machines {
			n := strings.ToLower(m.Name)
			if _, exists := s.machines[n]; exists {
				return fmt.Errorf("machine %q already defined", n)
			}
			s.machines[n] = m
		}
	}
	if err != nil {
		return fmt.Errorf("failed to read %q: %v", path, err)
	}
	return nil
}
