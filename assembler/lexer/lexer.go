package lexer

import (
	"github.com/peter-mount/go-kernel/v2/log"
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
		log.Printf("Exiting %q", l.curFile.Name())
		nl := len(l.fileStack)
		if nl > 0 {
			l.curFile = l.fileStack[nl-1]
			l.fileStack = l.fileStack[:nl-1]
		} else {
			l.curFile = nil
		}
		if l.curFile != nil {
			log.Printf("Resuming %q", l.curFile.Name())
		}
	}()

	log.Printf("Entering %q", l.curFile.Name())

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
