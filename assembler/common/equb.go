package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/errors"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/util"
)

func Equb(n *node.Node, ctx context.Context) error {
	var b []byte
	for cn := n.Right; cn != nil; cn = cn.Right {
		token := cn.Token
		switch token.Token {

		// Strings are just appended as-is
		case lexer.TokenString, lexer.TokenRawString:
			b = append(b, token.Text...)

		// Integers
		case lexer.TokenInt:
			a, err := util.Atoi(token.Text)
			if err != nil {
				return token.Pos.Error(err)
			}

			// Handle negative values
			if a < 0 {
				a = a + 256
			}
			if a < 0 || a > 255 {
				return token.Pos.Errorf("%q is not a byte", token.Text)
			}

			b = append(b, byte(a))

		case lexer.TokenIdent:
			v, err := ctx.Get(token.Text)
			if err != nil {
				return token.Pos.Error(err)
			}

			a, err := ToInt(v)
			if err != nil {
				return token.Pos.Error(err)
			}

			// Handle negative values
			if a < 0 {
				a = a + 256
			}
			if a < 0 || a > 255 {
				return token.Pos.Errorf("%q is not a byte", token.Text)
			}

			b = append(b, byte(a))

		// TODO If TokenIdent then do a variable/label lookup here

		// Ignore valid value separators
		case ',':

		default:
			return errors.UnsupportedError("unsupported token %q", cn.Token.Text)
		}
	}

	n.GetLine().SetData(b...)
	return nil
}

func EquW(n *node.Node, ctx context.Context) error {
	return equv(n, ctx, "word", 1<<16)
}

func EquL(n *node.Node, ctx context.Context) error {
	return equv(n, ctx, "long", 1<<32)
}

func equv(n *node.Node, ctx context.Context, t string, max int64) error {
	var b []byte
	for cn := n.Right; cn != nil; cn = cn.Right {
		token := cn.Token
		switch token.Token {

		// Integers
		case lexer.TokenInt:
			a, err := util.Atoi(token.Text)
			if err != nil {
				return token.Pos.Error(err)
			}

			// Handle negative values
			if a < 0 {
				a = a + max
			}
			if a < 0 || a >= max {
				return token.Pos.Errorf("%q is not a %s", token.Text, t)
			}

			b = append(b, byte(a))

		// TODO If TokenIdent then do a variable/label lookup here

		// Ignore valid value separators
		case ',':

		default:
			return errors.UnsupportedError("unsupported token %q", cn.Token.Text)
		}
	}

	n.GetLine().SetData(b...)
	return nil
}
