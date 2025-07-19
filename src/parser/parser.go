package parser

import (
	"fmt"
	AST "mikescript/src/ast"
	token "mikescript/src/token"
)

////////////////////////////////////////////////////////////
// 						Parser
////////////////////////////////////////////////////////////

type Parser interface {
	parse(tokens []token.Token) (AST.ExpNodeI, error)
}

type MSParser struct {
	src    string			// source code
	tokens []token.Token	// token list from tokenizer
	pos  int     			// current position in tokens
	pnc  bool    			// panic flag
	Errors []ParserError	// parser errors
}

func (parser *MSParser) SetSrc(src string) {
	parser.src = src
}

func (parser *MSParser) SetTokens(tokens []token.Token) {
	parser.tokens = tokens
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

func (parser *MSParser) advance() token.Token {
	// Peeks current token and advances next position.
	tok := parser.peek()
	parser.pos = parser.pos + 1
	return tok
}

func (parser *MSParser) peek() token.Token {
	if parser.atend() {
		return token.Token{Type: token.UNKNOWN, Lexeme: "UNKNOWN", Line: 0, Col: 0}
	}
	return parser.tokens[parser.pos]
}

func (parser *MSParser) atend() bool {
	// When past the token stack, we are at end.
	return parser.pos >= len(parser.tokens)
}

func (parser *MSParser) checkType(t token.TokenType) bool {
	return !parser.atend() && parser.peek().Type == t
}

func (parser *MSParser) match(t ...token.TokenType) (bool, token.Token) {
	// check if we matched the token types
	// If we did, we advance the position
	for _, tt := range t {
		if parser.checkType(tt) {
			return true, parser.advance()
		}
	}

	// Did not match any of the tokens
	// return false and an empty token
	return false, parser.peek()
}

func (parser *MSParser) expect(t token.TokenType) (bool, token.Token) {
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
	for tok.Type != token.SEMICOLON && !parser.atend() && tok.Type != token.EOF {
		tok = parser.advance()
	}

	fmt.Println("Synchronized to: ", tok)

	if tok.Type != token.EOF {
		parser.advance()
	}
	
}

////////////////////////////////////////////////////////////
// 						Parse
////////////////////////////////////////////////////////////

func (parser *MSParser) Parse(tokens []token.Token) (AST.Program, error) {
	// Parses: START -> Program EOF

	// parse block
	ast, err := parser.parseProgram()

	// If parsing failed, return the error
	if err != nil {
		return ast, err
	}

	// Expect token.EOF, otherwise return an error
	if ok, tok := parser.expect(token.EOF); !ok {
		msg := fmt.Sprintf("Expected 'token.EOF' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return ast, err
}

func (parser *MSParser) parseProgram() (AST.Program, error) {
	// parses program -> statement *

	// as long as we haven't reached the end of the tokens
	// we keep parsing statements.
	statements := []AST.StmtNodeI{}
	var err error
	var stmt AST.StmtNodeI

	for !parser.atend() && parser.peek().Type != token.EOF {

		// parse next statement and check for error
		stmt, err = parser.parseStatement()

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
	return AST.Program{Statements: statements}, err
}

// func (parser *MSParser) parseDeclaration() (AST.StmtNodeI, error) {

// 	// Found a declaration, parse it
// 	if ok, tok := parser.match(token.INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOLEAN_TYPE); ok {
// 		return parser.parseVarDeclaration(tok)
// 	}

// 	// fall through to statement parsing
// 	return parser.parseStatement()
// }

func (parser *MSParser) parseVarDeclaration(declarationType token.Token) (AST.DeclarationNodeS, error) {
	// arg: type of declartion;
	
	// We expect an identifier next so we parse it
	ident, err := parser.parseIdentifier()

	// check for errors
	if err != nil {
		// If we encounter an error parsing the identifier
		// We return the error and empty declaration node
		return AST.DeclarationNodeS{}, err
	}

	// expect a semicolon
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return AST.DeclarationNodeS{Identifier: ident, Vartype: declarationType}, err

}

func (parser *MSParser) parseStatement() (AST.StmtNodeI, error){
	// statement ->
	// [0]: block
	// [1]: if
	// [2]: variable_declaration
	// [3]: while
	// [-]: expression

	// BLOCK
	if ok, _ := parser.match(token.LEFT_BRACE); ok {
		return parser.parseBlock()
	}
	// IF
	if ok, _ := parser.match(token.IF); ok {
		return parser.parseIf()
	}
	// WHILE
	if ok, _ := parser.match(token.WHILE); ok {
		return parser.parseWhile()
	}
	// VARIABLE DECLARATION
	if ok, tk := parser.match(token.INT_TYPE, token.FLOAT_TYPE, token.STRING_TYPE, token.BOOLEAN_TYPE); ok {
		return parser.parseVarDeclaration(tk)
	}

	// CONTINUE
	if ok, tk := parser.match(token.CONTINUE); ok {
		err := parser.error("Continue statement not allowed outside of loops", tk.Line, tk.Col)
		return AST.ContinueNodeS{Tk: tk}, err
	}
	// BREAK
	if ok, tk := parser.match(token.BREAK); ok {
		err := parser.error("Break statement not allowed outside of loops", tk.Line, tk.Col)
		return AST.BreakNodeS{Tk: tk}, err
	}

	return parser.parseExpressionStatement()
}

func (parser *MSParser) parseBreak(tk token.Token) (AST.BreakNodeS, error) {
	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err := parser.error(msg, tok.Line, tok.Col)
		return AST.BreakNodeS{}, err
	}
	return AST.BreakNodeS{Tk: tk}, nil
}

func (parser *MSParser) parseContinue(tk token.Token) (AST.ContinueNodeS, error) {
	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err := parser.error(msg, tok.Line, tok.Col)
		return AST.ContinueNodeS{}, err
	}
	return AST.ContinueNodeS{Tk: tk}, nil
}


func (parser *MSParser) parseIf() (AST.IfNodeS, error) {

	// Parse conditional expression
	cond, err := parser.parseExpression()
	if err != nil {
		return AST.IfNodeS{}, err
	}

	// Expect a block statement so we check for opening brace.
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return AST.IfNodeS{}, err
	}
	stmt, err := parser.parseBlock()
	if err != nil {
		return AST.IfNodeS{}, err
	}

	// Check for if statement else branch
	// and parse it if it exists.
	if ok, _ := parser.match(token.ELSE); ok {

		// Expect a block statement
		if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
			msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
			err = parser.error(msg, tok.Line, tok.Col)
			return AST.IfNodeS{}, err
		}
		elsestmt, err := parser.parseBlock()
		if err != nil {
			return AST.IfNodeS{}, err
		}

		return AST.IfNodeS{Condition: cond, ThenStmt: stmt, ElseStmt: elsestmt}, err
	}

	return AST.IfNodeS{Condition: cond, ThenStmt: stmt, ElseStmt: nil}, err

}

