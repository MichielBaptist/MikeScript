package scanner

import (
	"fmt"
	token "mikescript/src/token"
	utils "mikescript/src/utils"
)

// TODO: fix this using regex instead?

const(
	SPACE byte = ' '
	TAB byte = '\t'
	NEWLINE byte = '\n'
	QUOTE byte = '"'
)

type Scanner interface {
	Scan(input string) []token.Token
}

type MSScanner struct {

	// src information
	src string 		// Source code to scan
	n int 			// Length of source code

	// scanner state
	tokens []token.Token 	// token.Tokens found in source code
	l int 			// Start of current token
	r int 			// Current position in source code
	line int 		// Current line number
	col int 		// Current column number

	// error information
	Errors []ScannerError
}

type ScannerError struct {
	msg 	string
	line 	int
	col 	int
}

func (err ScannerError) Compare(other ScannerError) bool {
	return err.msg == other.msg && err.line == other.line && err.col == other.col
}

func (err ScannerError) String() string {
	return fmt.Sprintf("[Scanner Error]: %v at line %v col %v", err.msg, err.line, err.col)
}

////////////////////////////////////////////////////////////////
// 							helpers
////////////////////////////////////////////////////////////////

func (scanner *MSScanner) advance() byte {
	c := scanner.atr()
	scanner.r++
	scanner.col++
	return c
}

func (scanner *MSScanner) newline() {
	scanner.line++
	scanner.col = 1
}

func (scanner *MSScanner) atl() byte {
	if scanner.l >= scanner.n { return 0 }
	return scanner.src[scanner.l]
}

func (scanner *MSScanner) atr() byte {
	if scanner.r >= scanner.n { return 0 }
	return scanner.src[scanner.r]
}

func (scanner *MSScanner) atrIsDigit() bool {
	return utils.IsDigit(scanner.atr())
}

////////////////////////////////////////////////////////////////
// 							Scan
////////////////////////////////////////////////////////////////

func (scanner *MSScanner) Scan(input string) []token.Token {
	scanner.reset()
	scanner.setSrc(input)
	scanner.scanTokens()
	return scanner.tokens
}

func (scanner *MSScanner) scanTokens() {

	for !scanner.atEnd() {
		if token, ok := scanner.nextToken(); ok {
			scanner.addToken(token)
		}
	}

	// finish with EOF token
	scanner.addToken(token.Token{Type: token.EOF, Lexeme: "", Line: scanner.line, Col: scanner.col})
}

