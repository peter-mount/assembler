package main

import (
	common "assembler"
	"github.com/peter-mount/go-kernel/v2"
	"log"
)

func main() {
	err := kernel.Launch(&common.VersionService{})
	if err != nil {
		log.Fatal(err)
	}
}
