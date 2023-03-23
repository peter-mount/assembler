package m6502

import (
	"assembler/assembler/common"
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/memory"
	"strings"
)

func JSR(n *node.Node, ctx context.Context) error {
	err := node.CallChildren(n, ctx)
	if err == nil {
		switch ctx.GetStage() {

		case context.StageCompile:
			// reserve 3 bytes
			ctx.AddAddress(3)

		case context.StageBackref:
			var r interface{}
			var addr memory.Address
			if err == nil {
				r, err = ctx.Pop()
				if err == nil {
					addr, err = common.ToAddr(r)
				}
			}

			if err != nil {
				return n.Token.Pos.Error(err)
			}

			b := addr.ToLittleEndian()
			n.GetLine().SetData(0x20, b[0], b[1])
		}
	}
	return err
}

var branchOpcodes = map[string]byte{
	"bcc": 0x90,
	"bcs": 0xb0,
	"beq": 0xf0,
	"bne": 0xd0,
	"bmi": 0x30,
	"bpl": 0x10,
	"bvc": 0x50,
	"bvs": 0x70,
}

func Branch(n *node.Node, ctx context.Context) error {
	// Resolve the opCode for this instruction
	opName := strings.ToLower(n.Token.Text)
	opCode, exists := branchOpcodes[opName]
	if !exists {
		return n.Token.Pos.Errorf("opcode %q not recognised", opName)
	}

	//log.Printf("Branch %d %q", n.Token.Token, n.Token.Text)
	err := node.CallChildren(n, ctx)
	if err == nil {
		switch ctx.GetStage() {

		case context.StageCompile:
			ctx.AddAddress(2)

		case context.StageBackref:
			var r interface{}
			var addr memory.Address
			if err == nil {
				r, err = ctx.Pop()
				if a, ok := r.(memory.Address); ok {
					addr = a
				} else if err == nil {
					addr, err = common.ToAddr(r)
				}
			}

			if err != nil {
				return n.Token.Pos.Error(err)
			}

			l := n.GetLine()
			// TODO check this is correct
			delta := int(addr) - int(l.Address)
			//log.Printf("Delta %d", delta)
			if delta < -127 || delta > 127 {
				return l.Pos.Errorf("Destination %q is %d bytes away, exceeding instruction")
			}

			n.GetLine().SetData(opCode, byte(delta))
		}
	}
	return nil
}
