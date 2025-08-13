package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)

func (evaluator *MSEvaluator) evaluateBinaryExpression(node *ast.BinaryExpNodeS) (MSVal, error) {

	// Evaluate left and right side of the expression
	lval, lerr := evaluator.evaluateExpression(node.Left)

	// quit on error
	if lerr != nil {
		return nil, lerr
	}

	rval, rerr := evaluator.evaluateExpression(node.Right)

	// quit on error
	if rerr != nil {
		return nil, rerr
	}

	switch node.Op.Type {
	case token.PLUS: 				return evalAdd(lval, rval)
	case token.MINUS:				return evalSub(lval, rval)
	case token.MULT:				return evalMult(lval, rval)
	case token.SLASH:				return evalDiv(lval, rval)
	case token.GREATER:				return evalGreater(lval, rval, node.Op.Type)
	case token.LESS:				return evalGreater(rval, lval, node.Op.Type)
	case token.GREATER_EQ:			return evalGreaterEq(lval, rval, node.Op.Type)
	case token.LESS_EQ:				return evalGreaterEq(rval, lval, node.Op.Type)
	case token.EQ_EQ:				return evalEq(lval, rval, node.Op.Type)
	case token.EXCLAMATION_EQ:		return evalNeq(lval, rval, node.Op.Type)
	case token.GREATER_GREATER: 	return evalGrGr(lval, rval)
	case token.COMMA:				return evalTuple(lval, rval)
	case token.PERCENT:				return evalMod(lval, rval)
	default:						return nil, &EvalError{unknownBinop(node)}
	}
}

func invalidBinop(left, right MSVal, op string) string {
	return fmt.Sprintf("Operator '%v' is not defined for types '%v' and '%v'", op, left, right)
}

func invalidLogop(left MSVal, op string) string {
	return fmt.Sprintf("Logical operator '%v' is not defined for type '%v'", op, left)
}

func unknownBinop(node *ast.BinaryExpNodeS) string {
	return fmt.Sprintf("Unknown binary operator: %v", node.Op.Lexeme)
}

func evalTuple(left, right MSVal) (MSVal, error) {
	return MSTuple{Values: []MSVal{left, right}}, nil
}

func evalMod(lval, rval MSVal) (MSVal, error) {

	var err error

	switch l := lval.(type) {
	case MSInt:
		switch r := rval.(type){
		case MSInt:
			if r.Val == 0 {
				return nil, &EvalError{message: "Division by zero."}
			}
			return MSInt{Val: l.Val % r.Val}, err
		}
	}

	return nil, &EvalError{invalidBinop(lval, rval, "%")}
}