package memory

import (
	"sort"
)

type Map struct {
	blocks []*Block
}

func (m *Map) AddBlock(b *Block) error {
	for _, b := range m.blocks {
		if b.Intersects(b) {
			return memoryOverlap
		}
	}
	m.blocks = append(m.blocks, b)
	sort.SliceStable(m.blocks, func(i, j int) bool {
		return m.blocks[i].End() < m.blocks[j].Start()
	})

	return nil
}

// GetBlock returns the Block containing a specific address
func (m *Map) GetBlock(addr Address) (*Block, error) {
	for _, b := range m.blocks {
		if b.Contains(addr) {
			return b, nil
		}
	}
	return nil, addressInvalid
}

func (m *Map) ReadByte(addr Address) (byte, error) {
	block, err := m.GetBlock(addr)
	if err != nil {
		return 0, err
	}
	return block.ReadByte(addr)
}

func (m *Map) WriteByte(addr Address, v byte) error {
	block, err := m.GetBlock(addr)
	if err == nil {
		err = block.WriteByte(addr, v)
	}
	return err
}
