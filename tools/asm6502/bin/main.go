package main

import (
	common "assembler"
	"assembler/assembler"
	_ "assembler/m6502"
	"assembler/machine"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	if err := kernel.Launch(
		&common.VersionService{},
		&machine.Service{},
		&assembler.Service{},
	); err != nil {
		log.Fatal(err)
	}
}
