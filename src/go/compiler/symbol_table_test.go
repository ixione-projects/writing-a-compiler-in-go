package compiler

import "testing"

type SymbolTest struct {
	name   string
	symbol Symbol
}

func TestSymbolTable(t *testing.T) {
	symbols :=
		[]SymbolTest{
			{"a", Symbol{Name: "a", Scope: GLOBAL_SCOPE, Index: 0}},
			{"b", Symbol{Name: "b", Scope: GLOBAL_SCOPE, Index: 1}},
		}

	globals := NewSymbolTable()
	for _, symbol := range symbols {
		globals.Define(symbol.name)
	}

	for i, expected := range symbols {
		actual, _ := globals.Lookup(expected.name)
		if expected.symbol != actual {
			t.Errorf("test[%d] - globals.Lookup(expected.name) ==> expected: <%v> but was: <%v>", i, expected.symbol, actual)
		}
	}
}
