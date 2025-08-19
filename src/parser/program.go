package parser

import (
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)



func (parser *MSParser) parseProgram() (*ast.Program, error) {
	// parses program -> statement *

	statements := []ast.StmtNodeI{}
	var err error
	var stmt ast.StmtNodeI

	for !parser.atend() && parser.peek().Type != token.EOF {

		stmt, err = parser.parseStatement()

		if err == nil {
			statements = append(statements, stmt)
		}

		if parser.pnc {
			parser.synchronize()
		}
		
	}

	return &ast.Program{Statements: statements}, err
}