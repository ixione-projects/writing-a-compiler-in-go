package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
)

type ObjectType int

const (
	FUNCTION ObjectType = iota
	BUILTIN
	NUMBER
	BOOLEAN
	STRING
	ARRAY
	HASH
	NULL
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Function struct {
	Literal *ast.FunctionLiteral
	Closure *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	params := []string{}
	for _, param := range f.Literal.Parameters {
		params = append(params, param.Value)
	}
	return "<fn (" + strings.Join(params, ", ") + ")>"
}

type BuiltinCall func(ctx *Environment, args ...Object) (Object, Interruption)

type BuiltinFunction struct {
	Fn BuiltinCall
}

func (bf *BuiltinFunction) Type() ObjectType {
	return BUILTIN
}

func (bf *BuiltinFunction) Inspect() string {
	return "<fn builtin>"
}

type (
	Number  float64
	Boolean bool
	String  string
)

func (n Number) Type() ObjectType {
	return NUMBER
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN
}

func (s String) Type() ObjectType {
	return STRING
}

func (n Number) Inspect() string {
	return fmt.Sprintf("%g", float64(n))
}

func (b Boolean) Inspect() string {
	return fmt.Sprintf("%t", bool(b))
}

func (s String) Inspect() string {
	return "\"" + string(s) + "\""
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elems := []string{}
	for _, elem := range a.Elements {
		elems = append(elems, elem.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
}

type Hashable interface {
	Object
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (n Number) HashKey() HashKey {
	return HashKey{n.Type(), uint64(n)}
}

func (b Boolean) HashKey() HashKey {
	var value uint64
	if b {
		value = 1
	} else {
		value = 0
	}
	return HashKey{b.Type(), value}
}

func (s String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s))
	return HashKey{s.Type(), h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type InterruptionType int

const (
	RETURN_VALUE InterruptionType = iota
	ERROR
)

type Interruption interface {
	Type() InterruptionType
	Inspect() string
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() InterruptionType {
	return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() InterruptionType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

var objects = map[ObjectType]string{
	FUNCTION: "FUNCTION",
	BUILTIN:  "BUILTIN",
	NUMBER:   "INTEGER",
	BOOLEAN:  "BOOLEAN",
	STRING:   "STRING",
	ARRAY:    "ARRAY",
	HASH:     "HASH",
	NULL:     "NULL",
}

func (ot ObjectType) String() string {
	return objects[ot]
}

var interruptions = map[InterruptionType]string{
	RETURN_VALUE: "RETURN_VALUE",
	ERROR:        "ERROR",
}

func (it InterruptionType) String() string {
	return interruptions[it]
}
