package interp

import (
	"fmt"
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeStatement(node *ast.StmtNodeI) EvalResult {
	switch node := (*node).(type) {
	case ast.Program:			return evaluator.executeStatements(&node)
	case ast.DeclarationNodeS:	return evaluator.executeDeclarationStatement(&node)
	case ast.BlockNodeS:		return evaluator.executeBlock(&node)
	case ast.IfNodeS:			return evaluator.executeIfstatement(&node)
	case ast.WhileNodeS:		return evaluator.executeWhileStatement(&node)
	case ast.ContinueNodeS:		return EvalResult{rt: RT_CONTINUE}
	case ast.BreakNodeS:		return EvalResult{rt: RT_BREAK}
	case ast.ExStmtNodeS:		return evaluator.executeExpressionStatement(&node)
	default:					return evalErr(fmt.Sprintf("Unknown statement type: %v", node))
	}
}

func (evaluator *MSEvaluator) executeExpressionStatement(node *ast.ExStmtNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Ex)
}


