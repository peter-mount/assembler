package m6502

import (
	"assembler/assembler/common"
	"assembler/assembler/context"
	"assembler/assembler/node"
	"github.com/peter-mount/go-kernel/v2/log"
)

type AddressMode uint8

// The detected addressing modes.
//
// Note: Some of these are multi-use.
// For example, Address represents Zero-Page, Absolute or Absolute-Long
const (
	AMUnknown                 AddressMode = iota // Holder to represent an unknown AddressMode
	AMImplied                                    // No argument, also represents stack
	AMImmediate                                  // #0x12			8-bit value
	AMAddress                                    // raw address, can be 8-bit (zero-page), 16-bit or 24-bit (65816 long)
	AMAccumulator                                // ASL A
	AMAbsoluteIndexedIndirect                    // (0x1234,X)		65816 only
	AMZeroPageIndirect                           // (0x12)			65c02 only
	AMZeroPageIndirectLong                       // [$12]			65816 only
	AMBlockMove                                  // source, dest	65816 only, consists of Address ',' Address
)

// Addressing is the parsed result from handling
type Addressing struct {
	AddressMode AddressMode // Type of addressing found
	Register    Register    // Register found in the addressing
	Value       uint        // First value found
	Value2      uint        // Second value found, for BlockMove
}

func GetAddressing(n *node.Node, ctx context.Context, accept ...AddressMode) (Addressing, error) {
	var err error
	addr := Addressing{}

	r := n.Right
	switch {
	// Implied has no more tokens
	case r == nil:
		addr.AddressMode = AMImplied

	// Immediate # value
	case r.Token.Token == '#' && r.Right != nil:
		addr.AddressMode = AMImmediate
		i, err := common.GetNodeInt(r, ctx)
		if err == nil {
			addr.Value = uint(i)
		}

	// Default to address
	default:
		addr.AddressMode = AMAddress
		i, err := common.GetNodeInt(r, ctx)
		if err == nil {
			addr.Value = uint(i)
		}
	}

	// If we are picky then verify we support the resolved addressing mode
	if err == nil && len(accept) > 0 {
		found := false
		for _, a := range accept {
			if a == addr.AddressMode {
				found = true
			}
		}
		if !found {
			err = n.Token.Pos.Errorf("Invalid addressing mode %d", addr.AddressMode)
		}
	}

	if err != nil {
		log.Printf("%02d %s", n.Token.Pos.Line, n.GetLine().String())
		err = n.Token.Pos.Error(err)
	}
	return addr, err
}
