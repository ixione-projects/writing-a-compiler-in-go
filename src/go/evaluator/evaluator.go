package evaluator

import (
	"fmt"
	"math"
	"slices"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
)

const (
	TRUE  = object.Boolean(true)
	FALSE = object.Boolean(false)
)

var NULL = &object.Null{}

func Evaluate(node ast.Node, env *object.Environment) (object.Object, object.Interruption) {
	switch node.Type() {
	case ast.PROGRAM:
		return evaluateProgram(node.(*ast.Program), env)
	case ast.ERROR:
		return nil, &object.Error{Message: node.(*ast.Error).Message}
	case ast.LET_DECLARATION:
		return evaluateLetDeclaration(node.(*ast.LetDeclaration), env)
	case ast.RETURN_STATEMENT:
		return evaluateReturnStatement(node.(*ast.ReturnStatement), env)
	case ast.EXPRESSION_STATEMENT:
		return evaluateExpressionStatement(node.(*ast.ExpressionStatement), env)
	case ast.BLOCK_STATEMENT:
		return evaluateBlockStatement(node.(*ast.BlockStatement), env)
	case ast.UNARY_EXPRESSION:
		return evaluateUnaryExpression(node.(*ast.UnaryExpression), env)
	case ast.BINARY_EXPRESSION:
		return evaluateBinaryExpression(node.(*ast.BinaryExpression), env)
	case ast.LOGICAL_EXPRESSION:
		return evaluateLogicalExpression(node.(*ast.LogicalExpression), env)
	case ast.CONDITIONAL_EXPRESSION:
		return evaluateConditionalExpression(node.(*ast.ConditionalExpression), env)
	case ast.FUNCTION_LITERAL:
		return evaluateFunctionLiteral(node.(*ast.FunctionLiteral), env)
	case ast.ASSIGNMENT_EXPRESSION:
		return evaluateAssignmentExpression(node.(*ast.AssignmentExpression), env)
	case ast.CALL_EXPRESSION:
		return evaluateCallExpression(node.(*ast.CallExpression), env)
	case ast.SUBSCRIPT_EXPRESSION:
		return evaluateSubscriptExpression(node.(*ast.SubscriptExpression), env)
	case ast.IDENTIFIER:
		return evaluateIdentifier(node.(*ast.Identifier), env)
	case ast.NUMBER_LITERAL:
		return object.Number(node.(*ast.NumberLiteral).Value), nil
	case ast.BOOLEAN_LITERAL:
		return toBoolean(node.(*ast.BooleanLiteral).Value), nil
	case ast.STRING_LITERAL:
		return object.String(node.(*ast.StringLiteral).Value), nil
	case ast.ARRAY_LITERAL:
		return evaluateArrayLiteral(node.(*ast.ArrayLiteral), env)
	case ast.HASH_LITERAL:
		return evaluateHashLiteral(node.(*ast.HashLiteral), env)
	case ast.NULL_LITERAL:
		return NULL, nil
	default:
		panic(fmt.Errorf("unexpected node type: %T", node))
	}
}

func evaluateProgram(node *ast.Program, env *object.Environment) (object.Object, object.Interruption) {
	var result object.Object
	var interrupt object.Interruption
	for _, stmt := range node.Statements {
		result, interrupt = Evaluate(stmt, env)
		if interrupt != nil {
			if interrupt.Type() == object.RETURN_VALUE {
				return interrupt.(*object.ReturnValue).Value, nil
			}
			return nil, interrupt
		}
	}
	return result, nil
}

func evaluateLetDeclaration(node *ast.LetDeclaration, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.Value, env)
	if interrupt != nil {
		return nil, interrupt
	}
	env.Set(node.Name.Value, value)
	return NULL, nil
}

func evaluateReturnStatement(node *ast.ReturnStatement, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.ReturnValue, env)
	if interrupt != nil {
		return nil, interrupt
	}
	return nil, &object.ReturnValue{Value: value}
}

func evaluateExpressionStatement(node *ast.ExpressionStatement, env *object.Environment) (object.Object, object.Interruption) {
	return Evaluate(node.Expression, env)
}

func evaluateBlockStatement(node *ast.BlockStatement, env *object.Environment) (object.Object, object.Interruption) {
	var result object.Object
	var interrupt object.Interruption
	for _, stmt := range node.Statements {
		result, interrupt = Evaluate(stmt, env)
		if interrupt != nil {
			return nil, interrupt
		}
	}
	return result, nil
}

