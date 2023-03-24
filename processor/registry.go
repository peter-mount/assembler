package processor

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"sort"
	"strings"
)

var processors = make(map[string]Processor)

func Register(procs ...Processor) {
	for _, proc := range procs {
		n := strings.ToLower(proc.ProcessorName())
		if _, exists := processors[n]; exists {
			panic(fmt.Errorf("processor %q already registered", n))
		}
		processors[n] = proc
	}
}

func Lookup(n string) Processor {
	return processors[strings.ToLower(n)]
}

// Registry is a Kernel service which simply lists the available processors on startup.
type Registry struct{}

func (p *Registry) Start() error {
	if log.IsVerbose() {
		var a []string
		for k, _ := range processors {
			a = append(a, k)
		}
		sort.SliceStable(a, func(i, j int) bool {
			return a[i] < a[j]
		})
		log.Printf("CPUs: %s", strings.Join(a, ", "))
	}
	return nil
}
