package token

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
	COLON							// : UNUSED
	PERCENT							// % 
	EXCLAMATION						// ! 
	LESS							// < 
	GREATER							// >
	BAR								// | (guard in xif) UNUSED
	EQ								// = (function call)

	// Double character tokens
	EXCLAMATION_EQ					// != 
	EQ_EQ							// == 
	INT_DIV							// // 
	DOT_DOT							// .. UNUSED
	LESS_EQ							// <=
	GREATER_EQ						// >=
	GREATER_GREATER					// >> (function bind)
	DOT_GREATER_GREATER				// .>> (mapped function bind)
	MULT_GREATER_GREATER			// *>> (unpack && bind)
	GREATER_GREATER_EQ				// >>= (function bind && call)
	DOT_GREATER_GREATER_EQ			// .>>= (mapped function bind && call)
	MULT_GREATER_GREATER_EQ			// *>>= (unpack && bind && call)
	LESS_LESS						// <<
	MINUS_GREAT						// -> (assignment)
	DOT_MINUS_GREAT					// .-> (for loop separator)
	EQ_GREATER						// => (decl & assignment)
	LESS_MINUS						// <- UNUSED
	AMP_AMP							// &&
	BAR_BAR							// ||
	DOT_EQ 							// .=

	// Literals
	IDENTIFIER						// Identifier (x, y, z, f, etc)
	STRING							// String literal
	NUMBER_INT						// Number literal (no dot)
	NUMBER_FLOAT					// Number literal (with dot)

	// Keywords
	FALSE 							// false
	TRUE 							// true
	IF								// if
	ELSE 							// else
	XIF 							// xif UNUSED
	OTHERWISE 						// otherwise UNUSED
	FOR 							// for
	WHILE 							// while
	FUNCTION 						// function
	RETURN 							// return
	PRINT 							// print
	CONTINUE 						// continue
	BREAK 							// break
	VAR								// var
	TYPE							// type

	// Types
	INT_TYPE 						// int (64)
	FLOAT_TYPE 						// float (64)
	STRING_TYPE 					// string
	BOOLEAN_TYPE 					// boolean
	NOTHING_TYPE					// nothing
	STRUCT 							// struct UNUSED

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
	COLON: ":",
	PERCENT: "%",
	EXCLAMATION: "!",
	LESS: "<",
	GREATER: ">",
	BAR: "|",
	EQ: "=",
	DOT_EQ: ".=",
	EXCLAMATION_EQ: "!=",
	EQ_EQ: "==",
	INT_DIV: "//",
	DOT_DOT: "..",
	LESS_EQ: "<=",
	GREATER_EQ: ">=",
	GREATER_GREATER: ">>",
	MULT_GREATER_GREATER: "*>>",
	GREATER_GREATER_EQ: ">>=",
	DOT_GREATER_GREATER: ".>>",
	DOT_GREATER_GREATER_EQ: ".>>=",
	MULT_GREATER_GREATER_EQ: "*>>=",
	MINUS_GREAT: "->",
	DOT_MINUS_GREAT: ".->",
	EQ_GREATER: "=>",
	LESS_LESS: "<<",
	AMP_AMP: "&&",
	BAR_BAR: "||",
	IDENTIFIER: "ID",
	STRING: "l_str",
	NUMBER_INT: "l_int",
	NUMBER_FLOAT: "l_float",
	FALSE: "false",
	TRUE: "true",
	XIF: "xif",
	IF: "if",
	ELSE: "else",
	OTHERWISE: "otherwise",
	FOR: "for",
	WHILE: "while",
	FUNCTION: "function",
	RETURN: "return",
	PRINT: "print",
	INT_TYPE: "t_int",
	FLOAT_TYPE: "t_float",
	STRING_TYPE: "t_str",
	BOOLEAN_TYPE: "t_bool",
	EOF: "EOF",
	UNKNOWN: "UNKNOWN",
	CONTINUE: "continue",
	BREAK: "break",
	VAR: "var",
	TYPE: "type",
}

// Map of keywords
var Keywords map[string]TokenType = map[string]TokenType{
	"false": FALSE,
	"true": TRUE,
	"xif": XIF,
	"if": IF,
	"else": ELSE,
	"otherwise": OTHERWISE,
	"for": FOR,
	"while": WHILE,
	"function": FUNCTION,
	"return": RETURN,
	"int": INT_TYPE,
	"float": FLOAT_TYPE,
	"string": STRING_TYPE,
	"bool": BOOLEAN_TYPE,
	"continue": CONTINUE,
	"break": BREAK,
	"var": VAR,
	"type": TYPE,
	"struct": STRUCT,
	"nothing": NOTHING_TYPE,
}

// implement stringer
func (t TokenType) String() string {
	return stmp[t]
}

// List of tokens which define a builtin type
var SimpleTypeKeywords []TokenType = []TokenType{
	INT_TYPE,
	FLOAT_TYPE,
	STRING_TYPE,
	BOOLEAN_TYPE,
	NOTHING_TYPE,
}