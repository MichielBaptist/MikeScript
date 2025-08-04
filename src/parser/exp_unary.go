package parser

import (
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
