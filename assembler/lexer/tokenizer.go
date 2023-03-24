package lexer

import (
	"assembler/util"
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
		lineText := line.Line

		// Strip any inline comment
		if ci := strings.IndexByte(lineText, ';'); ci > -1 {
			line.Comment = lineText[ci:]
			lineText = lineText[:ci]
		}

		l.scanner.Init(strings.NewReader(lineText))
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

		var tok2 *Token
		var text []string
		tid := 0
		tlen := len(tokens)
		setLabel := true
		for tid < tlen {
			token := tokens[tid]
			tok := token.Token
			hasMore := (tid + 1) < tlen
			if hasMore {
				tok2 = tokens[tid+1]
			} else {
				tok2 = nil
			}

			switch {
			// strip whitespace but clear setLabel as we can't now have one
			case tok == TokenWhitespace:
				setLabel = false

			// Strip "" from text as we want to handle this raw
			case tok == TokenString:
				token.Text = strings.Trim(token.Text, "\"")
				line.Tokens = append(line.Tokens, token)

			// Line starts with .ident then skip, the next pass will set it as the label
			case setLabel && tok == '.' && hasMore && tok2.Token == TokenIdent:
				// Skip, we will see it as an ident next

			// Line starts with an ident then it's the label
			case setLabel && tok == TokenIdent:
				setLabel = false
				line.Label = token.Text
				token.Token = TokenLabel
				line.Tokens = append(line.Tokens, token)

			// &1234 or &fedc are treated as hex values
			case (tok == '&' || tok == '$') && hasMore && tok2.Token == TokenInt,
				(tok == '&' || tok == '$') && hasMore && tok2.Token == TokenIdent && util.IsHex(tok2.Text):
				// Mark next token as an Int and merge the text with this one
				tok2.Token = TokenInt
				tok2.Text = token.Text + tok2.Text

			default:
				line.Tokens = append(line.Tokens, token)
				text = append(text, token.Text)
			}
			tid++
		}
		line.Line = strings.Join(text, " ")
	}

	// Handle line to Lexer
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
