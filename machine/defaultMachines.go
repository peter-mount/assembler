package machine

import "assembler/memory"

func init() {
	Register(NewMachine(
		"bbc",
		"6502",
		memory.AddressBlock{Start: 0x0000, End: 0x0dff, Name: "OS Workspace", ReadOnly: true},
		memory.AddressBlock{Start: 0x0e00, End: 0x7fff, Name: "Workspace+Screen"},
		memory.AddressBlock{Start: 0x8000, End: 0xbfff, Name: "Paged Rom", ReadOnly: true},
		memory.AddressBlock{Start: 0xc000, End: 0xffff, Name: "OS", ReadOnly: true},
	))

	// C64 default layout
	Register(NewMachine(
		"c64",
		"6510",
		memory.AddressBlock{Start: 0x0000, End: 0x03ff, Name: "Kernal Workspace", ReadOnly: true},
		memory.AddressBlock{Start: 0x0400, End: 0x07ff, Name: "Default Screen"},
		memory.AddressBlock{Start: 0x0800, End: 0x7fff, Name: "Workspace"},
		memory.AddressBlock{Start: 0x8000, End: 0xbfff, Name: "Basic Rom", ReadOnly: true},
		memory.AddressBlock{Start: 0xc000, End: 0xcfff, Name: "Upper memory"},
		memory.AddressBlock{Start: 0xd000, End: 0xffff, Name: "Kernal", ReadOnly: true},
	))

	// C64 with basic paged out
	Register(NewMachine(
		"c64-no-basic",
		"6510",
		memory.AddressBlock{Start: 0x0000, End: 0x03ff, Name: "Kernal Workspace", ReadOnly: true},
		memory.AddressBlock{Start: 0x0400, End: 0x07ff, Name: "Default Screen"},
		memory.AddressBlock{Start: 0x0800, End: 0x7fff, Name: "Workspace"},
		memory.AddressBlock{Start: 0x8000, End: 0xbfff, Name: "Basic swapped out"},
		memory.AddressBlock{Start: 0xc000, End: 0xcfff, Name: "Upper memory"},
		memory.AddressBlock{Start: 0xd000, End: 0xffff, Name: "Kernal", ReadOnly: true},
	))

	Register(NewMachine(
		"spectrum16k",
		"z80",
		memory.AddressBlock{Start: 0x0000, End: 0x3fff, Name: "ROM", ReadOnly: true},
		memory.AddressBlock{Start: 0x4000, End: 0x4fff, Name: "Screen"},
		memory.AddressBlock{Start: 0x5000, End: 0x7fff, Name: "Contended memory"},
	))

	Register(NewMachine(
		"spectrum48k",
		"z80",
		memory.AddressBlock{Start: 0x0000, End: 0x3fff, Name: "ROM", ReadOnly: true},
		memory.AddressBlock{Start: 0x4000, End: 0x4fff, Name: "Screen"},
		memory.AddressBlock{Start: 0x5000, End: 0x7fff, Name: "Contended memory"},
		memory.AddressBlock{Start: 0x8000, End: 0xffff, Name: "Uncontended memory"},
	))
}
