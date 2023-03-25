package m6502

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/common"
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/node"
)

type AddressMode uint8

func (am AddressMode) String() string {
	i := int(am)
	if i > len(amNames) {
		i = 0
	}
	return amNames[i]
}

// Acceptable returns true if it matches one of the provided parameters
func (am AddressMode) Acceptable(accepts ...AddressMode) bool {
	if len(accepts) == 0 {
		// No entries in accepts then we accept all
		return true
	}

	for _, a := range accepts {
		if am == a {
			return true
		}
	}
	return false
}

// Validate returns true if the provided value is within range for this
// AddressMode.
func (am AddressMode) Validate(i uint) bool {
	switch am {
	case AMImmediate,
		AMZeroPage, AMZeroPageIndirectX, AMZeroPageIndirectY,
		AMZeroPageIndexedX, AMZeroPageIndexedY,
		AMZeroPageIndirect:
		return i <= 0xff

	case AMAddress, AMAbsoluteIndexedX, AMAbsoluteIndexedY:
		return i <= 0xffff

	case AMAddressLong:
		return i <= 0xffffff

	default:
		return false
	}
}

// Opcode returns a formatted opcode with the supplied value.
func (am AddressMode) Opcode(opcode byte, val uint) ([]byte, error) {
	if !am.Validate(val) {
		return nil, errors.OutOfBounds()
	}

	switch am {
	// Implied instruction
	case AMImplied:
		return []byte{opcode}, nil

	// 8-bit data
	case AMImmediate,
		AMZeroPage, AMZeroPageIndirectX, AMZeroPageIndirectY,
		AMZeroPageIndexedX, AMZeroPageIndexedY,
		AMZeroPageIndirect:
		return []byte{opcode, byte(val & 0xff)}, nil

	// 16-bit data
	case AMAddress, AMAbsoluteIndexedX, AMAbsoluteIndexedY:
		return []byte{opcode, byte(val & 0xff), byte((val >> 8) & 0xff)}, nil

	default:
		return nil, errors.IllegalArgument()
	}
}

// The detected addressing modes.
const (
	AMUnknown                 AddressMode = iota // Holder to represent an unknown AddressMode
	AMImplied                                    // No argument, also represents stack
	AMImmediate                                  // #0x12			8-bit value
	AMAddress                                    // 0x1234			16-bit address, 0x0100-0xffff inclusive. Values <0x0100 appear as AMZeroPage
	AMZeroPage                                   // 0x12			Zero page address, e.g. 0..ff inclusive
	AMAddressLong                                // 0x123456		Long addresses, values 0x10000 and up
	AMAbsoluteIndexedX                           // 0x1234,X 		Absolute Indexed X
	AMAbsoluteIndexedY                           // 0x1234,Y 		Absolute Indexed X
	AMZeroPageIndirect                           // (0x12)			65c02 only
	AMZeroPageIndirectX                          // (0x12,X)   		Direct Page Indexed X
	AMZeroPageIndirectY                          // (0x12),Y 		Direct page Index Y
	AMZeroPageIndexedX                           // 0x12,X   		Direct Page Indexed X
	AMZeroPageIndexedY                           // 0x12,Y 		Direct page Index Y
	AMZeroPageIndirectLong                       // [$12]			65816 only
	AMAccumulator                                // ASL A
	AMAbsoluteIndexedIndirect                    // (0x1234,X)		65816 only
	AMBlockMove                                  // source, dest	65816 only, consists of Address ',' Address
	amEndMarker                                  // must be last, used in tests to identify how many are defined
)

var (
	// The 3 address modes that are parsed together, in order of precedence if multiple ones are acceptable
	addressAddressModes = []AddressMode{AMAddressLong, AMAddress, AMZeroPage}
	// The order here must match the constants
	amNames = []string{
		"AMUnknown",
		"AMImplied",
		"AMImmediate",
		"AMAddress",
		"AMZeroPage",
		"AMAddressLong",
		"AMAbsoluteIndexedX",
		"AMAbsoluteIndexedY",
		"AMZeroPageIndirect",
		"AMZeroPageIndirectX",
		"AMZeroPageIndirectY",
		"AMZeroPageIndexedX",
		"AMZeroPageIndexedY",
		"AMZeroPageIndirectLong",
		"AMAccumulator",
		"AMAbsoluteIndexedIndirect",
		"AMBlockMove",
	}
)

