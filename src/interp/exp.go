package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) evaluateExpression(node ast.ExpNodeI) (MSVal, error) {
	switch node := node.(type) {
	case *ast.BinaryExpNodeS:				return evaluator.evaluateBinaryExpression(node)
	case *ast.UnaryExpNodeS:				return evaluator.evaluateUnaryExpression(node)
	case *ast.LiteralExpNodeS:				return evaluator.evaluateLiteralExpression(node)
	case *ast.GroupExpNodeS:				return evaluator.evaluateGroupExpression(node)
	case *ast.AssignmentNodeS:				return evaluator.evaluateAssignmentExpression(node)
	case *ast.DeclAssignNodeS: 				return evaluator.evaluateDeclAssignExpression(node)
	case *ast.FuncAppNodeS:					return evaluator.evaluateFunctionApplication(node)
	case *ast.FuncCallNodeS:				return evaluator.evaluateFunctionCall(node)
	case *ast.IterableFuncCallNodeS:		return evaluator.evaluateIterableFunctionCall(node)
	case *ast.IterableFuncAppNodeS:			return evaluator.evaluateIterableFunctionApplication(node)
	case *ast.IterableFuncAppAndCallNodeS:	return evaluator.evaluateIterableFunctionAppAndCall(node)
	case *ast.VariableExpNodeS:				return evaluator.evalVariable(node)
	case *ast.LogicalExpNodeS:				return evaluator.evaluateLogicalExpression(node)
	case *ast.ArrayIndexNodeS:				return evaluator.evalArrayIndexExpression(node)
	case *ast.ArrayConstructorNodeS:		return evaluator.evaluateArrayConstructor(node)
	case *ast.ArrayAssignmentNodeS:			return evaluator.evaluateArrayAssignment(node)
	case *ast.TupleNodeS:					return evaluator.evaluateTuple(node)
	case *ast.FieldAccessNodeS:				return evaluator.evaluateFieldAccess(node)
	case *ast.FieldAssignmentNode:			return evaluator.evaluateFieldAssign(node)
	case *ast.RangeConstructorNodeS:		return evaluator.evaluateRangeConstructor(node)
	case *ast.StarredExpNodeS:				return evaluator.evaluateStarredExpression(node)
	default:								return nil, &EvalError{fmt.Sprintf("Unknown expression type: '%#v'", node)}
	}
}

