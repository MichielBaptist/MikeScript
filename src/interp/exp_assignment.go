package interp

import "mikescript/src/ast"


func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *ast.AssignmentNodeS) EvalResult {
	
	// Evaluate the expression
	res := evaluator.evaluateExpression(&node.Exp)

	if !res.Valid() {
		return res
	}

	// set the variable in current scope
	err := evaluator.env.SetVar(node.Identifier.Name.Lexeme, res)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return EvalResult{rt: res.rt, val: res.val}
}

func (evaluator *MSEvaluator) evalVariable(node *ast.VariableExpNodeS) EvalResult {

	// Get the value from the environment
	val, err := evaluator.env.GetVar(node.Name.Lexeme)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return val
	
}

func (evaluator *MSEvaluator) evaluateDeclAssignExpression(node *ast.DeclAssignNodeS) EvalResult {
	res := evaluator.evaluateExpression(&node.Exp)

	if !res.Valid() {
		return res
	}

	// set val, even if res contains error.
	name := node.Identifier.Name.Lexeme
	err := evaluator.env.NewVar(name, res, res.rt)

	if err != nil {
		// Add errors to res errors
		return EvalResult{err: append(res.err, err)}
	}

	return res
}