package lexer

import (
	"os"
	"path/filepath"
)

type File struct {
	fileName string
	lines    []*Line
}

func (f *File) Name() string {
	if f == nil {
		return "nil"
	}
	return f.fileName
}

func (f *File) BaseName() string {
	return filepath.Base(f.Name())
}

func (f *File) ForEach(h func(*Line) error) error {
	if f != nil {
		for _, line := range f.lines {
			if err := h(line); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadFile a file splitting the lines within it.
//
// A line is defined as text with either \n, \r\n or \r as the line terminators.
// Specifically:
// \n	Unix/Linux files
// \r\n	Windows formatted files
// \r	BBC Micro files
func ReadFile(fileName string) (*File, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	file := &File{fileName: fileName}
	p := 0
	scanning := true
	for scanning {
		np, line, eof := ScanLine(p, b)
		p = np
		scanning = !eof
		file.lines = append(file.lines, &Line{
			Pos: Position{
				File: file,
				Line: len(file.lines) + 1,
				Pos:  -1, // -1 as this is the entire line
			},
			Line: line,
		})
	}

	return file, nil
}

const (
	LF = '\n' // Line Feed
	VT = 0x0b // Vertical Tab
	FF = 0x0c // Form Feed
	CR = '\r' // Carriage Return
)

// ScanLine scans for the next line from a position in a byte slice.
// It returns the start of the next line, the line just found, and a boolean which
// is true when at the end of the file.
//
// Note we don't use bufio.Scanner here because we need to support different
// line encodings, whilst Scanner only supports \n and \r\n
//
// Here we look for both \r and \n and split the line at that point.
// The following sequences will produce blank lines:
// \n\n, \r\r, \r\n\r\n, \n\r\n\r
func ScanLine(p int, b []byte) (int, string, bool) {
	//bp := b[p:]

	l := len(b)
	e, r, l1 := l, l, l-1
	for i := p; r == l && i < l; i++ {
		switch b[i] {
		// CR or CRLF
		case CR:
			e, r = testChar(LF, i, l1, b)
			break

		// LF or LFCR
		case LF:
			e, r = testChar(CR, i, l1, b)
			break

		// Unicode standard for line terminators
		case FF, VT:
			e, r = i, i+1
			break
		}
	}

	return r, string(b[p:e]), r >= l
}

func testChar(c2 byte, i, l1 int, b []byte) (int, int) {
	if i < l1 && b[i+1] == c2 {
		return i, i + 2
	}
	return i, i + 1
}