func (scanner *MSScanner) nextToken() (token.Token, bool) {

	// get current char and advance r to next character
	c := scanner.advance()

	var tok token.Token      	// nil token
	var ok bool = true	// valid token

	switch {

	// handle witespace and newlines
	case c == SPACE || c == TAB:
		ok = false
	case c == NEWLINE:
		ok = false
		scanner.newline()
	// handle single character tokens
	case c == '(': tok = token.Token{Type: token.LEFT_PAREN, Lexeme: "(", Line: scanner.line, Col: scanner.col}
	case c == ')': tok = token.Token{Type: token.RIGHT_PAREN, Lexeme: ")", Line: scanner.line, Col: scanner.col}
	case c == '{': tok = token.Token{Type: token.LEFT_BRACE, Lexeme: "{", Line: scanner.line, Col: scanner.col}
	case c == '}': tok = token.Token{Type: token.RIGHT_BRACE, Lexeme: "}", Line: scanner.line, Col: scanner.col}
	case c == '[': tok = token.Token{Type: token.LEFT_SQUARE, Lexeme: "[", Line: scanner.line, Col: scanner.col}
	case c == ']': tok = token.Token{Type: token.RIGHT_SQUARE, Lexeme: "]", Line: scanner.line, Col: scanner.col}
	case c == ',': tok = token.Token{Type: token.COMMA, Lexeme: ",", Line: scanner.line, Col: scanner.col}
	case c == '+': tok = token.Token{Type: token.PLUS, Lexeme: "+", Line: scanner.line, Col: scanner.col}
	case c == '*': tok = token.Token{Type: token.MULT, Lexeme: "*", Line: scanner.line, Col: scanner.col}
	case c == ';': tok = token.Token{Type: token.SEMICOLON, Lexeme: ";", Line: scanner.line, Col: scanner.col}
	case c == ':': tok = token.Token{Type: token.COLON, Lexeme: ":", Line: scanner.line, Col: scanner.col}
	case c == '%': tok = token.Token{Type: token.PERCENT, Lexeme: "%", Line: scanner.line, Col: scanner.col}
	// handle two character tokens
	case c == '-' && scanner.advanceIfAtr('>'): tok = token.Token{Type: token.MINUS_GREAT, Lexeme: "<-", Line: scanner.line, Col: scanner.col}
	case c == '-': 								tok = token.Token{Type: token.MINUS, Lexeme: "-", Line: scanner.line, Col: scanner.col}
	case c == '/' && scanner.advanceIfAtr('/'):	ok, tok = scanner.skipComment()
	case c == '/':								tok = token.Token{Type: token.SLASH, Lexeme: "/", Line: scanner.line, Col: scanner.col}
	case c == '<' && scanner.advanceIfAtr('='):	tok = token.Token{Type: token.LESS_EQ, Lexeme: "<=", Line: scanner.line, Col: scanner.col}
	case c == '<' && scanner.advanceIfAtr('<'):	tok = token.Token{Type: token.LESS_LESS, Lexeme: "<<", Line: scanner.line, Col: scanner.col}
	case c == '<' && scanner.advanceIfAtr('-'): tok = token.Token{Type: token.LESS_MINUS, Lexeme: "<-", Line: scanner.line, Col: scanner.col}
	case c == '<':								tok = token.Token{Type: token.LESS, Lexeme: "<", Line: scanner.line, Col: scanner.col}
	case c == '>' && scanner.advanceIfAtr('='):	tok = token.Token{Type: token.GREATER_EQ, Lexeme: ">=", Line: scanner.line, Col: scanner.col}
	case c == '>' && scanner.advanceIfAtr('>'):
		if scanner.advanceIfAtr('='){
			tok = token.Token{Type: token.GREATER_GREATER_EQ, Lexeme: ">>=", Line: scanner.line, Col: scanner.col}
		} else {
			tok = token.Token{Type: token.GREATER_GREATER, Lexeme: ">>", Line: scanner.line, Col: scanner.col}
		}									
	case c == '>' && scanner.advanceIfAtr('>'):	
	case c == '>':								tok = token.Token{Type: token.GREATER, Lexeme: ">", Line: scanner.line, Col: scanner.col}
	case c == '|' && scanner.advanceIfAtr('|'):	tok = token.Token{Type: token.BAR_BAR, Lexeme: "||", Line: scanner.line, Col: scanner.col}
	case c == '|':								tok = token.Token{Type: token.BAR, Lexeme: "|", Line: scanner.line, Col: scanner.col}
	case c == '&' && scanner.advanceIfAtr('&'):	tok = token.Token{Type: token.AMP_AMP, Lexeme: "&&", Line: scanner.line, Col: scanner.col}
	case c == '!' && scanner.advanceIfAtr('='):	tok = token.Token{Type: token.EXCLAMATION_EQ, Lexeme: "!=", Line: scanner.line, Col: scanner.col}
	case c == '!':								tok = token.Token{Type: token.EXCLAMATION, Lexeme: "!", Line: scanner.line, Col: scanner.col}
	case c == '=' && scanner.advanceIfAtr('='):	tok = token.Token{Type: token.EQ_EQ, Lexeme: "==", Line: scanner.line, Col: scanner.col}
	case c == '=' && scanner.advanceIfAtr('>'): tok = token.Token{Type: token.EQ_GREATER, Lexeme: "=>", Line: scanner.line, Col: scanner.col}
	case c == '=':								tok = token.Token{Type: token.EQ, Lexeme: "=", Line: scanner.line, Col: scanner.col}
	// handle string literals
	case c == '"':								ok, tok = scanner.scanString()
	// handle numbers
	case c == '.' && scanner.advanceIfAtr('.'):	tok = token.Token{Type: token.DOT_DOT, Lexeme: "..", Line: scanner.line, Col: scanner.col}
	case c == '.' && scanner.atrIsDigit():		ok, tok = scanner.scanNumber()
	case c == '.':								tok = token.Token{Type: token.DOT, Lexeme: ".", Line: scanner.line, Col: scanner.col}
	case utils.IsDigit(c):							ok, tok = scanner.scanNumber()
	// handle identifiers
	case utils.IsAlpha(c):							ok, tok = scanner.scanIdentifierOrKeyword()
	default:
		ok, tok = false, token.Token{Type: token.UNKNOWN, Lexeme: "UNK", Line: scanner.line, Col: scanner.col}
		scanner.error("Unrecognized character", scanner.line, scanner.col)
	}

	// Set left idx to right idx for next token
	scanner.l = scanner.r

	return tok, ok
}

func (scanner *MSScanner) scanIdentifierOrKeyword() (bool, token.Token) {

	for !scanner.atEnd() && (utils.IsAlpha(scanner.atr()) || utils.IsDigit(scanner.atr())) {
		scanner.advance()
	}

	// extract string
	str := scanner.src[scanner.l:scanner.r]

	// check if the string is a keyword
	if tt, ok := token.Keywords[str]; ok {
		return true, token.Token{Type: tt, Lexeme: str, Line: scanner.line, Col: scanner.col}
	}

	// not a keyword, so it is an identifier
	return true, token.Token{Type: token.IDENTIFIER, Lexeme: str, Line: scanner.line, Col: scanner.col}

}

