package main

import "fmt"

type Parser interface {
	parse(tokens []Token) (ExpNodeI, error)
}

type MSParser struct {

	// src information
	src    string
	tokens []Token

	// parser state
	pos  int     // current position in tokens
	pnc  bool    // panic flag

	Errors []ParserError
}

////////////////////////////////////////////////////////////
// 						Errors
////////////////////////////////////////////////////////////

type ParserError struct {
	// Represents a parser error
	msg  string
	line int
	col  int
}

func (err ParserError) Error() string {
	return fmt.Sprintf("Parsing Error: %v at line %v col %v", err.msg, err.line, err.col)
}

////////////////////////////////////////////////////////////
// 						helpers
////////////////////////////////////////////////////////////

func (parser *MSParser) advance() Token {
	tok := parser.peek()
	parser.pos = parser.pos + 1
	return tok
}

func (parser *MSParser) peek() Token {
	if parser.atend() {
		return Token{UNKNOWN, "UNKNOWN", 0, 0}
	}
	return parser.tokens[parser.pos]
}

func (parser *MSParser) atend() bool {
	return parser.pos >= len(parser.tokens)
}

func (parser *MSParser) checkType(t TokenType) bool {
	return !parser.atend() && parser.peek().Type == t
}

func (parser *MSParser) match(t ...TokenType) (bool, Token) {
	// check if we matched the token types
	// If we did, we advance the position
	for _, tt := range t {
		if parser.checkType(tt) {
			fmt.Println("Matched: ", tt)
			return true, parser.advance()
		}
	}

	// Did not match any of the tokens
	// return false and an empty token
	return false, parser.peek()
}

func (parser *MSParser) expect(t TokenType) (bool, Token) {
	return parser.match(t)
}

////////////////////////////////////////////////////////////
// 						Error handling
////////////////////////////////////////////////////////////

func (parser *MSParser) panic() {
	parser.pnc = true
}

func (parser *MSParser) error(msg string, line, col int) error {
	err := ParserError{msg, line, col}
	parser.Errors = append(parser.Errors, err)
	parser.panic()
	return err
}


func (parser *MSParser) synchronize() {

	tok := parser.peek()

	// Synchronize the parser
	// by advancing the position
	for tok.Type != SEMICOLON && !parser.atend() && tok.Type != EOF {
		tok = parser.advance()
	}

	fmt.Println("Synchronized to: ", tok)

	if tok.Type != EOF {
		parser.advance()
	}
	
}

////////////////////////////////////////////////////////////
// 						Parse
////////////////////////////////////////////////////////////

