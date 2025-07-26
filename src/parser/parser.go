package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

////////////////////////////////////////////////////////////
// 						Parser
////////////////////////////////////////////////////////////

type Parser interface {
	parse(tokens []token.Token) (ast.ExpNodeI, error)
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

func (parser *MSParser) lookahead(t ...token.TokenType) (bool, token.Token) {
	for _, tt := range t {
		if parser.checkType(tt) {
			return true, parser.peek()
		}
	}
	return false, parser.peek()
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

func (parser *MSParser) synchronize() {

	tok := parser.peek()

	// Synchronize the parser by advancing the position
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

func (parser *MSParser) Parse(tokens []token.Token) (ast.Program, error) {
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

////////////////////////////////////////////////////////////
// Common parser errors
////////////////////////////////////////////////////////////

func (parser *MSParser) error(msg string, line, col int) error {
	err := ParserError{msg, line, col}
	parser.Errors = append(parser.Errors, err)
	parser.panic()
	return err
}

func (parser *MSParser) unexpectedToken(got token.Token, expected ...token.TokenType) error {
	msg := fmt.Sprintf("expected '%s' got '%s'", expected, got)
	return parser.error(msg, got.Line, got.Col)
}