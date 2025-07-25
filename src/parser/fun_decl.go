package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (parser *MSParser) parseFunctionDecl() (ast.FuncDeclNodeS, error) {

	// Parse arguments 'x, y, z, ...'
	args, err := parser.parseFunctionArgs()

	if (err != nil) {
		return ast.FuncDeclNodeS{Params: args}, err
	}

	// Parse '>>'
	ok, tok := parser.expect(token.GREATER_GREATER)
	if !ok {
		err = parser.error(fmt.Sprintf("Expected '%s' but received '%s'", token.GREATER_GREATER, tok.Lexeme), tok.Line, tok.Col)
		return ast.FuncDeclNodeS{Params: args}, err
	}

	// Parse identifier "f", "g", ...
	fname, err := parser.parseIdentifier()
	if (err != nil){
		return ast.FuncDeclNodeS{Params: args, Fname: fname}, err
	}

	// Parse '->'
	ok, arr := parser.expect(token.MINUS_GREAT)
	if !ok {
		err = parser.error(fmt.Sprintf("Expected '%s' but received '%s'", token.MINUS_GREAT, arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Params: args, Fname: fname}, err
	}

	// We expect a return type at this point
	ok, rt := parser.match(token.TypeKeywords...)

	// Convert token to mstype
	mrt, err := mstype.TokenToType(&rt)
	if err != nil {
		return ast.FuncDeclNodeS{}, err
	}

	if !ok {
		err := parser.error(fmt.Sprintf("Expected a return type but received '%s'", arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: mrt}, err
	}

	// Check which token comes next:
	// 1. '{' means we parse a block.
	// 2. ';' means we have a function declaration without body.
	// 3. Anything else is a syntax error.
	switch _, tok = parser.match(token.SEMICOLON, token.LEFT_BRACE); tok.Type {
	case token.LEFT_BRACE:

		// got '{' so expect we can parse a block
		block, err := parser.parseBlock()

		if err != nil {
			// Something went wrong while parsing the block
			return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: mrt}, err
		}
		
		return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: mrt, Body: &block}, err

	case token.SEMICOLON:
		return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: mrt, Body: nil}, nil
	default:
		err = parser.error(fmt.Sprintf("Expected a '%s' type but received '%s'", token.LEFT_BRACE, arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Params: args, Fname: fname, Rt: mrt, Body: nil}, err
	}
}

func (parser *MSParser) parseFunctionArgs() ([]ast.FuncParamS, error) {

	// Parses:
	// <empty>
	// int x
	// int x, int y

	// Allocate args list
	args := []ast.FuncParamS{}

	// Match a type, if no type found we don't enter the loop
	ok, ttok := parser.match(token.TypeKeywords...)

	// Keep going if we see 
	for ; ok ; ok, ttok = parser.match(token.TypeKeywords...){

		// Parse identifier
		vn, err := parser.parseIdentifier()

		// Check if parsing went ok, else we return
		if err != nil {
			return args, err
		}

		vtype, err := mstype.TokenToType(&ttok)

		if err != nil {
			return args, err
		}

		// add arg
		args = append(args, ast.FuncParamS{Type: vtype, Iden: vn})

		// Check if we see a ','. If so, we can continue the loop
		// else we have to break the loop. We don't expect a ">>"
		// in this function, it's up to the caller to expect it.
		if ok, _ := parser.match(token.COMMA) ; !ok{
			break
		}
	}

	return args, nil
}