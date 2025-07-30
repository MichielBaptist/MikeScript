package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

type ParamBindingS struct {
	Type mstype.MSType				// expected type of param
	Name ast.VariableExpNodeS 		// Name
	Value *MSVal					// Can be nil when unbound
}


func typeToBinding(t mstype.MSType) ParamBindingS {
	return ParamBindingS{
		Type: t,
		Value: nil,
	}
}


func paramToBinding(p ast.FuncParamS) ParamBindingS {
	// Converts a Function parameter to a binding
	// Initializes all bindings to nil

	return ParamBindingS{
		Type: p.Type,
		Name: p.Iden,
		Value: nil,
	}
}

func (b *ParamBindingS) copy() ParamBindingS {
	// Copies a binding struct, done to adjust bound values.

	return ParamBindingS{
		Type: b.Type,
		Name: b.Name,
		Value: b.Value,
	}
}

func (b *ParamBindingS) bind(val MSVal) (ParamBindingS, error) {

	// Check if param is already bound, we currently don't allow
	// re-binding of parameters; though it could be interesting???
	if b.Value != nil {
		msg := fmt.Sprintf("Cannot re-bind parameter '%s'", b.String())
		return *b, BindingError{msg: msg}
	}

	// Validate correct types
	if !b.ValidBindingEvalResult(&val) {
		vals := val.String()
		typs := val.Type().String()
		pname := b.Name.VarName()
		ptype := b.Type.String()
		msg := fmt.Sprintf("Cannot bind '%s' of type '%s' to parameter '%s' of type '%s'", vals, typs, pname, ptype)
		return *b, BindingError{msg: msg}
	}
	
	// Bind value
	b.Value = &val

	// done
	return *b, nil
}

func (b *ParamBindingS) strName() string {
	return b.Name.Name.Lexeme
}

func (b *ParamBindingS) String() string {
	var vals string
	if b.Value == nil {
		vals = "_"
	} else {
		vals = (*b.Value).String()
	}

	return fmt.Sprintf("%s %s = %s", b.Type.String(), b.Name.VarName(), vals)
}

func (b *ParamBindingS) ValidBindingEvalResult(t *MSVal) bool {
	expectedType := b.Type
	receivedType := (*t).Type()
	return expectedType.Eq(&receivedType)
}

// -----------------------------------------------------------
// BindingError
// -----------------------------------------------------------


type BindingError struct {
	msg string
}

func (e BindingError) Error() string {
	return "Binding error:" + e.msg
}