package m6502

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/processor"
)

type ldOpcode struct {
	name     string                // Instruction name
	operands map[AddressMode]uint8 // Map of AddressMode -> OpCode
}

var ldOpcodes6502 = []*ldOpcode{
	{name: "LDA", operands: map[AddressMode]uint8{
		AMImmediate:         0xA9,
		AMAddress:           0xAD,
		AMZeroPage:          0xA5,
		AMZeroPageIndirect:  0xA5,
		AMZeroPageIndirectX: 0xA1,
		AMZeroPageIndirectY: 0xB1,
		AMZeroPageIndexedX:  0xB5,
		AMAbsoluteIndexedX:  0xBD,
		AMAbsoluteIndexedY:  0xB9,
	}},
	{name: "LDX", operands: map[AddressMode]uint8{
		AMImmediate:        0xA2,
		AMAddress:          0xAE,
		AMZeroPage:         0xA6,
		AMAbsoluteIndexedY: 0xBE,
		AMZeroPageIndexedY: 0xB6,
	}},
	{name: "LDY", operands: map[AddressMode]uint8{
		AMImmediate:        0xA0,
		AMAddress:          0xAC,
		AMZeroPage:         0xA4,
		AMAbsoluteIndexedX: 0xBC,
		AMZeroPageIndexedX: 0xB4,
	}},
}

func addRegisterInstructions(defs []*ldOpcode) processor.BuilderInclude {
	return func(b processor.Builder) {
		for _, def := range defs {
			b.Handle(def.name, ld(def))
		}
	}
}

func ld(def *ldOpcode) node.Handler {
	// Slice used to pass to GetAddressing for the valid addressing modes
	// for this instruction
	var accept []AddressMode
	for k, _ := range def.operands {
		accept = append(accept, k)
	}

	return func(n *node.Node, ctx context.Context) error {
		switch ctx.GetStage() {

		case context.StageCompile:
			// Reserve 2 bytes as that's most common. we can expand in later stages
			n.GetLine().SetData(0, 0)

		case context.StageOptimise, context.StageBackref:
			params, err := GetAddressing(n, ctx, accept...)
			if err != nil {
				return n.Token.Pos.Error(err)
			}

			opcode, exists := def.operands[params.AddressMode]
			if !exists {
				return errors.IllegalArgument()
			}

			b, err := params.AddressMode.Opcode(opcode, params.Value)
			if err != nil {
				return err
			}
			n.GetLine().SetData(b...)
		}
		return nil
	}
}
