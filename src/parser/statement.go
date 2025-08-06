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
	if ok, _ := parser.match(token.VAR); ok {
		return parser.parseVarDeclaration()
	}
	if ok, _ := parser.match(token.TYPE) ; ok {
		return parser.parseTypeDeclaration()
	}
	// CONTINUE
	if ok, tk := parser.match(token.CONTINUE); ok {
		return parser.parseContinue(tk)
	}
	// BREAK
	if ok, tk := parser.match(token.BREAK); ok {
		return parser.parseBreak(tk)
	}
	// RETURN
	if ok, tk := parser.match(token.RETURN) ; ok {
		return parser.parseReturn(tk)
	}

	// Nothing matched, so we assume it must be an expression.
	return parser.parseExpressionStatement()
}

