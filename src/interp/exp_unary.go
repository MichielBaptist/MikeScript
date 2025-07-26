package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (evaluator *MSEvaluator) evaluateUnaryExpression(node *ast.UnaryExpNodeS) EvalResult {
	// evaluate the node
	res := evaluator.evaluateExpression(&node.Node)

	// check if the evaluation was.Valid()
	if !res.Valid() {
		return res
	}

	// handle unary operators
	switch node.Op.Type {
	case token.MINUS:		return evaluateMinus(&res)
	case token.EXCLAMATION:	return evaluateExcl(&res)
	default: 				return evalErr(unknownUnop(node.Op.Lexeme, res.Rt))
	}
	
}

func evaluateMinus(res *EvalResult) EvalResult {
	switch {
	case res.IsType(&mstype.MS_INT): 	return EvalResult{Rt: mstype.MS_INT, Val: -res.Val.(int)}
	case res.IsType(&mstype.MS_FLOAT):	return EvalResult{Rt: mstype.MS_FLOAT, Val: -res.Val.(float64)}
	default:							return evalErr(unknownUnop(token.MINUS.String(), res.Rt))
	}
}

func evaluateExcl(res *EvalResult) EvalResult {
	switch {
	case res.IsType(&mstype.MS_BOOL):	return EvalResult{Rt: mstype.MS_BOOL, Val: !res.Val.(bool)}
	default:							return evalErr(unknownUnop(token.MINUS.String(), res.Rt))
	}
}

func (evaluator *MSEvaluator) evaluateGroupExpression(node *ast.GroupExpNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Node)
}


func unknownUnop(lexeme string, tt mstype.MSType) string {
	return fmt.Sprintf("Operator %v is not defined for type %v", lexeme, tt)
}