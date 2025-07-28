package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *ast.AssignmentNodeS) (MSVal, error) {
	
	// Evaluate the expression
	res, err := evaluator.evaluateExpression(&node.Exp)

	if err != nil {
		return MSNothing{}, err
	}

	// set the variable in current scope
	name := node.Identifier.Name.Lexeme
	err = evaluator.env.SetVar(name, res)

	if err != nil {
		return MSNothing{}, err
	}

	return res, nil
}

func (evaluator *MSEvaluator) evalVariable(node *ast.VariableExpNodeS) (MSVal, error) {

	// Get the value from the environment
	val, err := evaluator.env.GetVar(node.Name.Lexeme)

	if err != nil {
		return MSNothing{}, err
	}

	return val, nil
	
}

func (evaluator *MSEvaluator) evaluateDeclAssignExpression(node *ast.DeclAssignNodeS) (MSVal, error) {
	
	
	res, err := evaluator.evaluateExpression(&node.Exp)

	if err != nil {
		return MSNothing{}, err
	}

	// set val, even if res contains error.
	name := node.Identifier.Name.Lexeme
	err = evaluator.env.NewVar(name, res)

	if err != nil {
		// Add errors to res errors
		return MSNothing{}, err
	}

	return res, nil
}