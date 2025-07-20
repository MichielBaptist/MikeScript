package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseWhile() (ast.WhileNodeS, error) {

	// 1. parse conditional expression
	cond, err := parser.parseExpression()

	// check for errors
	if err != nil {
		return ast.WhileNodeS{}, err
	}

	// Expect a block statement
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return ast.WhileNodeS{}, err
	}
	block, err := parser.parseBlock()

	// Pack into a while node
	return ast.WhileNodeS{Condition: cond, Body: block}, err

}

func (parser *MSParser) parseBreak(tk token.Token) (ast.BreakNodeS, error) {
	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err := parser.error(msg, tok.Line, tok.Col)
		return ast.BreakNodeS{}, err
	}
	return ast.BreakNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseContinue(tk token.Token) (ast.ContinueNodeS, error) {
	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err := parser.error(msg, tok.Line, tok.Col)
		return ast.ContinueNodeS{}, err
	}
	return ast.ContinueNodeS{Tk: tk}, nil
}
