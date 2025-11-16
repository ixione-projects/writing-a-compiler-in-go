package compiler

type SymbolScope int

const (
	GLOBAL_SCOPE SymbolScope = iota
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: map[string]Symbol{},
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	s := Symbol{
		Name:  name,
		Scope: GLOBAL_SCOPE,
		Index: len(st.store),
	}
	st.store[name] = s

	return s
}

func (st *SymbolTable) Lookup(name string) (Symbol, bool) {
	obj, found := st.store[name]
	return obj, found
}
