package main

import (
	common "assembler"
	"assembler/assembler"
	_ "assembler/m6502"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	if err := kernel.Launch(
		&common.VersionService{},
		&assembler.Service{},
	); err != nil {
		log.Fatal(err)
	}
}
