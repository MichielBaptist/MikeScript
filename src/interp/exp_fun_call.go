package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)



func NewMSFunction(decl *ast.FuncDeclNodeS) *MSFunction {

	// Create the bindings for the function
	bindings := make([]ParamBindingS, len(decl.Params))
	for i, p := range decl.Params {
		bindings[i] = paramToBinding(p)
	}

	return &MSFunction{
		fbody: decl.Body,
		unBoundParams: bindings,
		returnType: decl.Rt,
		name: &decl.Fname,
	}

}

func MSFunctionFromType(t *mstype.MSOperationTypeS, name *ast.VariableExpNodeS) *MSFunction {

	// Create the bindings for the function
	bindings := make([]ParamBindingS, len(t.Left))
	for i, t := range t.Left {
		bindings[i] = typeToBinding(t)
	}

	return &MSFunction{
		name: name,					// from var declaration
		fbody: nil,					// no body since used in "var (... -> ...) f;"
		unBoundParams: bindings,
		returnType: t.Right,
	}
}

// -----------------------------------------------------------
// Implements FunctionResult
// -----------------------------------------------------------

func (f MSFunction) Call(ev *MSEvaluator) (MSVal, error) {

	if !f.initialized() {
		err := EvalError{fmt.Sprintf("Cannot call uninitialized function '%s'", f.fname())}
		return MSNothing{}, &err
	}

	// Create a new environment with globals as base scope.
	env := NewEnvironment(ev.glb)

	// push all bindings in the env
	for _, bind := range f.boundParams {
		env.NewVar(bind.strName(), *bind.Value)
	}

	// Call the body using env
	res, err := ev.executeBlock(f.fbody, env)

	// Check if the block executed properly, if not,
	// we cannot expect an EvalResult with RT_RETURN
	if err != nil {
		return MSNothing{}, err
	}

	// Check if we can cast to MSReturn
	returnVal := res.(MSReturn)

	// Type check the return value against the
	// declared return type.
	if !returnVal.Type().Eq(f.GetOutputType()) {
		msg := fmt.Sprintf("Tried returning '%s' of type '%s', expected type '%s'", returnVal, returnVal.Type(), *f.GetOutputType())
		return MSNothing{}, &EvalError{msg}
	}

	return returnVal.Val, nil
}

func (f MSFunction) Arity() int {
	return len(f.unBoundParams)
}

func (f MSFunction) Bind(args []MSVal) (MSVal, error) {

	if !f.initialized() {
		return MSNothing{}, BindingError{msg: fmt.Sprintf("Cannot bind uninitialized function '%s'", f.fname())}
	}

	// Add all the args to the function bindings
	newf, err := f.bindArgs(args)

	// Binding error
	if err != nil {
		return MSNothing{}, err
	}

	return *newf, nil
}


// -----------------------------------------------------------
// helpers
// -----------------------------------------------------------

func (f *MSFunction) GetOutputType() *mstype.MSType{
	return &f.returnType
}

func (f *MSFunction) initialized() bool {
	return f.fbody != nil
}

func (f *MSFunction) fname() string {
	return f.name.String()
}


func (f *MSFunction) checkArity(args []MSVal) error {

	// info
	nargs := len(args)
	npars := f.Arity()

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


func (f *MSFunction) bindArgs(args []MSVal) (*MSFunction, error) {

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
		fbody: f.fbody,
		boundParams: newBound,
		unBoundParams: newUnbound,
		returnType: f.returnType,
		name: f.name,
	}

	return &fnew, nil
}

