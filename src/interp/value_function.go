package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"strings"
)

///////////////////////////////////////////////////////////////
// User defined function
///////////////////////////////////////////////////////////////

type MSFunction struct {
	fbody *ast.BlockNodeS				// How was the function declared
	boundParams []ParamBindingS			// bound params
	unBoundParams []ParamBindingS		// unbound
	returnType mstype.MSType			// return type for uninit functions
	name *ast.VariableExpNodeS			// function name
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
