package parser

import (
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

		// look for '.' or '['
		ok, tok := parser.lookahead(token.LEFT_SQUARE, token.DOT)

		// no more access tokens
		if !ok {
			break
		}

		switch tok.Type {
		case token.LEFT_SQUARE:	left, err = parser.parseIndexing(left)
		case token.DOT:			left, err = parser.parseStructFieldAccess(left)
		}

		if err != nil {
			return nil, err
		}
	}

	return left, err
}

func (p *MSParser) parseIndexing(target ast.ExpNodeI) (*ast.ArrayIndexNodeS, error) {
	// parses: primary '[' exp ']'
	// target (primary) is already parsed and given

	// '['
	ok, tok := p.match(token.LEFT_SQUARE)

	if !ok {
		return nil, p.unexpectedToken(tok, token.LEFT_SQUARE)
	}

	// exp
	index, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	// ']'
	if ok, op := p.expect(token.RIGHT_SQUARE) ; !ok {
		return nil, p.unexpectedToken(op, token.RIGHT_SQUARE)
	}

	return &ast.ArrayIndexNodeS{Target: target, Index: index}, nil

}

func (p *MSParser) parseStructFieldAccess(target ast.ExpNodeI) (*ast.FieldAccessNodeS, error) {
	// parses: primary '.' IDENTIFIER
	// target (primary) is parsed and given

	// '.'
	if ok, tok := p.match(token.DOT) ; !ok {
		return nil, p.unexpectedToken(tok, token.DOT)
	}

	// IDENTIFIER
	fieldName, err := p.parseIdentifier()

	if err != nil {
		return nil, err
	}

	return &ast.FieldAccessNodeS{Target: target, Field: fieldName}, nil

}

func (p *MSParser) parseArrayExpression() (ast.ExpNodeI, error) {
	// 1) exp? ']' type '{' {expression ','} * '}' --> array constructor
	// or 
	// 2) exp? '..' exp? ']'

	var n ast.ExpNodeI
	var err error

	// Check if there is an expression between '[' exp ']' or '[' '..'
	if ok, _ := p.lookahead(token.RIGHT_SQUARE, token.DOT_DOT) ; !ok {
		n, err = p.parseExpression()
	}

	if err != nil {
		return nil, err
	}

	// Check for '..' or ']'
	if ok, _ := p.match(token.DOT_DOT) ; ok {
		return p.parseRangeConstructor(n)
	} else if ok, _ := p.match(token.RIGHT_SQUARE) ; ok {
		return p.parseArrayConstructor(n)
	}

	return nil , err
}

func (p *MSParser) parseRangeConstructor(start ast.ExpNodeI) (ast.ExpNodeI, error) {
	// parses:  exp? ']'

	var to ast.ExpNodeI = nil
	var err error

	// expect an expression
	to, err = p.parseExpression()
	
	// exit on error
	if err != nil {
		return nil, err
	}

	// Expect ']'
	if ok, tk := p.match(token.RIGHT_SQUARE) ; !ok {
		return nil, p.unexpectedToken(tk, token.RIGHT_SQUARE)
	}

	// If start is nil, we add a 0 start
	if start == nil {
		start = &ast.LiteralExpNodeS{Tk: token.Token{Type: token.NUMBER_INT, Lexeme: "0", Line: 0, Col: 0}}
	}

	return &ast.RangeConstructorNodeS{From: start, To: to}, err
}

func (p *MSParser) parseArrayConstructor(n ast.ExpNodeI) (ast.ExpNodeI, error) {

	// Need type
	atype, err := p.parseType()

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