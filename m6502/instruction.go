package m6502

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/node"
)

func instruction(m map[AddressMode]byte, addressModes []AddressMode) node.Handler {

	// Validate we have the required entries
	for _, am := range addressModes {
		if _, exists := m[am]; !exists {
			panic(fmt.Errorf("address mode %s(%d) is not in definition", am.String(), am))
		}
	}

	return func(n *node.Node, ctx context.Context) error {
		var err error
		switch ctx.GetStage() {

		case context.StageCompile:
			// Do nothing until we have all labels set

		case context.StageOptimise, context.StageBackref:
			params, err := GetAddressing(n, ctx)
			if err != nil {
				return n.Token.Pos.Error(err)
			}

			if !params.AddressMode.Acceptable(addressModes...) {
				return n.Token.Pos.Error(errors.UnsupportedError("ADC %s", params.String()))
			}

			if opCode, exists := m[params.AddressMode]; exists {
				if !params.AddressMode.Validate(params.Value) {
					return errors.IllegalArgument()
				}
				switch params.AddressMode {
				case AMImplied:
					n.GetLine().SetData(opCode)

				case AMImmediate,
					AMZeroPage, AMZeroPageIndirect,
					AMZeroPageIndirectX, AMZeroPageIndirectY,
					AMZeroPageIndexedX, AMZeroPageIndexedY:
					n.GetLine().SetData(opCode, byte(params.Value))

				case AMAddress, AMAbsoluteIndexedX, AMAbsoluteIndexedY, AMAbsoluteIndexedIndirect:
					n.GetLine().SetData(opCode, byte(params.Value&0xff), byte((params.Value>>8)&0xff))

				case AMAddressLong:
					n.GetLine().SetData(opCode, byte(params.Value&0xff), byte((params.Value>>8)&0xff), byte((params.Value>>16)&0xff))

				}
			} else {
				return errors.UnsupportedError("%s %s", n.Token.Text, params.String())
			}
		}

		return err
	}
}
