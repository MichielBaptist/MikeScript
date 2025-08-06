package interp

import (
	"fmt"
	"mikescript/src/ast"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type MSEvaluator struct {
	ast *ast.Program				// The AST to evaluate
	err []error						// Evaluation errors
	env *Environment				// Environment of curresnt scope
	glb *Environment 				// Fixed reference to global scope (outermost env)
	types *TypeEnv					// Type definitions
	locals map[ast.ExpNodeI]int		// How deep do we need to go to resolve variables?
}

func NewMSEvaluator() *MSEvaluator {

	// Init gobal and local scope
	env := NewEnvironment(nil)
	glb := env

	// Add builtins to glb
	glb.NewVar("print", MSBuiltinPrint())
	glb.NewVar("print_env", MSBuiltinPrintEnv())

	// Create type env
	types := NewTypeEnv()

	return &MSEvaluator{
		env: env,
		glb: glb,
		locals: make(map[ast.ExpNodeI]int),
		types: types,
	}
}

func (evaluator *MSEvaluator) UpdateLocals(locals map[ast.ExpNodeI]int) {
	// add all locals to current locals
	for k, v := range locals {
		evaluator.locals[k] = v
	}
}

func (evaluator *MSEvaluator) Eval(ast *ast.Program) (MSVal, error) {

	// set the ast
	evaluator.ast = ast

	// evaluate the ast
	return evaluator.executeStatements(evaluator.ast)
	
}

func (evaluator *MSEvaluator) Errors() []error {
	return evaluator.err
}

func (evaluator *MSEvaluator) PrintEnv() {
	evaluator.env.printEnv()
}

func (e *MSEvaluator) PrintTypes() {
	fmt.Println("-------------------------------------------------------")
	e.types.printEnv()
	fmt.Println("-------------------------------------------------------")
}