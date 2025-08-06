package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
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
	// exp? ']' type '{' {expression ','} * '}' --> array constructor

	var atype mstype.MSType
	var n ast.ExpNodeI
	var err error

	// Check if there is an expression between '[' exp ']'
	if ok, _ := p.lookahead(token.RIGHT_SQUARE) ; !ok {
		n, err = p.parseExpression()
	}

	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", n)
	println(p.peek().Lexeme)

	// Need ']'
	if ok, tk := p.match(token.RIGHT_SQUARE) ; !ok {
		return nil, p.unexpectedToken(tk, token.RIGHT_SQUARE)
	}

	// Parse type
	atype, err = p.parseType()

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
		println("Closing brace found, empty initializer")
		vals := make([]ast.ExpNodeI, 0)
		return &ast.ArrayConstructorNodeS{Type: atype, Vals: vals, N: n}, nil
	}

	if n != nil {
		msg := "Cannot initialize an array with values if an initialize amount was provided."
		return nil, p.error(msg, 0, 0)
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

	return &ast.ArrayConstructorNodeS{Type: atype, Vals: exprs, N: n}, nil
}