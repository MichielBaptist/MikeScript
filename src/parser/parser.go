package parser

import (
	"fmt"
	ast "mikescript/src/ast"
	token "mikescript/src/token"
	"slices"
)

////////////////////////////////////////////////////////////
// 						Parser
////////////////////////////////////////////////////////////

type MSParser struct {
	src string				// source code
	tokens []token.Token	// token list from tokenizer
	pos int     			// current position in tokens
	pnc bool    			// panic flag
	Errors []ParserError	// parser errors
	context []ParserConext	// nothing, loop, function...
}

////////////////////////////////////////////////////////////
// 						helpers
////////////////////////////////////////////////////////////

func (parser *MSParser) SetSrc(src string) {
	parser.src = src
}

func (parser *MSParser) SetTokens(tokens []token.Token) {
	parser.tokens = tokens
}


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
	if slices.ContainsFunc(t, parser.checkType) {
			return true, parser.peek()
		}
	return false, parser.peek()
}

func (parser *MSParser) match(t ...token.TokenType) (bool, token.Token) {
	
	// Only advance if we match a type
	if slices.ContainsFunc(t, parser.checkType) {
			return true, parser.advance()
		}

	// No match
	return false, parser.peek()
}

func (parser *MSParser) expect(t token.TokenType) (bool, token.Token) {
	return parser.match(t)
}

func (p *MSParser) enterContext(ctx ParserConext) {
	p.context = append(p.context, ctx)
}

func (p *MSParser) leaveContext() ParserConext{
	// throw error if something goes wrong
	ctx := p.context[len(p.context)-1]
	p.context = p.context[:len(p.context)-1]
	return ctx
}

func (p *MSParser) inContext(ctx ParserConext) bool {
	// LOOP context: must be the last context seen
	// FUNCTION context: must be in the context stack (contains)
	switch ctx {
	case LOOP:		return p.context[len(p.context)-1] == LOOP
	case FUNCTION:	return slices.Contains(p.context, ctx)
	default: 		_ = []int{}[0]
	}
	return slices.Contains(p.context, ctx)
}

////////////////////////////////////////////////////////////
// 						Error handling
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

func (parser *MSParser) Parse(tokens []token.Token) (*ast.Program, error) {
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
