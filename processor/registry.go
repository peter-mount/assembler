package processor

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"sort"
	"strings"
)

var processors = make(map[string]Processor)

func Register(p Processor) {
	n := strings.ToLower(p.ProcessorName())
	//if _, exists := processors[n]; exists {
	//	panic(fmt.Errorf("processor %q already registered", n))
	//}
	processors[n] = p
}

func Lookup(n string) Processor {
	return processors[strings.ToLower(n)]
}

type ProcessorRegistry struct{}

func (p *ProcessorRegistry) Start() error {
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
