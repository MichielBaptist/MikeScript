package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
	"strconv"
)

func (evaluator *MSEvaluator) evaluateLiteralExpression(node *ast.LiteralExpNodeS) (MSVal, error) {
	switch node.Tk.Type {
	case token.NUMBER_INT:		return evalIntLiteral(node)
	case token.NUMBER_FLOAT:	return evalFloatLiteral(node)
	case token.STRING:			return MSString{Val: node.Tk.Lexeme}, nil
	case token.TRUE:			return MSBool{Val: true}, nil
	case token.FALSE:			return MSBool{Val: false}, nil
	case token.NOTHING_TYPE:	return MSNothing{}, nil
	case token.IDENTIFIER:		return nil, &EvalError{fmt.Sprintf("Trying to evaluate identifier '%v' as a literal.", node.Tk.Lexeme)}
	default:					return nil, &EvalError{fmt.Sprintf("Literal type '%#v' is not defined.", node)}
	}
}

func evalIntLiteral(node *ast.LiteralExpNodeS) (MSVal, error) {
	// convert the lexeme to an int
	val, err := strconv.Atoi(node.Tk.Lexeme)

	// Should never happen if parser works correctly.
	if err != nil {
		return nil, &EvalError{fmt.Sprintf("Could not convert '%v' to 'int'", node.Tk.Lexeme)}
	}
	
	return MSInt{Val: val}, nil
}

func evalFloatLiteral(node *ast.LiteralExpNodeS) (MSVal, error) {
	// convert the lexeme to a float
	val, err := strconv.ParseFloat(node.Tk.Lexeme, 64)

	if err != nil {
		return nil, &EvalError{fmt.Sprintf("Could not convert '%v' to 'int'", node.Tk.Lexeme)}
	}

	return MSFloat{Val: val}, nil
}

