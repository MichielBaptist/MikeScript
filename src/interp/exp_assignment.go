package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *ast.AssignmentNodeS) (MSVal, error) {
	
	// Evaluate the expression
	res, err := evaluator.evaluateExpression(node.Exp)

	if err != nil {
		return MSNothing{}, err
	}

	// set the variable in target scope
	name := node.Identifier.Name.Lexeme

	if depth, ok := evaluator.vlocals[node.Identifier] ; ok {
		err = evaluator.env.SetVar(name, res, depth)
	} else {
		err = evaluator.glb.SetVar(name, res, 0)
	}

	if err != nil {
		return MSNothing{}, err
	}

	return res, nil
}

func (evaluator *MSEvaluator) evalVariable(node *ast.VariableExpNodeS) (MSVal, error) {

	var val MSVal
	var err error

	if depth, ok := evaluator.vlocals[node] ; ok {
		val, err = evaluator.env.GetVar(node.VarName(), depth)
	} else {
		val, err = evaluator.glb.GetVar(node.VarName(), 0)
	}
	

	if err != nil {
		return MSNothing{}, err
	}

	return val, nil
	
}

func (evaluator *MSEvaluator) evaluateDeclAssignExpression(node *ast.DeclAssignNodeS) (MSVal, error) {
	
	
	res, err := evaluator.evaluateExpression(node.Exp)

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