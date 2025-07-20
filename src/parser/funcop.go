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

		// Check if we match '>>' or '->'
		ok, op := parser.match(token.GREATER_GREATER, token.MINUS_GREAT)
		
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

			// TODO: make sure right side is VariableExpNodeS
			switch v := right.(type) {
			case ast.VariableExpNodeS:
				left = ast.FuncAppNodeS{Args: lexpressions, Fun: right}
			default:
				return left, parser.error(fmt.Sprintf("Expected a function identifer, got '%v'", v), op.Line, op.Col)
			}

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
		}

		// check for errors
		if err != nil {
			return right, err
		}

	}

	return left, err
}