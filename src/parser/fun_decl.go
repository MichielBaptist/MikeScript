package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (parser *MSParser) parseFunctionDecl() (ast.FuncDeclNodeS, error) {
	// parses: arguments '>>' IDENTIFIER {'->' type}? '{' block
	// 1. arguments
	// 2. '>>'
	// 3. IDENTIFER
	// 4. {'->' type}?
	// 5. '{' block


	// 1. Parse arguments
	args, err := parser.parseFunctionArgs()
	if (err != nil) {
		return ast.FuncDeclNodeS{Params: args}, err
	}

	// 2. Parse '>>'
	ok, tok := parser.expect(token.GREATER_GREATER)
	if !ok {
		err = parser.error(fmt.Sprintf("Expected '%s' but received '%s'", token.GREATER_GREATER, tok.Lexeme), tok.Line, tok.Col)
		return ast.FuncDeclNodeS{Params: args}, err
	}

	// 3. Parse identifier "f", "g", ...
	fname, err := parser.parseIdentifier()
	if (err != nil){
		return ast.FuncDeclNodeS{Params: args, Fname: fname}, err
	}

	// 4. Parse {'->' type}?
	var returnType mstype.MSType
	if ok, _ := parser.match(token.MINUS_GREAT) ; ok {
		returnType, err = parser.parseType()
	} else {
		returnType = mstype.MS_NOTHING
		err = nil
	}
	if err != nil {
		return ast.FuncDeclNodeS{}, err
	}

	// 5. '{' block
	if ok, tok := parser.match(token.LEFT_BRACE) ; !ok {
		return ast.FuncDeclNodeS{}, parser.unexpectedToken(tok, token.LEFT_BRACE)
	}
	block, err := parser.parseBlock()
	if err != nil {
		return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: returnType}, err
	}

	return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: returnType, Body: &block}, err
}

func (parser *MSParser) parseFunctionArgs() ([]ast.FuncParamS, error) {
	// Parses:
	// <empty>
	// int x
	// int x, int y

	// Parse '('
	if ok, tok := parser.match(token.LEFT_PAREN) ; !ok {
		return []ast.FuncParamS{}, parser.unexpectedToken(tok, token.LEFT_PAREN)
	}

	// Allocate args list
	args := []ast.FuncParamS{}
	for parser.peek().Type != token.RIGHT_PAREN {

		// type
		paramType, err := parser.parseType()

		if err != nil {
			return args, err
		}

		// IDENTIFIER
		ident, err := parser.parseIdentifier()

		if err != nil {
			return args, err
		}

		// add to args
		args = append(args, ast.FuncParamS{Type: paramType, Iden: ident})

		// Check if we see a ','. If so, we can continue the loop
		// else we have to break the loop. We don't expect a ">>"
		// in this function, it's up to the caller to expect it.
		if ok, _ := parser.match(token.COMMA) ; !ok{
			break
		}
	}

	// We expect a closing paren
	ok, tk := parser.match(token.RIGHT_PAREN)

	if !ok {
		return args, parser.unexpectedToken(tk, token.RIGHT_BRACE)
	}

	return args, nil
}