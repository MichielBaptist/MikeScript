package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *ast.AssignmentNodeS) (MSVal, error) {
	
	res, err := evaluator.evaluateExpression(node.Exp)

	if err != nil {
		return nil, err
	}

	// If res is 'nothing', we need to assign the
	// nothing value if the target is nullable
	currentVal, err := evaluator.evalVariable(node.Identifier)

	if err != nil {
		return nil, err
	}

	if _, ok := res.(MSNothing) ; ok && currentVal.Nullable() {
		res = currentVal.NullVal()
	}

	// set the variable in target scope
	depth, ok := evaluator.vlocals[node.Identifier]
	name := node.Identifier.VarName()
	if ok {
		err = evaluator.env.SetVar(name, res, depth)
	} else {
		err = evaluator.glb.SetVar(name, res, 0)
	}

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return val, nil
	
}

func (evaluator *MSEvaluator) evaluateDeclAssignExpression(node *ast.DeclAssignNodeS) (MSVal, error) {
	
	
	res, err := evaluator.evaluateExpression(node.Exp)

	if err != nil {
		return nil, err
	}

	name := node.Identifier.Name.Lexeme
	err = evaluator.env.NewVar(name, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}