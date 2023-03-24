package main

import (
	common "assembler"
	"assembler/assembler"
	_ "assembler/m6502"
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
		&common.VersionService{},
		&assembler.Service{},
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
