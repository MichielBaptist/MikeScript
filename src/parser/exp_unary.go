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

		return ast.UnaryExpNodeS{Op: op, Node: right}, nil
	}

	return parser.parsePrimary()
}