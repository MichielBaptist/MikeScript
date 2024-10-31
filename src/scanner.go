package main

import (
	"fmt"
)


const(
	SPACE byte = ' '
	TAB byte = '\t'
	NEWLINE byte = '\n'
	QUOTE byte = '"'
)

type Scanner interface {
	Scan(input string) []Token
}

type MSScanner struct {

	// src information
	src string 		// Source code to scan
	n int 			// Length of source code

	// scanner state
	tokens []Token 	// Tokens found in source code
	l int 			// Start of current token
	r int 			// Current position in source code
	line int 		// Current line number
	col int 		// Current column number

	// error information
	errors []ScannerError
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
	return IsDigit(scanner.atr())
}

////////////////////////////////////////////////////////////////
// 							Scan
////////////////////////////////////////////////////////////////

func (scanner *MSScanner) Scan(input string) []Token {
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
	scanner.addToken(Token{EOF, "", scanner.line, scanner.col})
}

func (scanner *MSScanner) nextToken() (Token, bool) {

	// get current char and advance r to next character
	c := scanner.advance()

	var tok Token      	// nil token
	var ok bool = true	// valid token

	switch {

	// handle witespace and newlines
	case c == SPACE || c == TAB:
		ok = false
	case c == NEWLINE:
		ok = false
		scanner.newline()
	// handle single character tokens
	case c == '(': tok = Token{LEFT_PAREN, "(", scanner.line, scanner.col}
	case c == ')': tok = Token{RIGHT_PAREN, ")", scanner.line, scanner.col}
	case c == '{': tok = Token{LEFT_BRACE, "{", scanner.line, scanner.col}
	case c == '}': tok = Token{RIGHT_BRACE, "}", scanner.line, scanner.col}
	case c == '[': tok = Token{LEFT_SQUARE, "[", scanner.line, scanner.col}
	case c == ']': tok = Token{RIGHT_SQUARE, "]", scanner.line, scanner.col}
	case c == ',': tok = Token{COMMA, ",", scanner.line, scanner.col}
	case c == '+': tok = Token{PLUS, "+", scanner.line, scanner.col}
	case c == '*': tok = Token{MULT, "*", scanner.line, scanner.col}
	case c == ';': tok = Token{SEMICOLON, ";", scanner.line, scanner.col}
	case c == '%': tok = Token{PERCENT, "%", scanner.line, scanner.col}
	// handle two character tokens
	case c == '-' && scanner.advanceIfAtr('>'): tok = Token{MINUS_GREAT, "<-", scanner.line, scanner.col}
	case c == '-': 								tok = Token{MINUS, "-", scanner.line, scanner.col}
	case c == '/' && scanner.advanceIfAtr('/'):	ok, tok = scanner.skipComment()
	case c == '/':								tok = Token{SLASH, "/", scanner.line, scanner.col}
	case c == '<' && scanner.advanceIfAtr('='):	tok = Token{LESS_EQ, "<=", scanner.line, scanner.col}
	case c == '<' && scanner.advanceIfAtr('<'):	tok = Token{LESS_LESS, "<<", scanner.line, scanner.col}
	case c == '<' && scanner.advanceIfAtr('-'): tok = Token{LESS_MINUS, "<-", scanner.line, scanner.col}
	case c == '<':								tok = Token{LESS, "<", scanner.line, scanner.col}
	case c == '>' && scanner.advanceIfAtr('='):	tok = Token{GREATER_EQ, ">=", scanner.line, scanner.col}
	case c == '>' && scanner.advanceIfAtr('>'):	tok = Token{GREATER_GREATER, ">>", scanner.line, scanner.col}
	case c == '>':								tok = Token{GREATER, ">", scanner.line, scanner.col}
	case c == '|' && scanner.advanceIfAtr('|'):	tok = Token{BAR_BAR, "||", scanner.line, scanner.col}
	case c == '|':								tok = Token{BAR, "|", scanner.line, scanner.col}
	case c == '!' && scanner.advanceIfAtr('='):	tok = Token{EXCLAMATION_EQ, "!=", scanner.line, scanner.col}
	case c == '!':								tok = Token{EXCLAMATION, "!", scanner.line, scanner.col}
	case c == '=' && scanner.advanceIfAtr('='):	tok = Token{EQ_EQ, "==", scanner.line, scanner.col}
	case c == '=':								tok = Token{EQ, "=", scanner.line, scanner.col}
	// handle string literals
	case c == '"':								ok, tok = scanner.scanString()
	// handle numbers
	case c == '.' && scanner.advanceIfAtr('.'):	tok = Token{DOT_DOT, "..", scanner.line, scanner.col}
	case c == '.' && scanner.atrIsDigit():		ok, tok = scanner.scanNumber()
	case c == '.':								tok = Token{DOT, ".", scanner.line, scanner.col}
	case IsDigit(c):							ok, tok = scanner.scanNumber()
	// handle identifiers
	case IsAlpha(c):							ok, tok = scanner.scanIdentifierOrKeyword()
	default:
		ok, tok = false, Token{UNKNOWN, "UNK", scanner.line, scanner.col}
		scanner.error("Unrecognized character", scanner.line, scanner.col)
	}

	// Set left idx to right idx for next token
	scanner.l = scanner.r

	return tok, ok
}

