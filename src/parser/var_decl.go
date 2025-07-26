package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (parser *MSParser) parseVarDeclaration() (ast.VarDeclNodeS, error) {
	// type IDENTIFIER ';'

	var vtype mstype.MSType
	var err error

	vtype, err = parser.parseType()

	// For now, we cannor declare MSOperationTypeS
	// Because they are declared using function keyword
	// if _, ok := vtype.(*mstype.MSOperationTypeS) ; ok {
	// 	msg := "Cannot declare function types with 'var', functions can only be declared with 'function'"
	// 	return ast.VarDeclNodeS{}, parser.error(msg, 0, 0)
	// }

	if err != nil {
		return ast.VarDeclNodeS{}, nil
	}

	// Case 2 & 3: ignore for now
	
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

	return ast.VarDeclNodeS{Identifier: ident, Vartype: vtype}, err

}