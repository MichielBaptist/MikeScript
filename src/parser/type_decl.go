package parser

import (
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (p *MSParser) parseTypeDeclaration() (ast.StmtNodeI, error) {

	// Parses: "type" type identifier ";"
	// Parses: "type" "struct" identifier '{' ... '}'

	var node ast.StmtNodeI
	var err error

	if ok, _ := p.match(token.STRUCT) ; ok {
		node, err = p.parseStructDeclaration()
	} else {
		node, err = p.parseTypedefStatement()
	}

	if err != nil {
		return nil, err
	}

	return node, nil

}

func (p *MSParser) parseTypedefStatement() (*ast.TypeDefStatementS, error) {
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


	return &ast.TypeDefStatementS{Tname: v, Type: t}, nil
}
