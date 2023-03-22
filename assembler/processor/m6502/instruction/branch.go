package instruction

import (
	"assembler/assembler/common"
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/memory"
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

func Branch(n *node.Node, ctx context.Context) error {
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

			// TODO lookup op code from n.Token.Text

			n.GetLine().SetData(0x20, byte(delta))
		}
	}
	return nil
}
