package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseIf() (*ast.IfNodeS, error) {

	// Parse conditional expression
	cond, err := parser.parseExpression()
	if err != nil {
		return &ast.IfNodeS{}, err
	}

	// Expect a block statement so we check for opening brace.
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return &ast.IfNodeS{}, err
	}
	stmt, err := parser.parseBlock()
	if err != nil {
		return &ast.IfNodeS{}, err
	}

	// Check for if statement else branch
	// and parse it if it exists.
	if ok, _ := parser.match(token.ELSE); ok {

		// Expect a block statement
		if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
			msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
			err = parser.error(msg, tok.Line, tok.Col)
			return &ast.IfNodeS{}, err
		}
		elsestmt, err := parser.parseBlock()
		if err != nil {
			return &ast.IfNodeS{}, err
		}

		return &ast.IfNodeS{Condition: cond, ThenStmt: stmt, ElseStmt: elsestmt}, err
	}

	return &ast.IfNodeS{Condition: cond, ThenStmt: stmt, ElseStmt: nil}, err

}