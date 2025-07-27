package parser

import "fmt"

type ParserError struct {
	// Represents a parser error
	msg  string
	line int
	col  int
}

func (err ParserError) Error() string {
	return fmt.Sprintf("Parsing Error: %v at line %v col %v", err.msg, err.line, err.col)
}
