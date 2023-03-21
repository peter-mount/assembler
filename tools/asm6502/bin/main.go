package main

import (
	common "assembler"
	"assembler/assembler"
	"assembler/machine"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	if err := kernel.Launch(
		&common.VersionService{},
		&machine.Service{},
		&assembler.Assembler{},
	); err != nil {
		log.Fatal(err)
	}
}
