package interp

import "mikescript/src/ast"


func (evaluator *MSEvaluator) executeBlock(node *ast.BlockNodeS) EvalResult {

	// Create a new environment with the current 
	// environment as the enclosing environment
	evaluator.env = NewEnvironment(evaluator.env)

	// Pop the environment when we are done with block
	defer func() {
		evaluator.env = evaluator.env.enclosing
	}()

	return evaluator.executeStatements(&ast.Program{Statements: node.Statements})
}