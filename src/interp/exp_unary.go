package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)

func (evaluator *MSEvaluator) evaluateUnaryExpression(node *ast.UnaryExpNodeS) (MSVal, error) {
	
	// evaluate the node
	res, err := evaluator.evaluateExpression(node.Node)

	// check if the evaluation was.Valid()
	if err != nil {
		return MSNothing{}, err
	}

	// handle unary operators
	switch node.Op.Type {
	case token.MINUS:		return evaluateMinus(res)
	case token.EXCLAMATION:	return evaluateExcl(res)
	default: 				return MSNothing{}, &EvalError{unknownUnop(node.Op.Lexeme, res)}
	}
	
}

func evaluateMinus(res MSVal) (MSVal, error) {

	var err error

	switch v := res.(type){
	case MSInt:	return MSInt{Val: -v.Val}, err
	case MSFloat:	return MSFloat{Val: -v.Val}, err
	default:			return MSNothing{}, &EvalError{unknownUnop(token.MINUS.String(), res)}
	}
}

func evaluateExcl(res MSVal) (MSVal, error) {

	var err error

	switch v := res.(type){
	case MSBool:	return MSBool{Val: !v.Val}, err
	default:			return MSNothing{}, &EvalError{unknownUnop(token.EXCLAMATION.String(), res)}
	}
}

func (evaluator *MSEvaluator) evaluateGroupExpression(node *ast.GroupExpNodeS) (MSVal, error) {
	return evaluator.evaluateExpression(node.Node)
}


func unknownUnop(lexeme string, tt MSVal) string {
	return fmt.Sprintf("Operator '%v' is not defined for type '%v'", lexeme, tt.Type())
}