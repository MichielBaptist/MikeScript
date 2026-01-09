package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (parser *MSParser) parseVarDeclaration() (*ast.VarDeclNodeS, error) {
	// type IDENTIFIER ';'

	var vtype mstype.MSType
	var err error

	vtype, err = parser.parseType()

	// fmt.Printf("%p -- %+v\n", vtype, vtype)

	// For now, we cannor declare MSOperationTypeS
	// Because they are declared using function keyword
	// if _, ok := vtype.(*mstype.MSOperationTypeS) ; ok {
	// 	msg := "Cannot declare function types with 'var', functions can only be declared with 'function'"
	// 	return ast.VarDeclNodeS{}, parser.error(msg, 0, 0)
	// }

	if err != nil {
		return nil, err
	}

	// Case 2 & 3: ignore for now
	
	ident, err := parser.parseIdentifier()

	if err != nil {
		return nil, err
	}

	// expect a semicolon
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return &ast.VarDeclNodeS{Identifier: ident, Vartype: vtype}, err

}