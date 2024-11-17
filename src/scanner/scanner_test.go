package scanner

import (
	token "mikescript/src/token"
	"testing"
)

func tokenCompare(a, b token.Token) bool {
	return a.Type == b.Type && a.Lexeme == b.Lexeme
}

func TestScan(t *testing.T) {
	
	// test cases
	tests := []struct {
		input string
		tokens []token.Token
	}{
		{
			input: "123;",
			tokens: []token.Token{
				{Type: token.NUMBER_INT, Lexeme: "123"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "{.1};",
			tokens: []token.Token{
				{Type: token.LEFT_BRACE, Lexeme: "{"},
				{Type: token.NUMBER_FLOAT, Lexeme: ".1"},
				{Type: token.RIGHT_BRACE, Lexeme: "}"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "1/2;",
			tokens: []token.Token{
				{Type: token.NUMBER_INT, Lexeme: "1"},
				{Type: token.SLASH, Lexeme: "/"},
				{Type: token.NUMBER_INT, Lexeme: "2"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "123+456.789",
			tokens: []token.Token{
				{Type: token.NUMBER_INT, Lexeme: "123"},
				{Type: token.PLUS, Lexeme: "+"},
				{Type: token.NUMBER_FLOAT, Lexeme: "456.789"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "1 >> int x; \n",
			tokens: []token.Token{
				{Type: token.NUMBER_INT, Lexeme: "1"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.INT_TYPE, Lexeme: "int"},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "1, 2 >> +, y >> +;",
			tokens: []token.Token{
				{Type: token.NUMBER_INT, Lexeme: "1"},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.NUMBER_INT, Lexeme: "2"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.PLUS, Lexeme: "+"},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.IDENTIFIER, Lexeme: "y"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.PLUS, Lexeme: "+"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: "true >> bool x;",
			tokens: []token.Token{
				{Type: token.TRUE, Lexeme: "true"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.BOOLEAN_TYPE, Lexeme: "bool"},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},

		{
			input: "\"hello\" >> string x >> print;",
			tokens: []token.Token{
				{Type: token.STRING, Lexeme: "hello"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.STRING_TYPE, Lexeme: "string"},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: `
				xif {
				| x, true >> == {"Hello" >> print;}
				| otherwise {"Goodbye" >> print;}
				}`,
			tokens: []token.Token{
				{Type: token.XIF, Lexeme: "xif"},
				{Type: token.LEFT_BRACE, Lexeme: "{"},
				{Type: token.BAR, Lexeme: "|"},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.TRUE, Lexeme: "true"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.EQ_EQ, Lexeme: "=="},
				{Type: token.LEFT_BRACE, Lexeme: "{"},
				{Type: token.STRING, Lexeme: "Hello"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.RIGHT_BRACE, Lexeme: "}"},
				{Type: token.BAR, Lexeme: "|"},
				{Type: token.OTHERWISE, Lexeme: "otherwise"},
				{Type: token.LEFT_BRACE, Lexeme: "{"},
				{Type: token.STRING, Lexeme: "Goodbye"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.RIGHT_BRACE, Lexeme: "}"},
				{Type: token.RIGHT_BRACE, Lexeme: "}"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
		{
			input: `for [1 .. 10] >> int i{
						i >> print;
					}`,
			tokens: []token.Token{
				{Type: token.FOR, Lexeme: "for"},
				{Type: token.LEFT_SQUARE, Lexeme: "["},
				{Type: token.NUMBER_INT, Lexeme: "1"},
				{Type: token.DOT_DOT, Lexeme: ".."},
				{Type: token.NUMBER_INT, Lexeme: "10"},
				{Type: token.RIGHT_SQUARE, Lexeme: "]"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.INT_TYPE, Lexeme: "int"},
				{Type: token.IDENTIFIER, Lexeme: "i"},
				{Type: token.LEFT_BRACE, Lexeme: "{"},
				{Type: token.IDENTIFIER, Lexeme: "i"},
				{Type: token.GREATER_GREATER, Lexeme: ">>"},
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.RIGHT_BRACE, Lexeme: "}"},
				{Type: token.EOF, Lexeme: ""},
			},
		},
	}

	scanner := MSScanner{}

	for _, test := range tests {
		
		// scan the input
		tokens := scanner.Scan(test.input)

		// check the tokens
		if len(tokens) != len(test.tokens) {
			t.Errorf("Expected %d tokens, got %d", len(test.tokens), len(tokens))
		}

		for i := 0; i < len(tokens); i++ {
			if !tokenCompare(tokens[i], test.tokens[i]) {
				t.Errorf("Expected %v, got %v", test.tokens[i], tokens[i])
			}
		}
	}

}

func TestScannerErrors(t *testing.T) {

	var input string
	var expected []ScannerError
	var scanner MSScanner = MSScanner{}
	var received []ScannerError

	///////////////////////////////////////////////
	input = "1.2.3"
	expected = []ScannerError{
		{msg: "Invalid number literal", line: 1, col: 6},
	}

	scanner.Scan(input)
	received = scanner.Errors
	if !arraysEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}

	///////////////////////////////////////////////

	input = "hello \"world"
	expected = []ScannerError{
		{msg: "No matching \" found for string", line: 1, col: 13},
	}

	scanner.Scan(input)
	received = scanner.Errors
	if !arraysEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}

	///////////////////////////////////////////////

	input = "?"
	expected = []ScannerError{
		{msg: "Unrecognized character", line: 1, col: 2},
	}

	scanner.Scan(input)
	received = scanner.Errors
	if !arraysEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}


}

func arraysEqual(a, b []ScannerError) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !a[i].Compare(b[i]) {
			return false
		}
	}

	return true
}