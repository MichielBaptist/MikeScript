package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseExpressionStatement() (ast.ExStmtNodeS, error) {

	// Parse the expression.
	xpr, err := parser.parseExpression()

	// Check for errors
	if err != nil {
		return ast.ExStmtNodeS{Ex: xpr}, err
	}

	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ast.ExStmtNodeS{Ex: xpr}, err
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

		// Match
		ok, op := parser.match(token.BAR_BAR)

		// No match
		if !ok {
			break
		}
		
		// Match
		right, err := parser.parseLand()
		land = ast.LogicalExpNodeS{Left: land, Op: op, Right: right}

		// check for errors
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
			comp = ast.LogicalExpNodeS{Left: comp, Op: op, Right: right}
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
	
	// parse the first expression
	node, err := parser.parseEquality()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.COMMA); ok {
			right, err := parser.parseEquality()
			node = ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

			// check for errors
			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
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
			// found next "==" or "!="
			right, err := parser.parseComp()

			// check if neq or eq
			is_neq := op.Type == token.EXCLAMATION_EQ

			// build binary node
			eqop := token.Token{Type: token.EQ_EQ, Lexeme: "==", Col: op.Col, Line: op.Line}
			node = ast.BinaryExpNodeS{Left: node, Op: eqop, Right: right}

			// If the operator was neq, we wrap the binary node
			if is_neq {
				neg := token.Token{Type: token.EXCLAMATION, Lexeme: "!", Col: op.Col, Line: op.Line}
				node = ast.UnaryExpNodeS{Op: neg, Node: node}
			}

			if err != nil {
				return node, err
			}

		} else {
			// Next token is not "==" or "!="
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

			node = ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

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
				right = ast.UnaryExpNodeS{Op: op, Node: right}        // x (-y)
				op = token.Token{Type: token.PLUS, Lexeme: "+", Col: op.Col, Line: op.Line} // x + (-y)
			}
			node = ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

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
			node = ast.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseUnary() (ast.ExpNodeI, error) {

	if ok, op := parser.match(token.MINUS, token.EXCLAMATION); ok {
		right, err := parser.parseUnary()

		if err != nil {
			return right, err
		}

		return ast.UnaryExpNodeS{Op: op, Node: right}, nil
	}

	return parser.parsePrimary()
}

func (parser *MSParser) parsePrimary() (ast.ExpNodeI, error) {

	var err error = nil

	// matches a primary expression
	if ok, tok := parser.match(token.NUMBER_INT, token.NUMBER_FLOAT, token.STRING, token.TRUE, token.FALSE); ok {
		return ast.LiteralExpNodeS{Tk: tok}, err
	}

	// matches an identifier
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return ast.VariableExpNodeS{Name: id}, err
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
		return ast.GroupExpNodeS{Node: node, TokenLeft: lpar, TokenRight: rpar}, err
	}

	// If we reach this point, we couldn't match any
	// of the primary expressions, so we need to return an error.
	tok := parser.peek()
	msg := fmt.Sprintf("Expected primary expression got '%v'", tok.Type.String())
	err = parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return nil, err
}

func (parser *MSParser) parseIdentifier() (ast.VariableExpNodeS, error) {

	// Handles "x", "f", ...
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return ast.VariableExpNodeS{Name: id}, nil
	}

	// Handles "(x)", "(f)", "((x))", ...
	if ok, _ := parser.match(token.LEFT_PAREN); ok {

		// Recursively parse the identifier when 
		// We encounter a left parenthesis
		ident, err := parser.parseIdentifier()

		// On error, return the error and the identifier
		if err != nil {
			return ident, err
		}

		// Ensure we have a closing parenthesis
		if ok, rp := parser.expect(token.RIGHT_PAREN); !ok {
			msg := fmt.Sprintf("Expected ')' got '%v'", rp.Type.String())
			err = parser.error(msg, rp.Line, rp.Col)
		}

		return ident, err
	}

	tok := parser.peek()
	msg := fmt.Sprintf("Expected identifier got '%v'", tok.Type.String())
	err := parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return ast.VariableExpNodeS{}, err
}


func flattenExpNode(n *ast.ExpNodeI) []ast.ExpNodeI {

	// nil node might happen when for example calling a
	// function without arguments.
	if (*n == nil) {
		return []ast.ExpNodeI{}
	}

	// By default, flatten returns the node wrapped in a slice
	lexpressions := []ast.ExpNodeI{*n}
	
	// If the node is a tuple, we need to flatten
	// the left side and append the right side.
	switch node := (*n).(type) {
	case ast.BinaryExpNodeS:
		switch node.Op.Type {
		case token.COMMA:
			lexpressions = append(flattenExpNode(&node.Left), node.Right)
		}
	}

	return lexpressions
}