package interp

import (
	"mikescript/src/ast"
	"mikescript/src/token"
)



func (evaluator *MSEvaluator) evaluateLogicalExpression(node *ast.LogicalExpNodeS) (MSVal, error) {

	var err error

	// Evaluate the left and right side of the expression
	lval, lerr := evaluator.evaluateExpression(&node.Left)

	// error check
	if lerr != nil {
		return MSNothing{}, lerr
	}

	// Expect boolean value
	var leftb bool

	switch lt := lval.(type){
	case MSBool:	leftb = lt.Val
	default:			return MSNothing{}, &EvalError{invalidLogop(lval, node.Op.Lexeme)}
	}

	// short circuit evaluation
	switch  {
	case node.Op.Type == token.AMP_AMP && !leftb: return MSBool{Val: false}, err	// false && ...
	case node.Op.Type == token.BAR_BAR &&  leftb: return MSBool{Val: true}, err	// true || ...
	}

	// Means the first operand is inconclusive
	rval, rerr := evaluator.evaluateExpression(&node.Right)

	// error check
	if rerr != nil {
		return MSNothing{}, rerr
	}

	var rb bool

	switch rt := rval.(type){
	case MSBool:	rb = rt.Val
	default:			return MSNothing{}, &EvalError{invalidLogop(lval, node.Op.Lexeme)}
	}

	return MSBool{Val: rb}, err

}
