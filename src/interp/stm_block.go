package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeBlock(node *ast.BlockNodeS, env *Environment) (MSVal, error) {

	previous := evaluator.env

	evaluator.env = env

	// Pop the environment when we are done with block
	defer func() {
		evaluator.env = previous
	}()

	return evaluator.executeStatements(&ast.Program{Statements: node.Statements})
}