func evaluateUnaryExpression(node *ast.UnaryExpression, env *object.Environment) (object.Object, object.Interruption) {
	right, interrupt := Evaluate(node.Right, env)
	if interrupt != nil {
		return nil, interrupt
	}
	switch node.Operator {
	case "!":
		if isTruthy(right) {
			return FALSE, nil
		}
		return TRUE, nil
	case "-":
		switch right.Type() {
		case object.NUMBER:
			return object.Number(-right.(object.Number)), nil
		}
	}

	return nil, toError("unknown operator: %s%s", node.Operator, right.Type())
}

func evaluateBinaryExpression(node *ast.BinaryExpression, env *object.Environment) (object.Object, object.Interruption) {
	left, interrupt := Evaluate(node.Left, env)
	if interrupt != nil {
		return nil, interrupt
	}

	right, interrupt := Evaluate(node.Right, env)
	if interrupt != nil {
		return nil, interrupt
	}
	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		switch node.Operator {
		case "+":
			return object.Number(left.(object.Number) + right.(object.Number)), nil
		case "-":
			return object.Number(left.(object.Number) - right.(object.Number)), nil
		case "*":
			return object.Number(left.(object.Number) * right.(object.Number)), nil
		case "/":
			return object.Number(left.(object.Number) / right.(object.Number)), nil
		case "<":
			return toBoolean(left.(object.Number) < right.(object.Number)), nil
		case ">":
			return toBoolean(left.(object.Number) > right.(object.Number)), nil
		case "==":
			return toBoolean(left.(object.Number) == right.(object.Number)), nil
		case "!=":
			return toBoolean(left.(object.Number) != right.(object.Number)), nil
		}
	case left.Type() == object.STRING && right.Type() == object.STRING:
		switch node.Operator {
		case "+":
			return object.String(left.(object.String) + right.(object.String)), nil
		case "==":
			return toBoolean(left.(object.String) == right.(object.String)), nil
		case "!=":
			return toBoolean(left.(object.String) != right.(object.String)), nil
		}
	case left.Type() != right.Type():
		return nil, toError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	case node.Operator == "==":
		return toBoolean(left == right), nil
	case node.Operator == "!=":
		return toBoolean(left != right), nil
	}

	return nil, toError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

