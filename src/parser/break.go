package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseBreak(tk token.Token) (ast.BreakNodeS, error) {

	// Check if in a loop context
	if !parser.inContext(LOOP) {
		msg := fmt.Sprintf("Connot use '%s' outside of loop contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return ast.BreakNodeS{}, err
	}

	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		return ast.BreakNodeS{}, parser.unexpectedToken(tok, token.SEMICOLON)
	}
	return ast.BreakNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseContinue(tk token.Token) (ast.ContinueNodeS, error) {

	// Check if in a loop context
	if !parser.inContext(LOOP) {
		msg := fmt.Sprintf("Connot use '%s' outside of loop contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return ast.ContinueNodeS{}, err
	}

	// Check if semicolon
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		return ast.ContinueNodeS{}, parser.unexpectedToken(tok, token.SEMICOLON)
	}

	return ast.ContinueNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseReturn(tk token.Token) (ast.ReturnNodeS, error) {

	// Check if in a function context
	if !parser.inContext(FUNCTION) {
		msg := fmt.Sprintf("Connot use '%s' outside of function contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return ast.ReturnNodeS{}, err
	}

	// Check if semicolon
	var val ast.ExpNodeI
	var err error
	if ok, _ := parser.expect(token.SEMICOLON) ; !ok {

		// get return type
		val, err = parser.parseExpression()

		// on error stop
		if err != nil {
			return ast.ReturnNodeS{}, err
		}

		// Need semicolon
		if ok, tk := parser.expect(token.SEMICOLON) ; !ok {
			return ast.ReturnNodeS{Node: val}, parser.unexpectedToken(tk, token.SEMICOLON)
		}

	}

	return ast.ReturnNodeS{Node: val}, err
}