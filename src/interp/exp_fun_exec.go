package interp

import (
	"fmt"
	"mikescript/src/ast"
	"strings"
)


type FunctionResult interface {
	call(e *MSEvaluator, args []EvalResult) EvalResult	// native or foreign
	arity() int											// # args left
}

///////////////////////////////////////////////////////////////
// helper struct
///////////////////////////////////////////////////////////////

type ParamBindingS struct {
	Type ResultType				// expected type of param
	Name ast.VariableExpNodeS 	// Name
	Value *EvalResult			// Can be nil when unbound
}

func paramToBinding(p ast.FuncParamS) ParamBindingS {
	// Converts a Function parameter to a binding
	// Initializes all bindings to nil

	return ParamBindingS{
		Type: declaredTypeToReturnType(p.Type),
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

func (b *ParamBindingS) bind(val EvalResult) (ParamBindingS, error) {

	// Check if param is already bound, we currently don't allow
	// re-binding of parameters; though it could be interesting???
	if b.Value != nil {
		msg := fmt.Sprintf("Cannot re-bind parameter '%s'", b.String())
		return *b, BindingError{msg: msg}
	}

	// Validate correct types
	if b.Type != val.rt {
		msg := fmt.Sprintf("Cannot bind '%s' of type '%s' to parameter '%s' of type '%s'", val.String(), val.rt.String(), b.Name.String(), b.Type.String())
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
		vals = b.Value.String()
	}

	return fmt.Sprintf("%s %s = %s", b.Type.String(), b.Name.String(), vals)
}

// -----------------------------------------------------------
// BindingError
// -----------------------------------------------------------


type BindingError struct {
	msg string
}

func (e BindingError) Error() string {
	return e.msg
}


///////////////////////////////////////////////////////////////
// User defined function
///////////////////////////////////////////////////////////////

func NewMSFunction(decl *ast.FuncDeclNodeS) *MSFunction {

	// Create the bindings for the function
	bindings := make([]ParamBindingS, len(decl.Params))
	for i, p := range decl.Params {
		bindings[i] = paramToBinding(p)
	}

	return &MSFunction{decl: decl, unBoundParams: bindings}

}

type MSFunction struct {
	decl *ast.FuncDeclNodeS 			// How was the function declared
	boundParams []ParamBindingS			// bound params
	unBoundParams []ParamBindingS		// unbound
}

// -----------------------------------------------------------
// Implements FunctionResult
// -----------------------------------------------------------

func (f *MSFunction) call(ev *MSEvaluator, args []EvalResult) EvalResult {

	if !f.initialized() {
		return evalErr(fmt.Sprintf("Cannot call uninitialized function '%s'", f.fname()))
	}

	// Add all the args to the function bindings
	newf, err := f.bindArgs(args)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	if newf.execable() {
		// parameters full
		return newf.exec(ev)
	} else {
		// parameters not full
		return EvalResult{rt: RT_FUNCTION, val: newf}
	}
}

func (f *MSFunction) exec(evaluator *MSEvaluator) EvalResult {

	// Create a new environment with globals as base scope.
	env := NewEnvironment(evaluator.glb)

	// push all bindings in the env
	for _, bind := range f.boundParams {
		env.NewVar(bind.strName(), *bind.Value, bind.Type)
	}

	// Call the body using env
	return evaluator.executeBlock(f.body(), env)

}

func (f *MSFunction) arity() int {
	return len(f.unBoundParams)
}

// -----------------------------------------------------------
// Implements Stringer
// -----------------------------------------------------------

func (f *MSFunction) String() string {

	ps := []string{}
	for _, bp := range f.boundParams {
		ps = append(ps, bp.String())
	}
	for _, up := range f.unBoundParams {
		ps = append(ps, up.String())
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
	strs = append(strs, ">>", f.fname(), "->", f.decl.Rt.Type.String(), bodys)
	
	// bindings >> fname -> rt
	return strings.Join(strs, " ")

}

// -----------------------------------------------------------
// helpers
// -----------------------------------------------------------

func (f *MSFunction) initialized() bool {
	return f.body() != nil
}

func (f *MSFunction) body() *ast.BlockNodeS {
	return f.decl.Body
}

func (f *MSFunction) execable() bool {
	return f.arity() == 0
}

func (f *MSFunction) fname() string {
	// Name of function from declaration
	return f.decl.Fname.String()
}


func (f *MSFunction) checkArity(args []EvalResult) error {

	// info
	nargs := len(args)
	npars := f.arity()

	// Get amount of parameters (definition)
	if npars < nargs{
		// TODO: allow returning as tuple when exceeding arity????
		// Could be the move to be honest, but it MIGHT make the 
		// language unusable, donno...
		msg := fmt.Sprintf("Exceeded function arity. Expected %v arguments, but received %v (already bound: %v)", nargs, npars, len(f.boundParams))
		return BindingError{msg: msg}
	}

	return nil
}

func (f *MSFunction) copyBound() []ParamBindingS {
	new := make([]ParamBindingS, len(f.boundParams))
	for i, bp := range f.boundParams {
		new[i] = bp.copy()
	} 
	return new
}

func (f *MSFunction) copyUnBound() []ParamBindingS {
	new := make([]ParamBindingS, len(f.unBoundParams))
	for i, bp := range f.unBoundParams {
		new[i] = bp.copy()
	} 
	return new
}


func (f *MSFunction) bindArgs(args []EvalResult) (*MSFunction, error) {

	// check arity, if we don't have enough space
	// in unbounded parameters, we can't bind all args.
	if err := f.checkArity(args) ; err != nil {
		return nil, err
	}

	// Copy bindings. Note we can't change the values of the
	// binding directly as they are EvalResult pointers
	newBound := f.copyBound()
	newUnbound := f.copyUnBound()

	// Loop over args and compare the return type to the
	// expected type.
	for i, arg := range args {

		// Get unbound param
		up := newUnbound[i]

		// bind
		up, err := up.bind(arg)

		// Bind can fail on type fault
		if err != nil {
			return nil, err
		}

		// Add to bound params
		newBound = append(newBound, up)
	}

	// Remove the first len(args) from the unbound list
	newUnbound = newUnbound[len(args):]

	// Creates a new MSFunction struct containing the new bindings
	fnew := MSFunction{
		decl: f.decl,
		boundParams: newBound,
		unBoundParams: newUnbound,
	}

	return &fnew, nil
}

