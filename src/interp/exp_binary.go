package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func (evaluator *MSEvaluator) evaluateBinaryExpression(node *ast.BinaryExpNodeS) EvalResult {

	// Evaluate left and right side of the expression
	left := evaluator.evaluateExpression(&node.Left)
	right := evaluator.evaluateExpression(&node.Right)

	// If either of the evaluations failed, return the error
	if !left.Valid() || !right.Valid() {
		return EvalResult{err: append(left.err, right.err...)}
	}

	switch node.Op.Type {
	case token.PLUS: 				return evalAdd(left, right)
	case token.MINUS:				return evalSub(left, right)
	case token.MULT:				return evalMult(left, right)
	case token.SLASH:				return evalDiv(left, right)
	case token.GREATER:				return evalGreater(left, right, node.Op.Type)
	case token.LESS:				return evalGreater(right, left, node.Op.Type)
	case token.GREATER_EQ:			return evalGreaterEq(left, right, node.Op.Type)
	case token.LESS_EQ:				return evalGreaterEq(right, left, node.Op.Type)
	case token.EQ_EQ:				return evalEq(left, right, node.Op.Type)
	case token.EXCLAMATION_EQ:		return evalNeq(left, right, node.Op.Type)
	case token.GREATER_GREATER: 	return evalGrGr(left, right)
	case token.COMMA:				return evalTuple(left, right)
	case token.PERCENT:				return evalMod(left, right)
	default:						return evalErr(fmt.Sprintf("Unknown binary operator: %v", node.Op.Lexeme))
	}
}


func (evaluator *MSEvaluator) evaluateLogicalExpression(node *ast.LogicalExpNodeS) EvalResult {

	// Evaluate the left and right side of the expression
	left := evaluator.evaluateExpression(&node.Left)

	// If either of the evaluations failed, return the error
	if !left.Valid() {
		return left // EvalResult contains error(s)
	}

	// Check if the left side is a boolean
	if !left.IsType(&mstype.MS_BOOL) {
		return evalErr(invalidLogop(left, node.Op.Lexeme))
	}

	// Check the truth value of the left side
	leftb := left.val.(bool)

	// short circuit evaluation
	switch  {
	case node.Op.Type == token.AMP_AMP && !leftb: return EvalResult{rt: mstype.MS_BOOL, val: false} // false && ...
	case node.Op.Type == token.BAR_BAR &&  leftb: return EvalResult{rt: mstype.MS_BOOL, val: true}  // true || ...
	}

	// Means the first operand is inconclusive
	right := evaluator.evaluateExpression(&node.Right)

	if !right.Valid() {
		return right // EvalResult contains error(s)
	}

	// Check if the right side is a boolean
	if !right.IsType(&mstype.MS_BOOL) {
		return evalErr(invalidLogop(right, node.Op.Lexeme))
	}

	return right

}

func invalidBinop(left, right EvalResult, op string) string {
	return fmt.Sprintf("Operator '%v' is not defined for types '%v' and '%v'", op, left.rt, right.rt)
}

func invalidLogop(left EvalResult, op string) string {
	return fmt.Sprintf("Logical operator '%v' is not defined for type '%v'", op, left.rt)
}


func evalTuple(left, right EvalResult) EvalResult {
	valuetype := mstype.MSCompositeTypeS{
		Left: left.rt,
		Right: right.rt,
	}
	return EvalResult{rt: &valuetype, val: []EvalResult{left, right}}
}

func evalMod(left, right EvalResult) EvalResult {
	switch left.rt {
	case mstype.MS_INT:
		switch right.rt {
		case mstype.MS_INT:
			if right.val.(int) == 0 { return evalErr("Division by zero.")}
			return EvalResult{rt: mstype.MS_INT, val: left.val.(int) % right.val.(int)}
		}
	}
	return evalErr(invalidBinop(left, right, "%"))
}