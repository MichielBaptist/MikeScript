package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseVarDeclaration(declarationType token.Token) (ast.VarDeclNodeS, error) {
	// arg: type of declartion;
	
	// We expect an identifier next so we parse it
	ident, err := parser.parseIdentifier()

	// check for errors
	if err != nil {
		// If we encounter an error parsing the identifier
		// We return the error and empty declaration node
		return ast.VarDeclNodeS{}, err
	}

	// expect a semicolon
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ast.VarDeclNodeS{Identifier: ident, Vartype: declarationType}, err

}