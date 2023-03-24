package m6502

import "github.com/peter-mount/go-kernel/v2"

func init() {
	kernel.Register(
		&M6502{},
		&M65c02{},
	)
}
