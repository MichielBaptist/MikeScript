package parser

import (
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (parser *MSParser) parseUnary() (ast.ExpNodeI, error) {

	if ok, op := parser.match(token.MINUS, token.EXCLAMATION, token.EQ, token.DOT_EQ, token.MULT); ok {
		right, err := parser.parseUnary()

		if err != nil {
			return right, err
		}

		// Still make a distinction between unary and function calls
		// though they are the same "priority", namely the highest

		// Note: .= a, b, c; means =a, =b, =c

		switch op.Type {
		case token.EQ: 		return &ast.FuncCallNodeS{Op: op, Fun: right}, nil
		case token.DOT_EQ:	return &ast.IterableFuncCallNodeS{Op: op, Fun: right}, nil
		case token.MULT:	return &ast.StarredExpNodeS{Node: right}, nil
		default: 			return &ast.UnaryExpNodeS{Op: op, Node: right}, nil
		}
	}

	return parser.parseAccess()
}
