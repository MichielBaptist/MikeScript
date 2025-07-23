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
	if !parser.checkType(token.GREATER_GREATER) {
		left, err = parser.parseTuple()
	}

	// check for errors
	if err != nil {
		return left, err
	}

	for {

		// Check if we match '>>' or '->' or '=>'
		ok, op := parser.match(token.GREATER_GREATER, token.MINUS_GREAT, token.EQ_GREATER)
		
		// Check if we matched
		if !ok {
			break
		}

		// Found match, check which one we got
		var right ast.ExpNodeI
		var err error

		// Parse the right side of the '>>'. If this is
		// and identifier, we have either a function application
		// or a variable assignment. The parser cannot know
		// which one it is, we only know this at runtime time.
		switch op.Type {
		case token.GREATER_GREATER:

			// Parse the right side of the function application
			// This should resolve into an identifier.
			right, err = parser.parseTuple()

			// Flatten the left side of the function application
			// into a list of expressions
			lexpressions := flattenExpNode(&left)

			left = ast.FuncAppNodeS{Args: lexpressions, Fun: right}

		case token.MINUS_GREAT:

			// Parse the right side of the assignment
			// This should resolve into:
			// 1. A single variable
			// 2. A tuple of variables (not implemented)
			variable, verr := parser.parseTuple()
			err = verr

			// TODO: add support for tuple assignments
			// TODO:
			switch v := variable.(type) {
			case ast.VariableExpNodeS:
				left = ast.AssignmentNodeS{Identifier: v, Exp: left}
			default:
				return left, parser.error(fmt.Sprintf("Expected a variable, got '%v'", v), op.Line, op.Col)
			}
		case token.EQ_GREATER:
			variable, verr := parser.parseTuple()
			err = verr

			switch v := variable.(type) {
			case ast.VariableExpNodeS:
				left = ast.DeclAssignNodeS{Identifier: v, Exp: left}
			default:
				return left, parser.error(fmt.Sprintf("Expected a variable, got '%v'", v), op.Line, op.Col)
			}


		}

		// check for errors
		if err != nil {
			return right, err
		}

	}

	return left, err
}

func flattenExpNode(n *ast.ExpNodeI) []ast.ExpNodeI {

	// nil node might happen when for example calling a
	// function without arguments.
	if (*n == nil) {
		return []ast.ExpNodeI{}
	}

	// By default, flatten returns the node wrapped in a slice
	lexpressions := []ast.ExpNodeI{*n}
	
	// If the node is a tuple, we need to flatten
	// the left side and append the right side.
	switch node := (*n).(type) {
	case ast.BinaryExpNodeS:
		switch node.Op.Type {
		case token.COMMA:
			lexpressions = append(flattenExpNode(&node.Left), node.Right)
		}
	}

	return lexpressions
}