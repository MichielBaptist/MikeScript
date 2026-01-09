package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)

func (parser *MSParser) parseFor() (*ast.ForNodeS, error) {

	var err error = nil

	// Parse iterable variable
	iterable, err := parser.parseExpression()

	if err != nil {
		return nil, err
	}

	// Expect '.->' token
	if ok, tok := parser.expect(token.DOT_MINUS_GREAT); !ok {
		msg := "Expected '.->' in for loop"
		err = parser.error(msg, tok.Line, tok.Col)
		return nil, err
	}

	// Parse loop variable
	loop_var, err := parser.parseIdentifier()

	if err != nil {
		return nil, err
	}

	// Expect a block statement so we check for opening brace.
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return nil, err
	}

	block, err := parser.parseBlock()

	if err != nil {
		return nil, err
	}

	return &ast.ForNodeS{Iterable: iterable, LoopVar: loop_var, Body: block}, err

}