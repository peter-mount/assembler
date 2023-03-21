package lexer

import "unicode"

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// IsIdent returns true if a rune should be matched as an ident
func IsIdent(r rune) bool {
	return r == TokenIdent
	/*	return r == '=' || r == '<' || r == '>' ||
		r == '^' || r == '&' ||
		r == '%' || r == '!' ||
		r == ':' || r == ';'*/
}

// IsVariableStart true if the rune is valid for the first char of a variable name
func IsVariableStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// IsVariableSuccessor true if the rune is valid for the successive chars in a variable name
func IsVariableSuccessor(r rune) bool {
	return IsVariableStart(r) || r == '_' || (r >= '0' && r <= '9')
}

func IsPlusMinus(r rune) bool {
	return r == '+' || r == '-'
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
