package memory

// Address a generic memory address
type Address uint

func (a Address) Increment() Address {
	return a + 1
}

func (a Address) Add(n int) Address {
	return a + Address(n)
}

type AddressBlock struct {
	Name       string
	Start, End Address
	ReadOnly   bool
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
