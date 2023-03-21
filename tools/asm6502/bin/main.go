package main

import (
	common "assembler"
	"assembler/machine"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	if err := kernel.Launch(
		&common.VersionService{},
		&machine.Service{},
	); err != nil {
		log.Fatal(err)
	}
}
