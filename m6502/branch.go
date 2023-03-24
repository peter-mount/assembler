package m6502

import (
	"assembler/assembler/context"
	"assembler/assembler/node"
	"assembler/processor"
	"strings"
)

func JSR(n *node.Node, ctx context.Context) error {
	var err error
	switch ctx.GetStage() {

	case context.StageCompile:
		// reserve 3 bytes
		n.GetLine().SetData(0, 0, 0)

	case context.StageBackref:
		params, err := GetAddressing(n, ctx, AMAddress) // TODO AMAddressLong for JSL (JSR as alias)
		if err != nil {
			return err
		}

		b, err := params.AddressMode.Opcode(0x20, params.Value)
		if err != nil {
			return err
		}
		n.GetLine().SetData(b...)
	}

	return err
}

// conditional branch instruction opCodes, used by Branch
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

func addBranchOpcodes(b processor.Builder) {
	for k, _ := range branchOpcodes {
		b.Handle(k, Branch)
	}
}

// Branch handles the 6502 conditional branch instructions.
func Branch(n *node.Node, ctx context.Context) error {
	// Resolve the opCode for this instruction
	opName := strings.ToLower(n.Token.Text)
	opCode, exists := branchOpcodes[opName]
	if !exists {
		return n.Token.Pos.Errorf("opcode %q not recognised", opName)
	}
	return branchOp(opCode, n, ctx)
}

// BranchAlways is the 65c02 BRA instruction which shares the same
// underlying handler as the other branch relative instructions on the 6502
func BranchAlways(n *node.Node, ctx context.Context) error {
	return branchOp(0x80, n, ctx)
}

// branchOp common handler for Branch and BranchAlways
func branchOp(opCode byte, n *node.Node, ctx context.Context) error {
	switch ctx.GetStage() {

	case context.StageCompile:
		// Reserve 2 bytes
		n.GetLine().SetData(opCode, 0)

	case context.StageBackref:
		params, err := GetAddressing(n, ctx, AMAddress, AMAddressLong)
		if err != nil {
			return err
		}

		l := n.GetLine()
		// TODO check this is correct
		delta := int(params.Value) - int(l.Address)
		//log.Printf("Delta %d", delta)
		if delta < -127 || delta > 127 {
			return l.Pos.Errorf("Destination %q is %d bytes away, exceeding instruction", params.Value, delta)
		}

		n.GetLine().SetData(opCode, byte(delta))
	}
	return nil
}