func (parser *MSParser) parseWhile() (AST.WhileNodeS, error) {

	// 1. parse conditional expression
	cond, err := parser.parseExpression()

	// check for errors
	if err != nil {
		return AST.WhileNodeS{}, err
	}

	// Expect a block statement
	if ok, tok := parser.expect(token.LEFT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '{' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
		return AST.WhileNodeS{}, err
	}
	block, err := parser.parseBlock()

	// Pack into a while node
	return AST.WhileNodeS{Condition: cond, Body: block}, err

}

func (parser *MSParser) parseBlock() (AST.BlockNodeS, error) {

	// Parse the block
	stmts := []AST.StmtNodeI{}
	var err error
	var stmt AST.StmtNodeI

	// As long as we haven't reached the end of the block
	// we keep parsing statements.
	for !parser.atend() && parser.peek().Type != token.RIGHT_BRACE {

		// continue
		if ok, tk := parser.match(token.CONTINUE); ok {
			stmt, err = parser.parseContinue(tk)
		} else if ok, tk := parser.match(token.BREAK); ok {
			stmt, err = parser.parseBreak(tk)
		} else {
			stmt, err = parser.parseStatement()
		}

		// Only add to the statements if there was no error
		// Else we just synchronize and continue
		if err == nil {
			stmts = append(stmts, stmt)
		} else {
			return AST.BlockNodeS{}, err
		}
	}

	// Expect a closing brace
	if ok, tok := parser.expect(token.RIGHT_BRACE); !ok {
		msg := fmt.Sprintf("Expected '}' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return AST.BlockNodeS{Statements: stmts}, err
}


func (parser *MSParser) parseExpressionStatement() (AST.ExStmtNodeS, error) {

	// Parse the expression.
	xpr, err := parser.parseExpression()

	// Check for errors
	if err != nil {
		return AST.ExStmtNodeS{Ex: xpr}, err
	}

	// Expect a semicolon. If not found, return an error
	// set the error message to the expected token.
	if ok, tok := parser.expect(token.SEMICOLON); !ok {
		msg := fmt.Sprintf("Expected ';' got '%v'", tok.Type.String())
		err = parser.error(msg, tok.Line, tok.Col)
	}

	return AST.ExStmtNodeS{Ex: xpr}, err
}

func (parser *MSParser) parseExpression() (AST.ExpNodeI, error) {
	return parser.parseLor()
}

func (parser *MSParser) parseLor() (AST.ExpNodeI, error) {
	land, ok := parser.parseLand()

	// Error on parsing logical and
	if ok != nil {
		return land, ok
	}

	for {
		if ok, op := parser.match(token.BAR_BAR); ok {
			right, err := parser.parseLand()
			land = AST.LogicalExpNodeS{Left: land, Op: op, Right: right}

			// check for errors
			if err != nil {
				return land, err
			}
		} else {
			break
		}
	}

	return land, nil
}

func (parser *MSParser) parseLand() (AST.ExpNodeI, error) {
	comp, ok := parser.parseFuncop()
	if ok != nil {
		return comp, ok	
	}
	for {
		if ok, op := parser.match(token.AMP_AMP); ok {
			right, err := parser.parseFuncop()
			comp = AST.LogicalExpNodeS{Left: comp, Op: op, Right: right}
			if err != nil {
				return comp, err
			}
		} else {
			break
		}
	}
	return comp, nil
}

func (parser *MSParser) parseFuncop() (AST.ExpNodeI, error) {

	var left AST.ExpNodeI
	var err error

	// Check if we have an argument-less function
	// application. If we do, we return the function
	// application node.
	if ok, id := parser.match(token.GREATER_GREATER); ok {
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
		if ok, op := parser.match(token.GREATER_GREATER, token.MINUS_GREAT); ok {

			var right AST.ExpNodeI
			var err error

			// Parse the right side of the '>>'. If this is
			// and identifier, we have either a function application
			// or a variable assignment. The parser cannot know
			// which one it is, we only know this at runtime time.
			switch op.Type {
			case token.GREATER_GREATER:

				// Parse the right side of the function application
				// This should resolve into an identifier.
				right, err = parser.parseTuple()

				// Flatten the left side of the function application
				// into a list of expressions
				lexpressions := flattenExpNode(&left)

				// Create a function application node
				left = AST.FuncAppNodeS{Args: lexpressions, Fun: right}

			case token.MINUS_GREAT:

				// Parse the right side of the assignment
				// This should resolve into:
				// 1. A single variable
				// 2. A tuple of variables (not implemented)
				variable, verr := parser.parseTuple()
				err = verr

				// TODO: add support for tuple assignments
				switch v := variable.(type) {
				case AST.VariableExpNodeS:
					left = AST.AssignmentNodeS{Identifier: v, Exp: left}
				default:
					err = parser.error(fmt.Sprintf("Expected a variable, got '%v'", v), op.Line, op.Col)
					return left, err
				}
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

func (parser *MSParser) parseTuple() (AST.ExpNodeI, error) {
	
	// parse the first expression
	node, err := parser.parseEquality()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.COMMA); ok {
			right, err := parser.parseEquality()
			node = AST.BinaryExpNodeS{Left: node, Op: op, Right: right}

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

func (parser *MSParser) parseEquality() (AST.ExpNodeI, error) {
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
			node = AST.BinaryExpNodeS{Left: node, Op: eqop, Right: right}

			// If the operator was neq, we wrap the binary node
			if is_neq {
				neg := token.Token{Type: token.EXCLAMATION, Lexeme: "!", Col: op.Col, Line: op.Line}
				node = AST.UnaryExpNodeS{Op: neg, Node: node}
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

func (parser *MSParser) parseComp() (AST.ExpNodeI, error) {

	node, err := parser.parseTerm()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.LESS, token.GREATER, token.LESS_EQ, token.GREATER_EQ); ok {
			right, err := parser.parseTerm()

			node = AST.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}
		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseTerm() (AST.ExpNodeI, error) {

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
				right = AST.UnaryExpNodeS{Op: op, Node: right}        // x (-y)
				op = token.Token{Type: token.PLUS, Lexeme: "+", Col: op.Col, Line: op.Line} // x + (-y)
			}
			node = AST.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseFactor() (AST.ExpNodeI, error) {

	node, err := parser.parseUnary()

	if err != nil {
		return node, err
	}

	for {
		if ok, op := parser.match(token.MULT, token.SLASH, token.PERCENT); ok {
			right, err := parser.parseUnary()
			node = AST.BinaryExpNodeS{Left: node, Op: op, Right: right}

			if err != nil {
				return node, err
			}

		} else {
			break
		}
	}

	return node, err
}

func (parser *MSParser) parseUnary() (AST.ExpNodeI, error) {

	if ok, op := parser.match(token.MINUS, token.EXCLAMATION); ok {
		right, err := parser.parseUnary()

		if err != nil {
			return right, err
		}

		return AST.UnaryExpNodeS{Op: op, Node: right}, nil
	}

	return parser.parsePrimary()
}

func (parser *MSParser) parsePrimary() (AST.ExpNodeI, error) {

	var err error = nil

	// matches a primary expression
	if ok, tok := parser.match(token.NUMBER_INT, token.NUMBER_FLOAT, token.STRING, token.TRUE, token.FALSE); ok {
		return AST.LiteralExpNodeS{Tk: tok}, err
	}

	// matches an identifier
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return AST.VariableExpNodeS{Name: id}, err
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
		return AST.GroupExpNodeS{Node: node, TokenLeft: lpar, TokenRight: rpar}, err
	}

	// Matches variable declaration
	// if ok, typ := parser.match(token.INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOLEAN_TYPE); ok {
		
	// 	ident, err := parser.parseIdentifier()

	// 	// Check for errors
	// 	if err != nil {
	// 		return AST.DeclarationNodeS{}, err
	// 	}

	// 	return AST.DeclarationNodeS{ident, typ}, err
	// }

	// If we reach this point, we couldn't match any
	// of the primary expressions, so we need to return
	// an error.
	tok := parser.peek()
	msg := fmt.Sprintf("Expected primary expression got '%v'", tok.Type.String())
	err = parser.error(msg, tok.Line, tok.Col)
	parser.panic()

	return nil, err
}

func (parser *MSParser) parseIdentifier() (AST.VariableExpNodeS, error) {

	// Handles "x", "f", ...
	if ok, id := parser.match(token.IDENTIFIER); ok {
		return AST.VariableExpNodeS{Name: id}, nil
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

	return AST.VariableExpNodeS{}, err
}


func flattenExpNode(n *AST.ExpNodeI) []AST.ExpNodeI {

	// By default, flatten returns the node
	// wrapped in a slice
	lexpressions := []AST.ExpNodeI{*n}
	
	// If the node is a tuple, we need to flatten
	// the left side and append the right side.
	switch node := (*n).(type) {
	case AST.BinaryExpNodeS:
		switch node.Op.Type {
		case token.COMMA:
			lexpressions = append(flattenExpNode(&node.Left), node.Right)
		}
	}

	return lexpressions
}