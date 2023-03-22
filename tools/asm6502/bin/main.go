package main

import (
	common "assembler"
	"assembler/assembler"
	_ "assembler/assembler/processor/m6502"
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
