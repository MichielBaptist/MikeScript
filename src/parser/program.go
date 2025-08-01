package parser

import (
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)



func (parser *MSParser) parseProgram() (*ast.Program, error) {
	// parses program -> statement *

	// as long as we haven't reached the end of the tokens
	// we keep parsing statements.
	statements := []ast.StmtNodeI{}
	var err error
	var stmt ast.StmtNodeI

	for !parser.atend() && parser.peek().Type != token.EOF {

		// parse next statement and check for error.
		// Returns stmt is a pointer, not struct.
		stmt, err = parser.parseStatement()

		// Only add to the statements if there was no error
		// Else we just synchronize and continue
		if err == nil {
			statements = append(statements, stmt)
		}

		// check if we are in panic mode.
		// If we are, we need to synchronize
		if parser.pnc {
			parser.synchronize()
		}
		
	}

	// Returns all successfully parsed statements
	// and the last error encountered.
	return &ast.Program{Statements: statements}, err
}