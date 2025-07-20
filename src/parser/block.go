package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseBlock() (ast.BlockNodeS, error) {
	// This function expects that a '{' was already matched before
	// but WILL consume the closing '}'.

	// Parse the block
	stmts := []ast.StmtNodeI{}
	var err error
	var stmt ast.StmtNodeI

	// As long as we haven't reached the end of the block
	// we keep parsing statements.
	for !parser.atend() && parser.peek().Type != token.RIGHT_BRACE {

		// continue
		if ok, tk := parser.match(token.CONTINUE); ok {
			stmt, err = parser.parseContinue(tk)
		} else if ok, tk := parser.match(token.BREAK); ok {
			stmt, err = parser.parseBreak(tk)
		} else {
			stmt, err = parser.parseStatement()
		}

		// Only add to the statements if there was no error
		// Else we just synchronize and continue
		if err == nil {
			stmts = append(stmts, stmt)
		} else {
			// Provide the statements we were able to parse.
			return ast.BlockNodeS{Statements: stmts}, err
		}
	}

	// Expect a closing brace
	if ok, tok := parser.expect(token.RIGHT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '}' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ast.BlockNodeS{Statements: stmts}, err
}

