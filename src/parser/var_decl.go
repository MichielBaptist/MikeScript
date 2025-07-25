package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	"mikescript/src/mstype"
	token "mikescript/src/token"
)

func (parser *MSParser) parseVarDeclaration(tp token.Token) (ast.VarDeclNodeS, error) {
	// arg: type of declartion;
	// Need to consider 3 cases:
	// 1. basic types 'int', 'bool', 'float', 'string'
	// 2. composite types (type, type)
	// 3. function types (type, type, type -> type)

	var vtype mstype.MSType
	var err error

	// case 1:
	if tp.Type == token.INT_TYPE || tp.Type == token.FLOAT_TYPE || tp.Type == token.STRING_TYPE || tp.Type == token.BOOLEAN_TYPE {
		
		vtype, err = mstype.TokenToType(&tp)

		if err != nil {
			return ast.VarDeclNodeS{}, err
		}
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