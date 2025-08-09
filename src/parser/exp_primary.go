package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (parser *MSParser) parsePrimary() (ast.ExpNodeI, error) {
	// parses:
	// 1. literal
	// 2. IDENTIFIER
	// 3. IDENTIFIER '{' ... '}'
	// 4. '(' expr ')'
	// 5. '[' exp ']' type '{' exp ? {',' exp}* '}'

	var err error = nil

	// 1. Literal
	if ok, tok := parser.match(token.NUMBER_INT, token.NUMBER_FLOAT, token.STRING, token.TRUE, token.FALSE); ok {
		return &ast.LiteralExpNodeS{Tk: tok}, err
	}

	// 2. IDENTIFIER
	// 3. IDENTIFIER '{' ... '}'
	if ok, _ := parser.lookahead(token.IDENTIFIER); ok {
		return parser.parseStructConstructor()
	}

	// 4. '(' expr ')'
	if ok, _ := parser.lookahead(token.LEFT_PAREN); ok {
		return parser.parseGroupExpression()
	}

	// 5. '[' exp ']' type '{' exp ? {',' exp}* '}'
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
