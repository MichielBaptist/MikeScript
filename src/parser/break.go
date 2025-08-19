package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseBreak(tk token.Token) (*ast.BreakNodeS, error) {

	// Check if in a loop context
	if !parser.inContext(LOOP) {
		msg := fmt.Sprintf("Connot use '%s' outside of loop contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return &ast.BreakNodeS{}, err
	}

	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		return &ast.BreakNodeS{}, parser.unexpectedToken(tok, token.SEMICOLON)
	}
	return &ast.BreakNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseContinue(tk token.Token) (*ast.ContinueNodeS, error) {

	// Check if in a loop context
	if !parser.inContext(LOOP) {
		msg := fmt.Sprintf("Connot use '%s' outside of loop contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return &ast.ContinueNodeS{}, err
	}

	// Check if semicolon
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		return &ast.ContinueNodeS{}, parser.unexpectedToken(tok, token.SEMICOLON)
	}

	return &ast.ContinueNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseReturn(tk token.Token) (*ast.ReturnNodeS, error) {

	// Check if in a function context
	if !parser.inContext(FUNCTION) {
		msg := fmt.Sprintf("Connot use '%s' outside of function contexts", tk.Lexeme)
		err := parser.error(msg, tk.Line, tk.Col)
		return nil, err
	}

	// Check if semicolon
	var val ast.ExpNodeI
	var err error

	// set val to nothing, if no return value provided, this is returned
	val = &ast.LiteralExpNodeS{Tk: token.Token{Type: token.NOTHING_TYPE, Lexeme: "nothing"}}

	if ok, _ := parser.expect(token.SEMICOLON) ; !ok {

		val, err = parser.parseExpression()

		if err != nil {
			return nil, err
		}

		if ok, tk := parser.expect(token.SEMICOLON) ; !ok {
			return nil, parser.unexpectedToken(tk, token.SEMICOLON)
		}

	}

	return &ast.ReturnNodeS{Node: val}, err
}