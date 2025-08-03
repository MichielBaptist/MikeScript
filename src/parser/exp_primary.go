package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (parser *MSParser) parsePrimary() (ast.ExpNodeI, error) {

	var err error = nil

	// matches a primary expression
	if ok, tok := parser.match(token.NUMBER_INT, token.NUMBER_FLOAT, token.STRING, token.TRUE, token.FALSE); ok {
		return &ast.LiteralExpNodeS{Tk: tok}, err
	}

	// matches an identifier
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return &ast.VariableExpNodeS{Name: id}, err
	}

	// matches parenthesis
	if ok, lpar := parser.match(token.LEFT_PAREN); ok {

		// parse the expression inside the parenthesis
		node, err := parser.parseExpression()

		// When encountering an error, return the error
		// and the parser should synchronize this statment.
		if err != nil {
			return node, err
		}

		// We expect a closing parenthesis
		// after the expression
		ok, rpar := parser.expect(token.RIGHT_PAREN)

		if !ok {
			msg := fmt.Sprintf("Expected ')' got '%v'", rpar.Type.String())
			err = parser.error(msg, rpar.Line, rpar.Col)
		}

		// wrap node in parenthesis
		return &ast.GroupExpNodeS{Node: node, TokenLeft: lpar, TokenRight: rpar}, err
	}

	// matches '['
	if ok, _ := parser.match(token.LEFT_SQUARE) ; ok {
		return parser.parseArrayConstructor()
	}

	// If we reach this point, we couldn't match any
	// of the primary expressions, so we need to return an error.
	tok := parser.peek()
	msg := fmt.Sprintf("Expected primary expression got '%v'", tok.Type.String())
	err = parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return nil, err
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