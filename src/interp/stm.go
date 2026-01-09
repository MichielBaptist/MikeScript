package interp

import (
	"fmt"
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeStatement(node ast.StmtNodeI) (MSVal, error) {
	switch node := node.(type) {
	case *ast.Program:					return evaluator.executeStatements(node)
	case *ast.VarDeclNodeS:				return evaluator.executeDeclarationStatement(node)
	case *ast.FuncDeclNodeS:			return evaluator.executeFuncDeclStatement(node)
	case *ast.BlockNodeS:				return evaluator.executeBlock(node, NewEnvironment(evaluator.env))
	case *ast.IfNodeS:					return evaluator.executeIfstatement(node)
	case *ast.WhileNodeS:				return evaluator.executeWhileStatement(node)
	case *ast.ContinueNodeS:			return MSContinue{}, nil
	case *ast.BreakNodeS:				return MSBreak{}, nil
	case *ast.ExStmtNodeS:				return evaluator.executeExpressionStatement(node)
	case *ast.ReturnNodeS: 				return evaluator.executeReturnStatement(node)
	case *ast.TypeDefStatementS:		return evaluator.executeTypeDeclaration(node)
	case *ast.StructDeclarationNodeS:	return evaluator.executeStructDeclaration(node)
	case *ast.ForNodeS:					return evaluator.executeForStatement(node)
	default:							return MSNothing{}, &EvalError{fmt.Sprintf("Unknown statement type: %v", node)}
	}
}

func (evaluator *MSEvaluator) executeExpressionStatement(node *ast.ExStmtNodeS) (MSVal, error) {
	return evaluator.evaluateExpression(node.Ex)
}


