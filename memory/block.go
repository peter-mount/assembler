package memory

type Block struct {
	bounds AddressBlock
	data   []byte
}

func NewBlock(name string, addr Address, size int) *Block {
	return &Block{
		bounds: AddressBlock{
			Name:  name,
			Start: addr,
			End:   addr + Address(size-1),
		},
		data: make([]byte, size),
	}
}

func NewMemoryBlock(bounds AddressBlock) *Block {
	return &Block{
		bounds: bounds,
		data:   make([]byte, bounds.End-bounds.Start+1),
	}
}

func (b *Block) Name() string { return b.bounds.Name }

// Start address for this block
func (b *Block) Start() Address { return b.bounds.Start }

// End address of this block. This is the last address within the block.
func (b *Block) End() Address { return b.bounds.End }

// Size of this block
func (b *Block) Size() int { return len(b.data) }

func (b *Block) Contains(addr Address) bool {
	return b != nil && b.bounds.Contains(addr)
}

// Intersects returns true if a Block intersects (overlaps) this one
func (b *Block) Intersects(a *Block) bool {
	return a != nil && b != nil && b.bounds.Intersects(a.bounds)
}

// ContainsBlock returns true if this block is entirely contained within this block
func (b *Block) ContainsBlock(a *Block) bool {
	return a != nil && b != nil && b.bounds.ContainsBlock(a.bounds)
}

func (b *Block) ReadByte(addr Address) (byte, error) {
	if !b.Contains(addr) {
		return 0, addressInvalid
	}

	return b.data[addr-b.bounds.Start], nil
}

func (b *Block) WriteByte(addr Address, v byte) error {
	if !b.Contains(addr) {
		return addressInvalid
	}

	if b.bounds.ReadOnly {
		return readOnlyMemory
	}

	b.data[addr-b.bounds.Start] = v
	return nil
}
