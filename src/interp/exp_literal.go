package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
	"strconv"
)

func (evaluator *MSEvaluator) evaluateLiteralExpression(node *ast.LiteralExpNodeS) EvalResult {
	switch node.Tk.Type {
	case token.NUMBER_INT:		return evalIntLiteral(node)
	case token.NUMBER_FLOAT:	return evalFloatLiteral(node)
	case token.STRING:			return evalStringLiteral(node)
	case token.TRUE:			return EvalResult{Rt: mstype.MS_BOOL, Val: true}
	case token.FALSE:			return EvalResult{Rt: mstype.MS_BOOL, Val: false}
	case token.IDENTIFIER:		return evalErr(fmt.Sprintf("Trying to evaluate identifier '%v' as a literal.", node.Tk.Lexeme))
	default:					return evalErr(fmt.Sprintf("Literal type '%v' is not defined.", node.Tk.Type))
	}
}

func evalIntLiteral(node *ast.LiteralExpNodeS) EvalResult {
	// convert the lexeme to an int
	val, err := strconv.Atoi(node.Tk.Lexeme)

	if err != nil {
		return evalErr(fmt.Sprintf("Could not convert '%v' to int.", node.Tk.Lexeme))
	}

	return EvalResult{Rt: mstype.MS_INT, Val: val}
}

func evalFloatLiteral(node *ast.LiteralExpNodeS) EvalResult {
	// convert the lexeme to a float
	val, err := strconv.ParseFloat(node.Tk.Lexeme, 64)

	if err != nil {
		return evalErr(fmt.Sprintf("Could not convert '%v' to float64.", node.Tk.Lexeme))
	}

	return EvalResult{Rt: mstype.MS_FLOAT, Val: val}
}

func evalStringLiteral(node *ast.LiteralExpNodeS) EvalResult {
	return EvalResult{Rt: mstype.MS_STRING, Val: node.Tk.Lexeme}
}