func (scanner *MSScanner) scanIdentifierOrKeyword() (bool, Token) {

	for !scanner.atEnd() && (IsAlpha(scanner.atr()) || IsDigit(scanner.atr())) {
		scanner.advance()
	}

	// extract string
	str := scanner.src[scanner.l:scanner.r]

	// check if the string is a keyword
	if tt, ok := keywords[str]; ok {
		return true, Token{tt, str, scanner.line, scanner.col}
	}

	// not a keyword, so it is an identifier
	return true, Token{IDENTIFIER, str, scanner.line, scanner.col}

}

func (scanner *MSScanner) scanString() (bool, Token) {
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
		return false, Token{}
	}
	
	// Found the closing quote, add the string token
	str := scanner.src[scanner.l+1:scanner.r]
	tok := Token{STRING, str, scanner.line, scanner.col}

	// advance past the closing quote
	scanner.advance()

	return true, tok
}

func (scanner *MSScanner) scanNumber() (bool, Token) {

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
		} else if IsAlpha(scanner.atr()) {
			// Found a non-digit character, we have an error
			// And we know it is not a space, tab or newline
			// But we still continue the loop to find the end of the number
			valid = false
		} else if !IsDigit(scanner.atr()) {
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
		return false, Token{}
	}

	if !valid {
		scanner.error("Invalid number literal", scanner.line, scanner.col)
		return false, Token{}
	}

	// Check if the character at r is a quote '"', this is not allowed
	if scanner.atr() == QUOTE {
		scanner.error("Invalid number literal", scanner.line, scanner.col)
		return false, Token{}
	}


	// should be space, newline, tab or end of file
	// so we can add the number token
	var tt TokenType
	if ndot == 1 {
		tt = NUMBER_FLOAT
	} else {
		tt = NUMBER_INT
	}

	// create the number token
	num := scanner.src[scanner.l:scanner.r]
	tok := Token{tt, num, scanner.line, scanner.col}

	return true, tok
}

func (scanner *MSScanner) skipComment() (bool, Token) {

	// found a comment where l points to the first /
	// and r points to the second / Now we need to advance
	// r until we find a newline or the end of the file
	for !scanner.atEnd() && scanner.atr() != NEWLINE {
			scanner.advance()
	}

	return false, Token{}
}

func (scanner *MSScanner) advanceIfAtr(c byte) bool {

	// check if we are at the end of the file
	if scanner.atEnd() { return false }
	if scanner.atr() != c { return false }

	// increment r and return true
	scanner.advance()

	return true
}

func (scanner *MSScanner) addToken(token Token) {
	scanner.tokens = append(scanner.tokens, token)
}

func (scanner *MSScanner) atEnd() bool {
	return scanner.r >= scanner.n
}

func (scanner *MSScanner) error(msg string, line int, col int) {
	err := ScannerError{msg, line, col}
	scanner.errors = append(scanner.errors, err)
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
	scanner.tokens = make([]Token, 0)
	scanner.r = 0
	scanner.l = 0
	scanner.line = 1
	scanner.col = 1

	// reset errors
	scanner.errors = make([]ScannerError, 0)
	
}