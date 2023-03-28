package lexer

import (
	"fmt"
	"text/scanner"
)

const (
	TokenEOF        = scanner.EOF       // End of file
	TokenIdent      = scanner.Ident     // Identifier
	TokenInt        = scanner.Int       // Integer
	TokenFloat      = scanner.Float     // Float
	TokenChar       = scanner.Char      // Single character
	TokenString     = scanner.String    // String
	TokenRawString  = scanner.RawString // Raw string
	TokenComment    = scanner.Comment   // Comment
	TokenWhitespace = -(iota + 1)       // Whitespace
	TokenStart                          // Start of assembly - parser only
	TokenLine                           // Start of a line
	TokenLabel                          // Set Label
	TokenOpcode                         // A parsed OPCode
	TokenData                           // A parsed data block
	TokenCalculator                     // Calculator expression
)

type Token struct {
	Token rune
	Text  string
	Pos   Position
}

func (t *Token) String() string {
	if t == nil {
		return "nil"
	} /*
		switch t.Token {
		case TokenIdent:
			return fmt.Sprintf("ident(%s)", t.Text)
		case TokenInt, TokenFloat, TokenString:
			return t.Text
		case TokenChar:
			return fmt.Sprintf("'%s'", t.Text)
		}*/
	return fmt.Sprintf("Token[token=%d,text=%s]", t.Token, t.Text)
}
