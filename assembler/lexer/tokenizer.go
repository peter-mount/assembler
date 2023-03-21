package lexer

import (
	"strings"
	"unicode"
)

func (l *Lexer) tokenizeLine(line *Line) error {
	l.curLine = line

	c := strings.TrimSpace(line.Line)
	switch {
	case c == "":
		// Blank line so ignore
		return nil

	case c[0] == ';', c[0] == '*', c[0] == '#':
		// Line starts with a comment
		line.Comment = c
		line.Line = ""

	default:
		// TODO parse line based on style
		l.scanner.Init(strings.NewReader(line.Line))
		l.scanner.Whitespace = 0

		var tokens []*Token
		var token *Token
		for token == nil || token.Token != TokenEOF {
			token = l.scan()
			tok := token.Token
			if tok != TokenEOF {
				switch {
				case tok == TokenWhitespace:
					token.Text = token.Text + l.scanWhile(unicode.IsSpace)
				}
				tokens = append(tokens, token)
			}
		}

		var text []string
		tid := 0
		tlen := len(tokens)
		setLabel := true
		for tid < tlen {
			token := tokens[tid]
			tok := token.Token
			hasMore := (tid + 1) < tlen
			/*if hasMore {
				fmt.Printf("%d %q %d %q\n", tok, token.Text, tokens[tid+1].Token, tokens[tid+1].Text)
			} else {
				fmt.Printf("%d %q\n", tok, token.Text)
			}*/
			switch {
			// strip whitespace but clear setLabel as we can't now have one
			case tok == TokenWhitespace:
				setLabel = false

			// Line starts with .ident then skip, the next pass will set it as the label
			case setLabel && tok == '.' && hasMore && tokens[tid+1].Token == TokenIdent:
				// Skip, we will see it as an ident next

			// Line starts with an ident then it's the label
			case setLabel && tok == TokenIdent:
				setLabel = false
				line.Label = token.Text
				token.Token = TokenLabel

			default:
				line.Tokens = append(line.Tokens, token)
				text = append(text, token.Text)
			}
			tid++
		}
		line.Line = strings.Join(text, " ")

		/*for _, t := range line.Tokens {
			fmt.Printf("%s ", t.String())
		}
		fmt.Println()*/
	}

	// Add line to Lexer
	l.lines = append(l.lines, line)
	return nil
}

func (l *Lexer) scan() *Token {
	tok := &Token{Token: l.scanner.Scan(), Pos: l.curLine.Pos}
	t := tok.Token
	if t != TokenEOF {
		tok.Text = l.scanner.TokenText()
		if unicode.IsSpace(t) {
			tok.Token = TokenWhitespace
		}
	}
	return tok
}

func (l *Lexer) scanNext() string {
	l.scanner.Scan()
	return l.scanner.TokenText()
}

func (l *Lexer) scanWhile(f func(rune) bool) string {
	var s string
	for f(l.scanner.Peek()) {
		s = s + l.scanNext()
	}
	return s
}

func (l *Lexer) peek() rune {
	return l.scanner.Peek()
}
