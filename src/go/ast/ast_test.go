package ast

import (
	"testing"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetDeclaration{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "ident"},
					Value: "ident",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "value"},
					Value: "value",
				},
			},
		},
	}

	if "let ident=value;" != program.String() {
		t.Errorf("String() ==> expected: <%s> but was: <%s>", "let ident=value;", program.String())
	}
}
