package memory

// Address a generic memory address
type Address uint

func (a Address) Increment() Address {
	return a + 1
}

func (a Address) Add(n int) Address {
	return a + Address(n)
}

// ToLittleEndian returns the address as a 4 byte 64-bit address
// in Little Endian format.
func (a Address) ToLittleEndian() []byte {
	i := uint64(a)
	return []byte{
		byte(i & 0xff),
		byte((i >> 8) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 24) & 0xff),
	}
}

// ToBigEndian returns the address as a 4 byte 64-bit address
// in Big Endian format.
func (a Address) ToBigEndian() []byte {
	i := uint64(a)
	return []byte{
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	}
}

type AddressBlock struct {
	Name     string  `yaml:"name"`
	Start    Address `yaml:"start"`
	End      Address `yaml:"end"`
	ReadOnly bool    `yaml:"readOnly"`
	Notes    string  `yaml:"notes"`
}

func (r AddressBlock) Contains(addr Address) bool {
	return r.Start <= addr && addr <= r.End
}

// Intersects returns true if a Block intersects (overlaps) this one
func (r AddressBlock) Intersects(a AddressBlock) bool {
	// If start or end is within this block then there's some overlap
	return r.Contains(a.Start) || r.Contains(a.End)
}

// ContainsBlock returns true if this block is entirely contained within this block
func (r AddressBlock) ContainsBlock(a AddressBlock) bool {
	return r.Contains(a.Start) && r.Contains(a.End)
}
