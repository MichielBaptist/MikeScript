package parser

import (
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (p *MSParser) parseTypeDeclaration() (*ast.TypeDeclarationNode, error) {

	// Parses: "type" type identifier ";"

	// parse type
	t, err := p.parseType()

	if err != nil {
		return nil, err
	}

	// parse identifier
	v, err := p.parseIdentifier()

	if err != nil {
		return nil, err
	}

	// expect ';'
	if ok, tok := p.expect(token.SEMICOLON) ; !ok {
		return nil, p.unexpectedToken(tok, token.SEMICOLON)
	}

	return &ast.TypeDeclarationNode{Tname: v, Type: t}, nil
}