// Addressing is the parsed result from handling
type Addressing struct {
	AddressMode AddressMode // Type of addressing found
	Value       uint        // First value found
	Value2      uint        // Second value found, for BlockMove
}

func GetAddressing(n *node.Node, ctx context.Context, accept ...AddressMode) (Addressing, error) {
	var err error
	addr := Addressing{}

	r := n.Right
	children := r.GetChildren()

	switch {
	// Implied has no more tokens
	case r == nil:
		addr.AddressMode = AMImplied

	// Immediate # value
	case node.MatchPattern(children, "#", ""):
		err = addr.getInt(children, 1, AMImmediate, ctx)

	// addr,X
	case node.MatchPattern(children, "", ",", "X"):
		err = addr.getInt(children, 0, AMAbsoluteIndexedX, ctx)

	// addr,Y
	case node.MatchPattern(children, "", ",", "Y"):
		err = addr.getInt(children, 0, AMAbsoluteIndexedY, ctx)

	// (addr)
	case node.MatchPattern(children, "(", "", ")"):
		err = addr.getInt(children, 1, AMZeroPageIndirect, ctx)

	// (addr),X
	case node.MatchPattern(children, "(", "", ",", "X", ")"):
		err = addr.getInt(children, 1, AMZeroPageIndirectX, ctx)

	// (addr),Y
	case node.MatchPattern(children, "(", "", ")", ",", "Y"):
		err = addr.getInt(children, 1, AMZeroPageIndirectY, ctx)

	// Default to address
	default:
		i, err := common.GetNodeInt(r, ctx)
		if err == nil {
			addr.Value = uint(i)

			// Now work out which one of addressAddressModes to use
			for _, a := range addressAddressModes {
				if a.Validate(addr.Value) && a.Acceptable(accept...) {
					addr.AddressMode = a
				}
			}

			// Verify the choice is correct. if none are acceptable then this should fail
			if !addr.AddressMode.Validate(addr.Value) {
				return addr, errors.OutOfBounds()
			}
		}
	}

	// If we are picky then verify we support the resolved addressing mode
	if err == nil && !addr.AddressMode.Acceptable(accept...) {
		err = n.Token.Pos.Errorf("Invalid addressing mode %d", addr.AddressMode)
	}

	if err != nil {
		err = n.Token.Pos.Error(err)
	}

	return addr, err
}

func (addr *Addressing) getInt(children []*node.Node, index int, am AddressMode, ctx context.Context) error {
	i, err := common.GetNodeInt(children[index], ctx)
	if err == nil {
		addr.AddressMode = am
		addr.Value = uint(i)
	}
	return err
}

func (addr *Addressing) String() string {
	if addr == nil {
		return "nil"
	}
	switch addr.AddressMode {
	case AMUnknown, AMImplied:
		return ""
	case AMImmediate:
		return fmt.Sprintf("#%x", addr.Value)
	case AMAddress, AMAddressLong, AMZeroPage:
		return fmt.Sprintf("%x", addr.Value)
	case AMAbsoluteIndexedX:
		return fmt.Sprintf("%x,X", addr.Value)
	case AMAbsoluteIndexedY:
		return fmt.Sprintf("%x,Y", addr.Value)
	case AMZeroPageIndirect:
		return fmt.Sprintf("(%x)", addr.Value)
	case AMZeroPageIndirectX:
		return fmt.Sprintf("(%x,X)", addr.Value)
	case AMZeroPageIndirectY:
		return fmt.Sprintf("(%x),Y", addr.Value)
	default:
		// FIXME catch all if we don't implement
		return fmt.Sprintf("?? %x ??", addr.Value)
	}
}
