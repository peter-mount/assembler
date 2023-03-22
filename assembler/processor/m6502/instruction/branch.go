package instruction

import (
	"assembler/assembler/node"
	"context"
	"log"
)

func RTS(node *node.Node, ctx context.Context) error {
	log.Println("RTS")
	node.GetLine().SetData(0x60)
	return nil
}

func JSR(node *node.Node, ctx context.Context) error {
	log.Println("JSR")
	// TODO placeholder
	node.GetLine().SetData(0xff)
	return nil
}
