package context

import (
	"assembler/memory"
	"bytes"
)

type Block struct {
	address memory.Address
	buffer  bytes.Buffer
}

func NewBlock(address memory.Address, b ...byte) *Block {
	block := &Block{address: address}
	_, _ = block.Write(b)
	return block
}

func (b *Block) Address() memory.Address {
	return b.address
}

func (b *Block) Bytes() []byte {
	return b.buffer.Bytes()
}

func (b *Block) Write(d []byte) (int, error) {
	if b == nil || len(d) == 0 {
		return 0, nil
	}
	return b.buffer.Write(d)
}

func BlocksEquals(a, b []*Block) bool {
	if len(a) != len(b) {
		return false
	}

	return true
}
