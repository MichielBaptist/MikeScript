package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"strings"
)

///////////////////////////////////////////////////////////////
// User defined function
///////////////////////////////////////////////////////////////

/*
Value that represents a function in Mikescript. The struct needs
to contain all information required to run "run":
	- An AST representing the program the function runs. This is 'nil' when declared using 'var (...->...) f;'
	- A list of unbound parameters. A parameter is: a name, a type and MSVal (nil for unbound)
	- A list of bound parameters. The MSVal of bound parameters is always of the correct type
	- A return type. Either derived from function or var declaration
	- A Name.
	- A closure, the environment when the function was declared:
		- When declared using "function () >> f {...}"
		- When declared using "var (->) f;" (bodyless).

Note: 	using "var (->) f;" allows you to capture a closure for
		an "unknown" or "to be declared" function, which may or may not be
		good idea? You can do some unhinged stuff using this concept, like
		using variables in functions that "come out of nowhere".
*/

type MSFunction struct {
	fbody *ast.BlockNodeS				// How was the function declared
	boundParams []ParamBindingS			// bound params
	unBoundParams []ParamBindingS		// unbound
	returnType mstype.MSType			// return type for uninit functions
	name *ast.VariableExpNodeS			// function name
	closure *Environment				// env at declaration time
}




// --------------------------------------------------------
//	constructors
// --------------------------------------------------------

func NewMSFunction(decl *ast.FuncDeclNodeS, closure *Environment) *MSFunction {

	bindings := make([]ParamBindingS, len(decl.Params))
	for i, p := range decl.Params {
		bindings[i] = paramToBinding(p)
	}

	return &MSFunction{
		fbody: decl.Body,			// function body
		unBoundParams: bindings,	// current bincings
		returnType: decl.Rt,		// declared return type
		name: decl.Fname,			// name
		closure: closure,			// env at declaration
	}

}

func MSFunctionFromType(t *mstype.MSOperationTypeS, name *ast.VariableExpNodeS, closure *Environment) *MSFunction {

	bindings := make([]ParamBindingS, len(t.Left))
	for i, t := range t.Left {
		bindings[i] = typeToBinding(t)
	}

	return &MSFunction{
		name: name,					// from var declaration
		fbody: nil,					// no body since used in "var (... -> ...) f;"
		unBoundParams: bindings,	// 
		returnType: t.Right,
		closure: closure,
	}
}

// --------------------------------------------------------
// Implements MSValue
// --------------------------------------------------------


func (f MSFunction) Type() mstype.MSType {

	// Get list of unbound types
	ptypes := make([]mstype.MSType, len(f.unBoundParams))
	for i, p := range f.unBoundParams {
		ptypes[i] = p.Type
	}

	return &mstype.MSOperationTypeS{
		Left: ptypes,
		Right: f.returnType,
	}
}

func (f MSFunction) String() string {

	ps := []string{}
	for _, bp := range f.boundParams {
		ps = append(ps, "(" + bp.String() + ")")
	}
	for _, up := range f.unBoundParams {
		ps = append(ps, "(" + up.String() + ")")
	}

	// check if body is nil
	var bodys string
	if f.initialized() {
		bodys = "{...}"
	} else {
		bodys = "{}"
	}

	strs := []string{}

	// Join strings
	pss := strings.Join(ps, ", ")
	if pss != "" {
		strs = append(strs, pss)
	}
	strs = append(strs, ">>", f.fname(), "->", f.returnType.String(), bodys)
	
	// bindings >> fname -> rt
	return strings.Join(strs, " ")

}

func (f MSFunction) Initialized() bool {
	return f.fbody != nil
}
