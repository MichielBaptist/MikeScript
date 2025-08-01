package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseWhile() (*ast.WhileNodeS, error) {

	// 1. parse conditional expression
	cond, err := parser.parseExpression()

	// check for errors
	if err != nil {
		return &ast.WhileNodeS{}, err
	}

	// Expect a block statement
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return &ast.WhileNodeS{}, err
	}
	parser.enterContext(LOOP)
	block, err := parser.parseBlock()
	ctx := parser.leaveContext()
	if ctx != LOOP {
		_ = []int{}[0] // force error
	}

	// Pack into a while node
	return &ast.WhileNodeS{Condition: cond, Body: block}, err

}

