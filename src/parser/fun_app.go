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
			// >> function application (parameter binding)
			// >>= function application && call

			// TODO: remove tuple alltogether?
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
			default:
				err = parser.error(fmt.Sprintf("Expected a variable, got '%v'", v), op.Line, op.Col)
			}
		case token.EQ_GREATER:
			// => create && assignment

			switch v := right.(type) {
			case *ast.VariableExpNodeS:
				left = &ast.DeclAssignNodeS{Identifier: v, Exp: left}
			default:
				err = parser.error(fmt.Sprintf("Expected a variable, got '%v'", v), op.Line, op.Col)
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

	// nil node might happen when for example calling a
	// function without arguments.
	if (n == nil) {
		return []ast.ExpNodeI{}
	}

	// By default, flatten returns the node wrapped in a slice
	lexpressions := []ast.ExpNodeI{n}
	
	// If the node is a tuple, we need to flatten
	// the left side and append the right side.
	switch node := n.(type) {
	case *ast.BinaryExpNodeS:
		switch node.Op.Type {
		case token.COMMA:
			lexpressions = append(flattenExpNode(node.Left), node.Right)
		}
	}

	return lexpressions
}