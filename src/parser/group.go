package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)

func (p *MSParser) parseGroupExpression() (ast.ExpNodeI, error) {

	ok, lpar := p.expect(token.LEFT_PAREN)

	// expect '('
	if !ok {
		return nil, p.unexpectedToken(lpar, token.LEFT_PAREN)
	}

	// parse the expression inside the parenthesis
	node, err := p.parseExpression()

	// When encountering an error, return the error
	// and the parser should synchronize this statment.
	if err != nil {
		return node, err
	}

	// We expect a closing parenthesis
	// after the expression
	ok, rpar := p.expect(token.RIGHT_PAREN)

	if !ok {
		msg := fmt.Sprintf("Expected ')' got '%v'", rpar.Type.String())
		err = p.error(msg, rpar.Line, rpar.Col)
	}

	// wrap node in parenthesis
	return &ast.GroupExpNodeS{Node: node, TokenLeft: lpar, TokenRight: rpar}, err
}