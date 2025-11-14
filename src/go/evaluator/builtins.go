package evaluator

import (
	"fmt"
	"strings"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
)

var builtins map[string]*object.BuiltinFunction

func init() {
	builtins = map[string]*object.BuiltinFunction{
		"len": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("len", args)
				}

				switch args[0].Type() {
				case object.STRING:
					return object.Number(len(args[0].(object.String))), nil
				case object.ARRAY:
					return object.Number(len(args[0].(*object.Array).Elements)), nil
				default:
					return nil, toBuiltinError("len", args)
				}
			},
		},
		"puts": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL, nil
			},
		},
		"first": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("first", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					if len(elements) == 0 {
						return NULL, nil
					}
					return elements[0], nil
				default:
					return nil, toBuiltinError("first", args)
				}
			},
		},
		"last": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("last", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					if len(elements) == 0 {
						return NULL, nil
					}
					return elements[len(elements)-1], nil
				default:
					return nil, toBuiltinError("last", args)
				}
			},
		},
		"rest": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("rest", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					length := len(elements)
					if length == 0 {
						return NULL, nil
					}
					result := make([]object.Object, length-1)
					copy(result, elements[1:])
					return &object.Array{Elements: result}, nil
				default:
					return nil, toBuiltinError("rest", args)
				}
			},
		},
		"push": {
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 2 {
					return nil, toBuiltinError("push", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					length := len(elements)
					result := make([]object.Object, length+1)
					copy(result, elements)
					result[length] = args[1]
					return &object.Array{Elements: result}, nil
				default:
					return nil, toBuiltinError("push", args)
				}
			},
		},
	}
}

func toBuiltinError(name string, args []object.Object) *object.Error {
	types := []string{}
	for _, arg := range args {
		types = append(types, arg.Type().String())
	}
	return toError("argument(s) to `%s` not supported: (%s)", name, strings.Join(types, ", "))
}
