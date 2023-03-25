package context

import (
	"bytes"
	"github.com/peter-mount/assembler/memory"
	"testing"
)

func TestNewBlock(t *testing.T) {
	testBlock(t, func(test testBlockTest) *Block {
		return NewBlock(test.Address, test.Data...)
	})
}

func TestBlock_Write(t *testing.T) {
	testBlock(t, func(test testBlockTest) *Block {
		block := NewBlock(test.Address)
		n, err := block.Write(test.Data)
		if err != nil {
			t.Error(err)
		}
		if n != len(test.Data) {
			t.Errorf("Expected to write %d only wrote %d", len(test.Data), n)
		}
		return block
	})
}

type testBlockTest struct {
	Address memory.Address
	Data    []byte
}

func testBlock(t *testing.T, f func(testBlockTest) *Block) {
	tests := []testBlockTest{
		{Address: 0x8000, Data: []byte{1, 2, 3, 4, 5}},
		{Address: 0x0e00, Data: []byte{0, 2, 4, 8, 16, 32, 64, 128}},
	}
	for i, tt := range tests {
		t.Run("NewBlock", func(t *testing.T) {
			block := f(tt)

			bb := block.Bytes()
			if !bytes.Equal(bb, tt.Data) {
				t.Errorf("block %d differs\ngot: %d [%x]\nexp: %d [%x]", i, len(bb), bb, len(tt.Data), tt.Data)
			}
		})
	}
}
