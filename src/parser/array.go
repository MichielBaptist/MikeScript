package parser

import (
	"fmt"
	"mikescript/src/ast"
	token "mikescript/src/token"
)


func (parser *MSParser) parseAccess() (ast.ExpNodeI, error) {
	// Parses: primary { '.' IDENTIFIER | '[' expression ']'  }*

	var left ast.ExpNodeI
	var err error

	left, err = parser.parsePrimary()

	if err != nil {
		return left, err
	}
	
	for {

		// match '['
		ok, _ := parser.match(token.LEFT_SQUARE)

		if !ok {
			break;
		}
		
		// parse expression
		index, err := parser.parseExpression()

		if err != nil {
			return &ast.ArrayIndexNodeS{}, err
		}

		fmt.Printf("target: %s\n", left)
		fmt.Printf("index: %s\n", index)
		left = &ast.ArrayIndexNodeS{Target: left, Index: index}

		// Expect closing brace
		if ok, op := parser.expect(token.RIGHT_SQUARE) ; !ok {
			return &ast.ArrayIndexNodeS{}, parser.unexpectedToken(op, token.RIGHT_SQUARE)
		}
	}

	return left, err
}


func (p *MSParser) parseArrayConstructor() (ast.ExpNodeI, error) {
	// ']' type '{' {expression ','} * '}' --> array constructor

	// Need ']'
	if ok, tk := p.match(token.RIGHT_SQUARE) ; !ok {
		return nil, p.unexpectedToken(tk, token.RIGHT_SQUARE)
	}

	// Parse type
	atype, err := p.parseType()

	fmt.Printf("Type: %+v\n", atype)

	if err != nil {
		return nil, err
	}

	// Need '{'
	if ok, tk := p.match(token.LEFT_BRACE) ; !ok {
		return nil, p.unexpectedToken(tk, token.LEFT_BRACE)
	}

	// check for empty constructor
	if ok, _ := p.match(token.RIGHT_BRACE) ; ok {
		vals := make([]ast.ExpNodeI, 0)
		return &ast.ArrayConstructorNodeS{Type: atype, Vals: vals}, nil
	}

	// Parse expressions
	tuple, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	// Flatten tuple expression into list of expressions
	exprs := flattenExpNode(tuple)

	// Need '}'
	if ok, tk := p.match(token.RIGHT_BRACE) ; !ok {
		return nil, p.unexpectedToken(tk, token.RIGHT_BRACE)
	}

	return &ast.ArrayConstructorNodeS{Type: atype, Vals: exprs}, nil
}