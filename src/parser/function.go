package parser

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)

func (parser *MSParser) parseFunctionDecl() (ast.FuncDeclNodeS, error) {

	// Parse arguments
	args, err := parser.parseFunctionArgs()

	if (err != nil) {
		return ast.FuncDeclNodeS{Args: args}, err
	}

	ok, grgr := parser.expect(token.GREATER_GREATER)
	if !ok {
		err = parser.error(fmt.Sprintf("Expected '%s' but received '%s'", token.GREATER_GREATER, grgr.Lexeme), grgr.Line, grgr.Col)
		return ast.FuncDeclNodeS{Args: args}, err
	}

	// Parse identifier
	fname, err := parser.parseIdentifier()
	
	if (err != nil){
		return ast.FuncDeclNodeS{Args: args, Fname: fname}, err
	}

	// Expect a '->' into a return type
	ok, arr := parser.expect(token.MINUS_GREAT)
	if !ok {
		err = parser.error(fmt.Sprintf("Expected '%s' but received '%s'", token.MINUS_GREAT, arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Args: args, Fname: fname}, err
	}

	// We expect a return type at this point
	ok, rt := parser.match(token.TypeKeywords...)
	if !ok {
		err := parser.error(fmt.Sprintf("Expected a return type but received '%s'", arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Args: args, Fname: fname, Rt: rt}, err
	}

	// We now expect that the function body is inside "{...}"
	ok, _ = parser.expect(token.LEFT_BRACE)
	if !ok {
		err := parser.error(fmt.Sprintf("Expected a '%s' type but received '%s'", token.LEFT_BRACE, arr.Lexeme), arr.Line, arr.Col)
		return ast.FuncDeclNodeS{Args: args, Fname: fname, Rt: rt}, err
	}

	// Finally, parse block
	block, err := parser.parseBlock()

	// Return whatever
	return ast.FuncDeclNodeS{Args: args, Fname: fname, Rt: rt, Body: block}, err

}

func (parser *MSParser) parseFunctionArgs() ([]ast.FuncArgS, error) {

	// Parses:
	// <empty>
	// int x
	// int x, int y

	// Allocate args list
	args := []ast.FuncArgS{}

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

		// add arg
		args = append(args, ast.FuncArgS{Type: ttok, Iden: vn})

		// Check if we see a ','. If so, we can continue the loop
		// else we have to break the loop. We don't expect a ">>"
		// in this function, it's up to the caller to expect it.
		if ok, _ := parser.match(token.COMMA) ; !ok{
			break
		}
	}

	return args, nil
}