package main

type TokenType uint8

type Token struct {
	Type 	TokenType	// Type of token
	Lexeme 	string		// Lexeme of token (string representation)
	Line 	int			// Line number of token
	Col 	int			// Column number of token
}

const(
	// Single character tokens
	INVALID_TOKEN TokenType = iota	// Invalid token
	LEFT_PAREN 						// ( 
	RIGHT_PAREN						// ) 
	LEFT_BRACE						// { 
	RIGHT_BRACE						// } 
	LEFT_SQUARE						// [ 
	RIGHT_SQUARE					// ] 
	COMMA							// , 
	DOT								// . 
	PLUS							// + 
	MINUS							// - 
	MULT							// * 
	SLASH							// / 
	SEMICOLON						// ; 
	PERCENT							// % 
	EXCLAMATION						// ! 
	LESS							// < 
	GREATER							// >
	BAR								// | (guard in xif)
	EQ								// = (unused currently)

	// Double character tokens
	EXCLAMATION_EQ					// != 
	EQ_EQ							// == 
	INT_DIV							// // 
	DOT_DOT							// ..
	LESS_EQ							// <=
	GREATER_EQ						// >=
	GREATER_GREATER					// >> (function calls)
	LESS_LESS						// <<
	MINUS_GREAT						// -> (assignment)
	LESS_MINUS						// <-
	AMP_AMP							// &&
	BAR_BAR							// ||

	// Literals
	IDENTIFIER						// Identifier (x, y, z, f, etc)
	STRING							// String literal
	NUMBER_INT						// Number literal (no dot)
	NUMBER_FLOAT					// Number literal (with dot)

	// Keywords
	AND 							// and
	OR 								// or
	FALSE 							// false
	TRUE 							// true
	XIF 							// xif
	OTHERWISE 						// otherwise
	FOR 							// for
	WHILE 							// while
	FUNCTION 						// function
	RETURN 							// return
	PRINT 							// print

	// Types
	INT_TYPE 						// int (64)
	FLOAT_TYPE 						// float (64)
	STRING_TYPE 					// string
	BOOLEAN_TYPE 					// boolean
	VOID_TYPE 						// nothing
	STRUCT 							// struct

	// End of file
	EOF								// End of file

	// Unknown token
	UNKNOWN							// Unknown token
)

var stmp map[TokenType]string = map[TokenType]string{
	LEFT_PAREN: "(",
	RIGHT_PAREN: ")",
	LEFT_BRACE: "{",
	RIGHT_BRACE: "}",
	LEFT_SQUARE: "[",
	RIGHT_SQUARE: "]",
	COMMA: ",",
	DOT: ".",
	PLUS: "+",
	MINUS: "-",
	MULT: "*",
	SLASH: "/",
	SEMICOLON: ";",
	PERCENT: "%",
	EXCLAMATION: "!",
	LESS: "<",
	GREATER: ">",
	BAR: "|",
	EXCLAMATION_EQ: "!=",
	EQ_EQ: "==",
	INT_DIV: "//",
	DOT_DOT: "..",
	LESS_EQ: "<=",
	GREATER_EQ: ">=",
	GREATER_GREATER: ">>",
	LESS_LESS: "<<",
	AMP_AMP: "&&",
	BAR_BAR: "||",
	IDENTIFIER: "IDENTIFIER",
	STRING: "STRING",
	NUMBER_INT: "int",
	NUMBER_FLOAT: "float",
	AND: "and",
	OR: "or",
	FALSE: "false",
	TRUE: "true",
	XIF: "xif",
	OTHERWISE: "otherwise",
	FOR: "for",
	WHILE: "while",
	FUNCTION: "function",
	RETURN: "return",
	PRINT: "print",
	INT_TYPE: "int",
	FLOAT_TYPE: "float",
	STRING_TYPE: "string",
	BOOLEAN_TYPE: "bool",
	VOID_TYPE: "void",
	EOF: "EOF",
	UNKNOWN: "UNKNOWN",
}

// Map of keywords
var keywords map[string]TokenType = map[string]TokenType{
	"and": AND,
	"or": OR,
	"false": FALSE,
	"true": TRUE,
	"xif": XIF,
	"otherwise": OTHERWISE,
	"for": FOR,
	"while": WHILE,
	"function": FUNCTION,
	"return": RETURN,
	"int": INT_TYPE,
	"float": FLOAT_TYPE,
	"string": STRING_TYPE,
	"bool": BOOLEAN_TYPE,
	"nothing": VOID_TYPE,
}


// implement stringer
func (t TokenType) String() string {
	return stmp[t]
}

func (t Token) String() string {
	return "(" + t.Type.String() + ", " + t.Lexeme + ")"
}
