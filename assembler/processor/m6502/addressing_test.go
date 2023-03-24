package m6502

import (
	"testing"
)

func TestAddressMode_Acceptable(t *testing.T) {
	for am := AddressMode(0); am < amEndMarker; am++ {
		t.Run(am.String(), func(t *testing.T) {
			// Must accept itself
			if !am.Acceptable(am) {
				t.Errorf("AddressMode %d not matching itself", am)
			}

			// Test it does not match another AddressMode
			amn := (am + 3) % amEndMarker
			if am.Acceptable(amn) {
				t.Errorf("AddressMode %d matches %d when it shouldnt", am, amn)
			}
		})
	}
}

func TestAddressMode_Opcode(t *testing.T) {
}

func TestAddressMode_Validate(t *testing.T) {
}

func TestAddressing_getInt(t *testing.T) {
}

func TestGetAddressing(t *testing.T) {
}
