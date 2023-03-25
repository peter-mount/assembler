package main

import (
	"fmt"
	common "github.com/peter-mount/assembler"
	"github.com/peter-mount/assembler/assembler"
	_ "github.com/peter-mount/assembler/m6502"
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
