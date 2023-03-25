package m6502

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/node"
)

const (
	adcImmediate                   = 0x69
	adcAddress                     = 0x6d
	adcAddressLong                 = 0x6f
	adcAddressZeroPage             = 0x65
	adcAddressZeroPageIndirect     = 0x72
	adcAddressZeroPageIndirectLong = 0x67
)

func adc(addressModes ...AddressMode) node.Handler {
	return func(n *node.Node, ctx context.Context) error {
		var err error
		switch ctx.GetStage() {

		// As this instruction can be between 2 and 6 bytes, we don't set any data here.
		// We will set the correct length in the later stages
		case context.StageCompile:

		// Create the actual opCode's with the appropriate lengths
		case context.StageOptimise, context.StageBackref:
			params, err := GetAddressing(n, ctx)
			if err != nil {
				return n.Token.Pos.Error(err)
			}

			if !params.AddressMode.Acceptable(addressModes...) {
				return n.Token.Pos.Error(errors.UnsupportedError("ADC %s", params.String()))
			}

			switch params.AddressMode {
			case AMImmediate:
				n.GetLine().SetData(adcImmediate, byte(params.Value))

			case AMZeroPage:
				n.GetLine().SetData(adcAddressZeroPage, byte(params.Value&0xff))

			case AMAddress:
				n.GetLine().SetData(adcAddress, byte(params.Value&0xff), byte((params.Value>>8)&0xff))

			case AMAddressLong:
				n.GetLine().SetData(adcAddressLong, byte(params.Value&0xff), byte((params.Value>>8)&0xff), byte((params.Value>>16)&0xff))

			case AMZeroPageIndirect:
				n.GetLine().SetData(adcAddressZeroPageIndirect, byte(params.Value))

			case AMZeroPageIndirectLong:
				n.GetLine().SetData(adcAddressZeroPageIndirectLong, byte(params.Value))
				return nil

			default:
				return errors.UnsupportedError("ADC %s", params.String())
			}
		}

		return err
	}
}
