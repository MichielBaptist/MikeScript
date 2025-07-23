package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) evaluateExpression(node *ast.ExpNodeI) EvalResult {
	switch node := (*node).(type) {
	case ast.BinaryExpNodeS:	return evaluator.evaluateBinaryExpression(&node)
	case ast.UnaryExpNodeS:		return evaluator.evaluateUnaryExpression(&node)
	case ast.LiteralExpNodeS:	return evaluator.evaluateLiteralExpression(&node)
	case ast.GroupExpNodeS:		return evaluator.evaluateGroupExpression(&node)
	case ast.AssignmentNodeS:	return evaluator.evaluateAssignmentExpression(&node)
	case ast.DeclAssignNodeS: 	return evaluator.evaluateDeclAssignExpression(&node)
	case ast.FuncAppNodeS:		return evaluator.evaluateFunctionApplication(&node)
	case ast.VariableExpNodeS:	return evaluator.evalVariable(&node)
	case ast.LogicalExpNodeS:	return evaluator.evaluateLogicalExpression(&node)
	default:					return evalErr(fmt.Sprintf("Unknown expression type: %v", node))
	}
}