func (scanner *MSScanner) scanString() (bool, token.Token) {
	// advance r untill we find the matching "
	// make sure we don't go past the end of the file
	// also make sure to increment newlines occuring
	for scanner.atr() != QUOTE && !scanner.atEnd() {
		if scanner.atr() == NEWLINE { scanner.newline() }
		scanner.advance()
	}

	// Check cause of loop exit, if from EOF we have an error
	if scanner.atEnd() {
		scanner.error("No matching \" found for string", scanner.line, scanner.col)
		return false, token.Token{}
	}
	
	// Found the closing quote, add the string token
	str := scanner.src[scanner.l+1:scanner.r]
	tok := token.Token{Type: token.STRING, Lexeme: str, Line: scanner.line, Col: scanner.col}

	// advance past the closing quote
	scanner.advance()

	return true, tok
}

func (scanner *MSScanner) scanNumber() (bool, token.Token) {

	// scanner.l points to either the first digit or the dot
	// scanner.r points to next character after the first digit or dot
	// example: 242.213
	//          ^^
	//          lr
	// We know for a factc that r is also digit

	// keep track if the number contains a dot,
	// can also be the first character
	var ndot int = 0
	if scanner.atl() == '.' { ndot = ndot + 1 }

	valid := true

	// advance r until we find a non-digit character
	// make sure we don't go past the end of the file
	for !scanner.atEnd() {

		// cif there is a dot, we need to increment ndot
		// if we don't have a digit, we break (end of number)
		if scanner.atr() == '.' {
			ndot = ndot + 1
		} else if scanner.atr() == SPACE || scanner.atr() == TAB {
			break
		} else if utils.IsAlpha(scanner.atr()) {
			// Found a non-digit character, we have an error
			// And we know it is not a space, tab or newline
			// But we still continue the loop to find the end of the number
			valid = false
		} else if !utils.IsDigit(scanner.atr()) {
			// Not a digit, space, tab or newline, but also not
			// an alpha character, so this is still a valid number
			// Example: {123}; is valid and 42; is valid
			break
		}

		// advance r
		scanner.advance()
	}

	// Check if the number contains more than one dot
	// if it does, we have an error
	if ndot > 1 {
		scanner.error("Invalid number literal", scanner.line, scanner.col)
		return false, token.Token{}
	}

	if !valid {
		scanner.error("Invalid number literal", scanner.line, scanner.col)
		return false, token.Token{}
	}

	// Check if the character at r is a quote '"', this is not allowed
	if scanner.atr() == QUOTE {
		scanner.error("Invalid number literal", scanner.line, scanner.col)
		return false, token.Token{}
	}


	// should be space, newline, tab or end of file
	// so we can add the number token
	var tt token.TokenType
	if ndot == 1 {
		tt = token.NUMBER_FLOAT
	} else {
		tt = token.NUMBER_INT
	}

	// create the number token
	num := scanner.src[scanner.l:scanner.r]
	tok := token.Token{Type: tt, Lexeme: num, Line: scanner.line, Col: scanner.col}

	return true, tok
}

func (scanner *MSScanner) skipComment() (bool, token.Token) {

	// found a comment where l points to the first /
	// and r points to the second / Now we need to advance
	// r until we find a newline or the end of the file
	for !scanner.atEnd() && scanner.atr() != NEWLINE {
			scanner.advance()
	}

	return false, token.Token{}
}

func (scanner *MSScanner) advanceIfAtr(c byte) bool {

	// check if we are at the end of the file
	if scanner.atEnd() { return false }
	if scanner.atr() != c { return false }

	// increment r and return true
	scanner.advance()

	return true
}

func (scanner *MSScanner) addToken(token token.Token) {
	scanner.tokens = append(scanner.tokens, token)
}

func (scanner *MSScanner) atEnd() bool {
	return scanner.r >= scanner.n
}

func (scanner *MSScanner) error(msg string, line int, col int) {
	err := ScannerError{msg, line, col}
	scanner.Errors = append(scanner.Errors, err)
}

func (scanner *MSScanner) setSrc(input string) {
	scanner.src = input
	scanner.n = len(input)
}

func (scanner *MSScanner) reset() {

	// reset src information
	scanner.src = ""
	scanner.n = 0

	// reset scanner state
	scanner.tokens = make([]token.Token, 0)
	scanner.r = 0
	scanner.l = 0
	scanner.line = 1
	scanner.col = 1

	// reset errors
	scanner.Errors = make([]ScannerError, 0)
	
}