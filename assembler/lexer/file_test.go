package lexer

import (
	"testing"
)

func TestScanLine(t *testing.T) {
	crSrc := []byte("line1\rline2\rline 3")
	lfSrc := []byte("line1\nline2\nline 3")
	crLfSrc := []byte("line1\r\nline2\r\nline 3")
	lfCrSrc := []byte("line1\n\rline2\n\rline 3")

	tests := []struct {
		name   string
		offset int
		src    []byte
		want   int
		want1  string
		want2  bool
	}{
		// CR 		Acorn BBC, Commodore 8-Bit, ZX Spectrum
		{name: "CR1", offset: 0, src: crSrc, want: 6, want1: "line1", want2: false},
		{name: "CR2", offset: 6, src: crSrc, want: 12, want1: "line2", want2: false},
		{name: "CR3", offset: 12, src: crSrc, want: 18, want1: "line 3", want2: true},
		// LF 		Unix, Linux, Amiga, RiscOS
		{name: "LF1", offset: 0, src: lfSrc, want: 6, want1: "line1", want2: false},
		{name: "LF2", offset: 6, src: lfSrc, want: 12, want1: "line2", want2: false},
		{name: "LF3", offset: 12, src: lfSrc, want: 18, want1: "line 3", want2: true},
		// CR,LF	CP/M, MS-DOS, Windows, Atari TOS, Amstrad CPC
		{name: "CRLF1", offset: 0, src: crLfSrc, want: 7, want1: "line1", want2: false},
		{name: "CRLF2", offset: 7, src: crLfSrc, want: 14, want1: "line2", want2: false},
		{name: "CRLF3", offset: 14, src: crLfSrc, want: 20, want1: "line 3", want2: true},
		// LF,CR	Acorn BBC, Risc OS spooled output, BBC MOS OSASCII
		{name: "LFCR1", offset: 0, src: lfCrSrc, want: 7, want1: "line1", want2: false},
		{name: "LFCR2", offset: 7, src: lfCrSrc, want: 14, want1: "line2", want2: false},
		{name: "LFCR3", offset: 14, src: lfCrSrc, want: 20, want1: "line 3", want2: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ScanLine(tt.offset, tt.src)
			if got != tt.want {
				t.Errorf("ScanLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ScanLine() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ScanLine() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
