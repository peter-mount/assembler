package lexer

import (
	"text/scanner"
)

type Lexer struct {
	curFile   *File           // Current file being read
	fileStack []*File         // Stack of files, used when including additional files
	lines     []*Line         // Lines parsed
	scanner   scanner.Scanner // Tokenizer
	curLine   *Line           // Line being tokenized
}

func (l *Lexer) Lines() []*Line { return l.lines }

func (l *Lexer) Parse(fileName string) error {
	file, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	if l.curFile != nil {
		l.fileStack = append(l.fileStack, l.curFile)
	}

	l.curFile = file
	defer func() {
		nl := len(l.fileStack)
		if nl > 0 {
			l.curFile = l.fileStack[nl-1]
			l.fileStack = l.fileStack[:nl-1]
		} else {
			l.curFile = nil
		}
		if l.curFile != nil {
		}
	}()

	return l.curFile.ForEach(l.tokenizeLine)
}

func (l *Lexer) ForEach(h func(*Line) error) error {
	if l != nil {
		for _, line := range l.lines {
			if err := h(line); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *Lexer) ParseLines(script []string) error {
	for si, s := range script {
		line := &Line{
			Pos:  Position{File: &File{fileName: "unittest"}, Line: si + 1, Pos: -1},
			Line: s,
		}
		if err := l.tokenizeLine(line); err != nil {
			return err
		}
	}
	return nil
}
