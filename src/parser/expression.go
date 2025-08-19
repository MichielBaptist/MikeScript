package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseExpressionStatement() (*ast.ExStmtNodeS, error) {

	xpr, err := parser.parseExpression()

	if err != nil {
		return &ast.ExStmtNodeS{Ex: xpr}, err
	}

	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected '%v' got '%v'", token.SEMICOLON, tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return &ast.ExStmtNodeS{Ex: xpr}, err
}

func (parser *MSParser) parseExpression() (ast.ExpNodeI, error) {
	return parser.parseLor()
}

func (parser *MSParser) parseLor() (ast.ExpNodeI, error) {
	land, ok := parser.parseLand()

	// Error on parsing logical and
	if ok != nil {
		return land, ok
	}

	for {

		// Match ||
		ok, op := parser.match(token.BAR_BAR)

		if !ok {
			break
		}

		right, err := parser.parseLand()
		land = &ast.LogicalExpNodeS{Left: land, Op: op, Right: right}

		if err != nil {
			return land, err
		}
	}

	return land, nil
}

func (parser *MSParser) parseLand() (ast.ExpNodeI, error) {
	comp, ok := parser.parseFuncop()
	if ok != nil {
		return comp, ok	
	}
	for {
		if ok, op := parser.match(token.AMP_AMP); ok {
			right, err := parser.parseFuncop()
			comp = &ast.LogicalExpNodeS{Left: comp, Op: op, Right: right}
			if err != nil {
				return comp, err
			}
		} else {
			break
		}
	}
	return comp, nil
}

func (parser *MSParser) parseTuple() (ast.ExpNodeI, error) {

	var node ast.ExpNodeI
	var err error
	var exprs []ast.ExpNodeI
	
	node, err = parser.parseEquality()

	if err != nil {
		return node, err
	}

	exprs = append(exprs, node)

	for {
		if ok, _ := parser.match(token.COMMA); ok {
			// Keep parsing as long as we see commas
			
			right, err := parser.parseEquality()

			if err != nil {
				return node, err
			}

			exprs = append(exprs, right)

		} else {
			break
		}
	}

	if len(exprs) == 1 {
		return exprs[0], nil
	} else {
		return &ast.TupleNodeS{Expressions: exprs}, nil
	}

}

func (parser *MSParser) parseEquality() (ast.ExpNodeI, error) {
	// 1. parse comparison
	node, err := parser.parseComp()

	if err != nil {
		return node, err
	}

	// 2. While we see a '==' or '!='
	// we keep building the tree from
	// left to right
	for {
		if ok, op := parser.match(token.EQ_EQ, token.EXCLAMATION_EQ); ok {

			right, err := parser.parseComp()

			is_neq := op.Type == token.EXCLAMATION_EQ

			eqop := token.Token{Type: token.EQ_EQ, Lexeme: "==", Col: op.Col, Line: op.Line}
			node = &ast.BinaryExpNodeS{Left: node, Op: eqop, Right: right}

			if is_neq {
				neg := token.Token{Type: token.EXCLAMATION, Lexeme: "!", Col: op.Col, Line: op.Line}
				node = &ast.UnaryExpNodeS{Op: neg, Node: node}
			}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseComp() (ast.ExpNodeI, error) {

	node, err := parser.parseTerm()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.LESS, token.GREATER, token.LESS_EQ, token.GREATER_EQ); ok {
			right, err := parser.parseTerm()

			node = &ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}
		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseTerm() (ast.ExpNodeI, error) {

	node, err := parser.parseFactor()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.PLUS, token.MINUS); ok {
			right, err := parser.parseFactor()

			// If the operator is minus we just wrap the right
			// node in a unary node with - operator
			// x - y => x + (-y)
			if op.Type == token.MINUS {
				right = &ast.UnaryExpNodeS{Op: op, Node: right}        // x (-y)
				op = token.Token{Type: token.PLUS, Lexeme: "+", Col: op.Col, Line: op.Line} // x + (-y)
			}
			node = &ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseFactor() (ast.ExpNodeI, error) {

	node, err := parser.parseUnary()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.MULT, token.SLASH, token.PERCENT); ok {
			right, err := parser.parseUnary()
			node = &ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}



func (parser *MSParser) parseIdentifier() (*ast.VariableExpNodeS, error) {

	// Handles "x", "f", ...
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return &ast.VariableExpNodeS{Name: id}, nil
	}

	// // Handles "(x)", "(f)", "((x))", ...
	// if ok, _ := parser.match(token.LEFT_PAREN); ok {

	// 	// Recursively parse the identifier when 
	// 	// We encounter a left parenthesis
	// 	ident, err := parser.parseIdentifier()

	// 	// On error, return the error and the identifier
	// 	if err != nil {
	// 		return ident, err
	// 	}

	// 	// Ensure we have a closing parenthesis
	// 	if ok, rp := parser.expect(token.RIGHT_PAREN); !ok {
	// 		msg := fmt.Sprintf("Expected ')' got '%v'", rp.Type.String())
	// 		err = parser.error(msg, rp.Line, rp.Col)
	// 	}

	// 	return ident, err
	// }

	tok := parser.peek()
	msg := fmt.Sprintf("Expected identifier got '%v': '%v'", tok.Type.String(), tok.Lexeme)
	err := parser.error(msg, tok.Line, tok.Col)

	return &ast.VariableExpNodeS{}, err
}


