package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseWhile() (*ast.WhileNodeS, error) {

	cond, err := parser.parseExpression()

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

	return &ast.WhileNodeS{Condition: cond, Body: block}, err

}

