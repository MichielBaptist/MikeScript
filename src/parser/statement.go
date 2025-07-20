package parser

import (
	ast "mikescript/src/ast"
	token "mikescript/src/token"
)

func (parser *MSParser) parseStatement() (ast.StmtNodeI, error){
	// statement ->
	// [0]: block
	// [1]: if
	// [2]: variable_declaration
	// [3]: while
	// [-]: expression

	// FUNCDECL
	if ok, _ := parser.match(token.FUNCTION); ok {
		return parser.parseFunctionDecl()
	}
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
	if ok, tk := parser.match(token.TypeKeywords...); ok {
		return parser.parseVarDeclaration(tk)
	}
	// CONTINUE
	if ok, tk := parser.match(token.CONTINUE); ok {
		err := parser.error("Continue statement not allowed outside of loops", tk.Line, tk.Col)
		return ast.ContinueNodeS{Tk: tk}, err
	}
	// BREAK
	if ok, tk := parser.match(token.BREAK); ok {
		err := parser.error("Break statement not allowed outside of loops", tk.Line, tk.Col)
		return ast.BreakNodeS{Tk: tk}, err
	}

	// Nothing matched, so we assume it must be an expression.
	return parser.parseExpressionStatement()
}

