package main

import "testing"

func tokenCompare(a, b Token) bool {
	return a.Type == b.Type && a.Lexeme == b.Lexeme
}

func TestScan(t *testing.T) {
	
	// test cases
	tests := []struct {
		input string
		tokens []Token
	}{
		{
			input: "123;",
			tokens: []Token{
				{Type: NUMBER_INT, Lexeme: "123"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "{.1};",
			tokens: []Token{
				{Type: LEFT_BRACE, Lexeme: "{"},
				{Type: NUMBER_FLOAT, Lexeme: ".1"},
				{Type: RIGHT_BRACE, Lexeme: "}"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "1/2;",
			tokens: []Token{
				{Type: NUMBER_INT, Lexeme: "1"},
				{Type: SLASH, Lexeme: "/"},
				{Type: NUMBER_INT, Lexeme: "2"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "123+456.789",
			tokens: []Token{
				{Type: NUMBER_INT, Lexeme: "123"},
				{Type: PLUS, Lexeme: "+"},
				{Type: NUMBER_FLOAT, Lexeme: "456.789"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "1 >> int x; \n",
			tokens: []Token{
				{Type: NUMBER_INT, Lexeme: "1"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: INT_TYPE, Lexeme: "int"},
				{Type: IDENTIFIER, Lexeme: "x"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "1, 2 >> +, y >> +;",
			tokens: []Token{
				{Type: NUMBER_INT, Lexeme: "1"},
				{Type: COMMA, Lexeme: ","},
				{Type: NUMBER_INT, Lexeme: "2"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: PLUS, Lexeme: "+"},
				{Type: COMMA, Lexeme: ","},
				{Type: IDENTIFIER, Lexeme: "y"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: PLUS, Lexeme: "+"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: "true >> bool x;",
			tokens: []Token{
				{Type: TRUE, Lexeme: "true"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: BOOLEAN_TYPE, Lexeme: "bool"},
				{Type: IDENTIFIER, Lexeme: "x"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},

		{
			input: "\"hello\" >> string x >> print;",
			tokens: []Token{
				{Type: STRING, Lexeme: "hello"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: STRING_TYPE, Lexeme: "string"},
				{Type: IDENTIFIER, Lexeme: "x"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: IDENTIFIER, Lexeme: "print"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: `
				xif {
				| x, true >> == {"Hello" >> print;}
				| otherwise {"Goodbye" >> print;}
				}`,
			tokens: []Token{
				{Type: XIF, Lexeme: "xif"},
				{Type: LEFT_BRACE, Lexeme: "{"},
				{Type: BAR, Lexeme: "|"},
				{Type: IDENTIFIER, Lexeme: "x"},
				{Type: COMMA, Lexeme: ","},
				{Type: TRUE, Lexeme: "true"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: EQ_EQ, Lexeme: "=="},
				{Type: LEFT_BRACE, Lexeme: "{"},
				{Type: STRING, Lexeme: "Hello"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: IDENTIFIER, Lexeme: "print"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: RIGHT_BRACE, Lexeme: "}"},
				{Type: BAR, Lexeme: "|"},
				{Type: OTHERWISE, Lexeme: "otherwise"},
				{Type: LEFT_BRACE, Lexeme: "{"},
				{Type: STRING, Lexeme: "Goodbye"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: IDENTIFIER, Lexeme: "print"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: RIGHT_BRACE, Lexeme: "}"},
				{Type: RIGHT_BRACE, Lexeme: "}"},
				{Type: EOF, Lexeme: ""},
			},
		},
		{
			input: `for [1 .. 10] >> int i{
						i >> print;
					}`,
			tokens: []Token{
				{Type: FOR, Lexeme: "for"},
				{Type: LEFT_SQUARE, Lexeme: "["},
				{Type: NUMBER_INT, Lexeme: "1"},
				{Type: DOT_DOT, Lexeme: ".."},
				{Type: NUMBER_INT, Lexeme: "10"},
				{Type: RIGHT_SQUARE, Lexeme: "]"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: INT_TYPE, Lexeme: "int"},
				{Type: IDENTIFIER, Lexeme: "i"},
				{Type: LEFT_BRACE, Lexeme: "{"},
				{Type: IDENTIFIER, Lexeme: "i"},
				{Type: GREATER_GREATER, Lexeme: ">>"},
				{Type: IDENTIFIER, Lexeme: "print"},
				{Type: SEMICOLON, Lexeme: ";"},
				{Type: RIGHT_BRACE, Lexeme: "}"},
				{Type: EOF, Lexeme: ""},
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
	received = scanner.errors
	if !arraysEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}

	///////////////////////////////////////////////

	input = "hello \"world"
	expected = []ScannerError{
		{msg: "No matching \" found for string", line: 1, col: 13},
	}

	scanner.Scan(input)
	received = scanner.errors
	if !arraysEqual(received, expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}

	///////////////////////////////////////////////

	input = "?"
	expected = []ScannerError{
		{msg: "Unrecognized character", line: 1, col: 2},
	}

	scanner.Scan(input)
	received = scanner.errors
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