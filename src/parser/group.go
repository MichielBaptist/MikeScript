package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)

func (p *MSParser) parseGroupExpression() (ast.ExpNodeI, error) {

	ok, lpar := p.expect(token.LEFT_PAREN)

	if !ok {
		return nil, p.unexpectedToken(lpar, token.LEFT_PAREN)
	}

	node, err := p.parseExpression()

	if err != nil {
		return node, err
	}

	ok, rpar := p.expect(token.RIGHT_PAREN)

	if !ok {
		msg := fmt.Sprintf("Expected ')' got '%v'", rpar.Type.String())
		err = p.error(msg, rpar.Line, rpar.Col)
	}

	return &ast.GroupExpNodeS{Node: node, TokenLeft: lpar, TokenRight: rpar}, err
}