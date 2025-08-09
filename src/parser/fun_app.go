package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseFuncop() (ast.ExpNodeI, error) {

	// Parses:
	// >> tuple
	// tuple >> tuple
	// tuple >> tuple >> ...

	var left ast.ExpNodeI
	var err error

	// parse the first expression, this is either
	// a comma separated list of expressions or empty (nil)
	left, err = parser.parseTuple()

	// check for errors
	if err != nil {
		return left, err
	}

	for {

		ok, op := parser.match(
			token.GREATER_GREATER,		// >>  param binding
			token.MINUS_GREAT,			// ->  assignment
			token.EQ_GREATER,			// =>  declaration & assignment
			token.GREATER_GREATER_EQ,	// >>= binding & call
		)
		
		// Check if we matched
		if !ok {
			break
		}

		var right ast.ExpNodeI
		var err error

		// parse right
		right, err = parser.parseTuple()

		// stop on fail
		if err != nil {
			return left, err
		}

		switch op.Type {
		case token.GREATER_GREATER, token.GREATER_GREATER_EQ:
			// >>  function application (parameter binding)
			// >>= function application && call. 'e1, e2 >>= f' is syntactic sugar for '=(e1, e2 ... >> f)'

			lexpressions := flattenExpNode(left)

			// Function application
			left = &ast.FuncAppNodeS{Args: lexpressions, Fun: right}

			// also wrap with call?
			if op.Type == token.GREATER_GREATER_EQ {
				left = &ast.FuncCallNodeS{Op: op, Fun: left}
			}

		case token.MINUS_GREAT:
			// -> assignment

			switch v := right.(type) {
			case *ast.VariableExpNodeS:
				left = &ast.AssignmentNodeS{Identifier: v, Exp: left}
			case *ast.ArrayIndexNodeS:
				left = &ast.ArrayAssignmentNodeS{Target: v.Target, Index: v.Index, Value: left}
			case *ast.FieldAccessNodeS:
				left = &ast.FieldAssignmentNode{Target: v.Target, Field: v.Field, Value: left}
			default:
				err = parser.error(fmt.Sprintf("Expected an assignable target, got '%v'", v), op.Line, op.Col)
			}
		case token.EQ_GREATER:
			// => create && assignment

			switch v := right.(type) {
			case *ast.VariableExpNodeS:
				left = &ast.DeclAssignNodeS{Identifier: v, Exp: left}
			default:
				err = parser.error(fmt.Sprintf("Expected an assignable target, got '%v'", v), op.Line, op.Col)
			}
		}

		// Break loop if we see an error
		if err != nil {
			break
		}
	}

	return left, err
}

func flattenExpNode(n ast.ExpNodeI) []ast.ExpNodeI {
	switch t := n.(type) {
	case *ast.TupleNodeS: 	return t.Expressions
	case nil:				return []ast.ExpNodeI{}
	default: 				return []ast.ExpNodeI{n}
	}
}