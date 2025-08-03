package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (parser *MSParser) parseUnary() (ast.ExpNodeI, error) {

	if ok, op := parser.match(token.MINUS, token.EXCLAMATION, token.EQ); ok {
		right, err := parser.parseUnary()

		if err != nil {
			return right, err
		}

		// Still make a distinction between unary and function calls
		// though they are the same "priority", namely the highest

		if op.Type == token.EQ {
			return &ast.FuncCallNodeS{Op: op, Fun: right}, nil
		} else {
			return &ast.UnaryExpNodeS{Op: op, Node: right}, nil
		}
	}

	return parser.parseAccess()
}

func (parser *MSParser) parseAccess() (ast.ExpNodeI, error) {
	// Parses: primary { '.' IDENTIFIER | '[' expression ']'  }*

	var left ast.ExpNodeI
	var err error

	left, err = parser.parsePrimary()

	if err != nil {
		return left, err
	}
	
	for {

		// match '['
		ok, _ := parser.match(token.LEFT_SQUARE)

		if !ok {
			break;
		}
		
		// parse expression
		index, err := parser.parseExpression()

		if err != nil {
			return &ast.ArrayIndexNodeS{}, err
		}

		fmt.Printf("target: %s\n", left)
		fmt.Printf("index: %s\n", index)
		left = &ast.ArrayIndexNodeS{Target: left, Index: index}

		// Expect closing brace
		if ok, op := parser.expect(token.RIGHT_SQUARE) ; !ok {
			return &ast.ArrayIndexNodeS{}, parser.unexpectedToken(op, token.RIGHT_SQUARE)
		}
	}

	return left, err
}