func (parser *MSParser) parse(tokens []Token) (Program, error) {

	// parse block
	ast, err := parser.parseProgram()

	// If parsing failed, return the error
	if err != nil {
		return ast, err
	}

	// Expect EOF, otherwise return an error
	if ok, tok := parser.expect(EOF); !ok {
		msg := fmt.Sprintf("Expected 'EOF' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ast, err
}

func (parser *MSParser) parseProgram() (Program, error) {

	// as long as we haven't reached the end of the tokens
	// we keep parsing statements.
	statements := []StmtNodeI{}
	var err error
	var stmt StmtNodeI

	for !parser.atend() && parser.peek().Type != EOF {

		// parse next statement and check for error
		stmt, err = parser.parseDeclaration()

		// Only add to the statements if there was no error
		// Else we just synchronize and continue
		if err == nil {
			statements = append(statements, stmt)
		}

		// check if we are in panic mode.
		// If we are, we need to synchronize
		if parser.pnc {
			parser.synchronize()
		}
		
	}

	// Returns all successfully parsed statements
	// and the last error encountered.
	return Program{statements}, err
}

func (parser *MSParser) parseDeclaration() (StmtNodeI, error) {

	// Found a declaration, parse it
	if ok, tok := parser.match(INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOLEAN_TYPE); ok {
		return parser.parseVarDeclaration(tok)
	}

	// fall through to statement parsing
	return parser.parseStatement()
}

func (parser *MSParser) parseVarDeclaration(declarationType Token) (DeclarationNodeS, error) {
	
	// We expect an identifier next so we parse it
	ident, err := parser.parseIdentifier()

	// check for errors
	if err != nil {
		// If we encounter an error parsing the identifier
		// We return the error and empty declaration node
		return DeclarationNodeS{}, err
	}

	// expect a semicolon
	if ok, tok := parser.expect(SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return DeclarationNodeS{ident, declarationType}, err

}

func (parser *MSParser) parseStatement() (StmtNodeI, error){
	// parse expression statement
	exst, err := parser.parseExpressionStatement()

	// Wrap the expression statement in a statement node
	// and return the error if any.
	return exst, err
}

func (parser *MSParser) parseExpressionStatement() (ExStmtNodeS, error) {

	// Parse the expression.
	xpr, err := parser.parseExpression()

	// Check for errors
	if err != nil {
		return ExStmtNodeS{xpr}, err
	}

	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ExStmtNodeS{xpr}, err
}

func (parser *MSParser) parseExpression() (ExpNodeI, error) {
	return parser.parseFuncop()
}

func (parser *MSParser) parseFuncop() (ExpNodeI, error) {

	var left ExpNodeI
	var err error

	// Check if we have an argument-less function
	// application. If we do, we return the function
	// application node.
	if ok, id := parser.match(GREATER_GREATER); ok {
		msg := "Function call with no arguments is not implemented yet"
		return nil, ParserError{msg, id.Line, id.Col}
	}
	
	// parse the first expression, this is either
	// a comma separated list of expressions or empty (nil)
	left, err = parser.parseTuple()

	// check for errors
	if err != nil {
		return left, err
	}

	for {
		if ok, op := parser.match(GREATER_GREATER, MINUS_GREAT); ok {

			var right ExpNodeI
			var err error

			// Parse the right side of the '>>'. If this is
			// and identifier, we have either a function application
			// or a variable assignment. The parser cannot know
			// which one it is, we only know this at runtime time.
			switch op.Type {
			case GREATER_GREATER:

				// Parse the right side of the function application
				// This should resolve into a function. Or a tuple
				// of functions.
				right, err = parser.parseTuple()

				// Flatten left side of the function application
				// into a slice of expressions
				lexpressions := flattenExpNode(&left)

				// Create a function application node
				left = FuncAppNodeS{lexpressions, right}

			case MINUS_GREAT:

				// Parse the right side of the assignment
				// This should resolve into an identifier
				variable, verr := parser.parseIdentifier()
				err = verr

				// Assignment node (also an expression)
				left = AssignmentNodeS{Identifier: variable, Exp: left}
			}

			// check for errors
			if err != nil {
				return right, err
			}

		} else {
			break
		}
	}

	return left, err
}

func (parser *MSParser) parseTuple() (ExpNodeI, error) {
	
	// parse the first expression
	node, err := parser.parseEquality()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(COMMA); ok {
			right, err := parser.parseEquality()
			node = BinaryExpNodeS{node, op, right}

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

func (parser *MSParser) parseEquality() (ExpNodeI, error) {
	// 1. parse comparison
	node, err := parser.parseComp()

	if err != nil {
		return node, err
	}

	// 2. While we see a '==' or '!='
	// we keep building the tree from
	// left to right
	for {
		if ok, op := parser.match(EQ_EQ, EXCLAMATION_EQ); ok {
			// found next "==" or "!="
			right, err := parser.parseComp()

			// check if neq or eq
			is_neq := op.Type == EXCLAMATION_EQ

			// build binary node
			eqop := Token{EQ_EQ, "==", op.Col, op.Line}
			node = BinaryExpNodeS{node, eqop, right}

			// If the operator was neq, we wrap the binary node
			if is_neq {
				neg := Token{EXCLAMATION, "!", op.Col, op.Line}
				node = UnaryExpNodeS{neg, node}
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

func (parser *MSParser) parseComp() (ExpNodeI, error) {

	node, err := parser.parseTerm()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(LESS, GREATER, LESS_EQ, GREATER_EQ); ok {
			right, err := parser.parseTerm()

			node = BinaryExpNodeS{node, op, right}

			if err != nil {
				return node, err
			}
		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseTerm() (ExpNodeI, error) {

	node, err := parser.parseFactor()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(PLUS, MINUS); ok {
			right, err := parser.parseFactor()

			// If the operator is minus we just wrap the right
			// node in a unary node with - operator
			// x - y => x + (-y)
			if op.Type == MINUS {
				right = UnaryExpNodeS{op, right}        // x (-y)
				op = Token{PLUS, "+", op.Col, op.Line} // x + (-y)
			}
			node = BinaryExpNodeS{node, op, right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseFactor() (ExpNodeI, error) {

	node, err := parser.parseUnary()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(MULT, SLASH); ok {
			right, err := parser.parseUnary()
			node = BinaryExpNodeS{node, op, right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseUnary() (ExpNodeI, error) {

	if ok, op := parser.match(MINUS, EXCLAMATION); ok {
		right, err := parser.parseUnary()

		if err != nil {
			return right, err
		}

		return UnaryExpNodeS{op, right}, nil
	}

	return parser.parsePrimary()
}

func (parser *MSParser) parsePrimary() (ExpNodeI, error) {

	var err error = nil

	// matches a primary expression
	if ok, tok := parser.match(NUMBER_INT, NUMBER_FLOAT, STRING, TRUE, FALSE); ok {
		return LiteralExpNodeS{tok}, err
	}

	// matches an identifier
	if ok, id := parser.match(IDENTIFIER); ok {
		return VariableExpNodeS{id}, err
	}

	// matches parenthesis
	if ok, lpar := parser.match(LEFT_PAREN); ok {

		// parse the expression inside the parenthesis
		node, err := parser.parseExpression()

		// When encountering an error, return the error
		// and the parser should synchronize this statment.
		if err != nil {
			return node, err
		}

		// We expect a closing parenthesis
		// after the expression
		ok, rpar := parser.expect(RIGHT_PAREN)

		if !ok {
			msg := fmt.Sprintf("Expected ')' got '%v'", rpar.Type.String())
			err = parser.error(msg, rpar.Line, rpar.Col)
		}

		// wrap node in parenthesis
		return GroupExpNodeS{node, lpar, rpar}, err
	}

	// If we reach this point, we couldn't match any
	// of the primary expressions, so we need to return
	// an error.
	tok := parser.peek()
	msg := fmt.Sprintf("Expected primary expression got '%v'", tok.Type.String())
	err = parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return nil, err
}

func (parser *MSParser) parseIdentifier() (VariableExpNodeS, error) {

	// Handles "x", "f", ...
	if ok, id := parser.match(IDENTIFIER); ok {
		return VariableExpNodeS{id}, nil
	}

	// Handles "(x)", "(f)", "((x))", ...
	if ok, _ := parser.match(LEFT_PAREN); ok {

		// Recursively parse the identifier when 
		// We encounter a left parenthesis
		ident, err := parser.parseIdentifier()

		// On error, return the error and the identifier
		if err != nil {
			return ident, err
		}

		// Ensure we have a closing parenthesis
		if ok, rp := parser.expect(RIGHT_PAREN); !ok {
			msg := fmt.Sprintf("Expected ')' got '%v'", rp.Type.String())
			err = parser.error(msg, rp.Line, rp.Col)
		}

		return ident, err
	}

	tok := parser.peek()
	msg := fmt.Sprintf("Expected identifier got '%v'", tok.Type.String())
	err := parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return VariableExpNodeS{}, err
}


func flattenExpNode(n *ExpNodeI) []ExpNodeI {

	// By default, flatten returns the node
	// wrapped in a slice
	lexpressions := []ExpNodeI{*n}
	
	// If the node is a tuple, we need to flatten
	// the left side and append the right side.
	switch node := (*n).(type) {
	case BinaryExpNodeS:
		switch node.Op.Type {
		case COMMA:
			lexpressions = append(flattenExpNode(&node.Left), node.Right)
		}
	}

	return lexpressions
}