func evaluateLogicalExpression(node *ast.LogicalExpression, env *object.Environment) (object.Object, object.Interruption) {
	left, interrupt := Evaluate(node.Left, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch node.Operator {
	case "or":
		if isTruthy(left) {
			return left, nil
		}
		return Evaluate(node.Right, env)
	case "and":
		if !isTruthy(left) {
			return left, nil
		}
		return Evaluate(node.Right, env)
	default:
		right, interrupt := Evaluate(node.Right, env)
		if interrupt != nil {
			return nil, interrupt
		}
		return nil, toError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evaluateConditionalExpression(node *ast.ConditionalExpression, env *object.Environment) (object.Object, object.Interruption) {
	condition, interrupt := Evaluate(node.Condition, env)
	if interrupt != nil {
		return nil, interrupt
	}
	if isTruthy(condition) {
		return Evaluate(node.Consequence, env)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative, env)
	}
	return NULL, nil
}

func evaluateAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) (object.Object, object.Interruption) {
	rvalue, interrupt := Evaluate(node.RValue, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch lvalue := node.LValue.(type) {
	case *ast.Identifier:
		if _, found := env.Get(lvalue.Value); found {
			env.Set(lvalue.Value, rvalue)
			return rvalue, nil
		}

		return nil, toError("unknown identifier: %s", lvalue.Value)
	case *ast.SubscriptExpression:
		baseValue, interrupt := Evaluate(lvalue.Base, env)
		if interrupt != nil {
			return nil, interrupt
		}

		subscriptValue, interrupt := Evaluate(lvalue.Subscript, env)
		if interrupt != nil {
			return nil, interrupt
		}

		switch base := baseValue.(type) {
		case *object.Array:
			subscript, ok := subscriptValue.(object.Number)
			if !ok {
				return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
			}

			index, valid := toNativeInt(subscript)
			if !valid || index < 0 {
				return nil, toError("subscript value must be a positive whole number: %s", subscript.Inspect())
			}

			if index >= cap(base.Elements) {
				current := cap(base.Elements)
				base.Elements = slices.Grow(base.Elements, index+1-current)
				base.Elements = base.Elements[:cap(base.Elements)]
				for i := current; i < index; i += 1 {
					base.Elements[i] = NULL
				}
			} else {
				base.Elements = base.Elements[:cap(base.Elements)]
			}
			base.Elements[index] = rvalue
		case *object.Hash:
			key, ok := subscriptValue.(object.Hashable)
			if !ok {
				return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
			}
			base.Pairs[key.HashKey()] = object.HashPair{Key: key, Value: rvalue}
		default:
			return nil, toError("unknown operator: %s[%s]", base.Type(), subscriptValue.Type())
		}
	default:
		panic(fmt.Errorf("unknown lvalue type: %s", node.LValue.Type()))
	}

	return rvalue, nil
}

func evaluateCallExpression(node *ast.CallExpression, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.Callee, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch callee := value.(type) {
	case *object.Function:
		environment := object.NewEnvironment(callee.Closure)
		for i, parameter := range callee.Literal.Parameters {
			value, interrupt := Evaluate(node.Arguments[i], env)
			if interrupt != nil {
				return nil, interrupt
			}
			environment.Set(parameter.Value, value)
		}

		result, interrupt := Evaluate(callee.Literal.Body, environment)
		if interrupt != nil {
			if interrupt.Type() == object.RETURN_VALUE {
				return interrupt.(*object.ReturnValue).Value, nil
			}
			return nil, interrupt
		}
		return result, nil
	case *object.BuiltinFunction:
		args := []object.Object{}
		for _, arg := range node.Arguments {
			value, interrupt := Evaluate(arg, env)
			if interrupt != nil {
				return nil, interrupt
			}
			args = append(args, value)
		}
		return callee.Fn(env, args...)
	default:
		return nil, toError("unknown operator: %s()", callee.Type())
	}
}

func evaluateSubscriptExpression(node *ast.SubscriptExpression, env *object.Environment) (object.Object, object.Interruption) {
	baseValue, interrupt := Evaluate(node.Base, env)
	if interrupt != nil {
		return nil, interrupt
	}

	subscriptValue, interrupt := Evaluate(node.Subscript, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch base := baseValue.(type) {
	case *object.Array:
		subscript, ok := subscriptValue.(object.Number)
		if !ok {
			return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
		}

		index, valid := toNativeInt(subscript)
		if !valid || index < 0 || index >= len(base.Elements) {
			return NULL, nil
		}
		return base.Elements[index], nil
	case *object.Hash:
		key, ok := subscriptValue.(object.Hashable)
		if !ok {
			return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
		}

		value, found := base.Pairs[key.HashKey()]
		if !found {
			return NULL, nil
		}
		return value.Value, nil
	default:
		return nil, toError("unknown operator: %s[%s]", base.Type(), subscriptValue.Type())
	}
}

func evaluateFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) (object.Object, object.Interruption) {
	return &object.Function{Literal: node, Closure: env}, nil
}

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, object.Interruption) {
	if value, found := env.Get(node.Value); found {
		return value, nil
	}

	if builtin, found := builtins[node.Value]; found {
		return builtin, nil
	}

	return nil, toError("unknown identifier: %s", node.Value)
}

func evaluateArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) (object.Object, object.Interruption) {
	elements := []object.Object{}
	for _, element := range node.Elements {
		value, interrupt := Evaluate(element, env)
		if interrupt != nil {
			return nil, interrupt
		}
		elements = append(elements, value)
	}
	return &object.Array{Elements: elements}, nil
}

func evaluateHashLiteral(node *ast.HashLiteral, env *object.Environment) (object.Object, object.Interruption) {
	hash := &object.Hash{Pairs: map[object.HashKey]object.HashPair{}}
	for _, keyNode := range node.Keys {
		key, interrupt := Evaluate(keyNode, env)
		if interrupt != nil {
			return nil, interrupt
		}

		hashable, ok := key.(object.Hashable)
		if !ok {
			return nil, toError("unknown operator: HASH[%s]", key.Type())
		}

		value, interrupt := Evaluate(node.Pairs[keyNode], env)
		if interrupt != nil {
			return nil, interrupt
		}

		hash.Pairs[hashable.HashKey()] = object.HashPair{Key: key, Value: value}
	}
	return hash, nil
}

func toBoolean(value bool) object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func toError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func toNativeInt(subscript object.Number) (int, bool) {
	if subscript == object.Number(math.Trunc(float64(subscript))) {
		return int(subscript), true
	}
	return 0, false
}

func isTruthy(o object.Object) object.Boolean {
	switch {
	case o == FALSE:
		return FALSE
	case o.Type() == object.NUMBER && o.(object.Number) == 0.0:
		return FALSE
	case o.Type() == object.STRING && o.(object.String) == "":
		return FALSE
	case o.Type() == object.ARRAY && len(o.(*object.Array).Elements) == 0:
		return FALSE
	case o.Type() == object.HASH && len(o.(*object.Hash).Pairs) == 0:
		return FALSE
	case o == NULL:
		return FALSE
	default:
		return TRUE
	}
}
