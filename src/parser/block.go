package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseBlock() (*ast.BlockNodeS, error) {
	// This function expects that a '{' was already matched before
	// but WILL consume the closing '}'.

	stmts := []ast.StmtNodeI{}
	var err error
	var stmt ast.StmtNodeI

	for !parser.atend() && parser.peek().Type != token.RIGHT_BRACE {

		stmt, err = parser.parseStatement()

		if err == nil {
			stmts = append(stmts, stmt)
		} else {
			return nil, err
		}
	}

	if ok, tok := parser.expect(token.RIGHT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '}' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return &ast.BlockNodeS{Statements: stmts}, err
}

