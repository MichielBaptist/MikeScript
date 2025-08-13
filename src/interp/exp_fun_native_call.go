package interp

import (
	"fmt"
	"mikescript/src/mstype"
)

// -----------------------------------------------------------
// Implements MSCallable
// -----------------------------------------------------------

func (f MSFunction) Call(ev *MSEvaluator) (MSVal, error) {

	if !f.initialized() {
		err := EvalError{fmt.Sprintf("Cannot call uninitialized function '%s'", f.fname())}
		return nil, &err
	}

	// Create a new environment with globals as base scope.
	env := NewEnvironment(f.closure)

	// push all bindings in the env
	for _, bind := range f.boundParams {
		env.NewVar(bind.strName(), bind.Value)
	}

	// Call the body using env
	res, err := ev.executeBlock(f.fbody, env)

	// Check if the block executed properly, if not,
	// we cannot expect an EvalResult with RT_RETURN
	if err != nil {
		return nil, err
	}

	// Check if we can cast to MSReturn
	returnVal := res.(MSReturn)

	// Type check the return value against the
	// declared return type.
	if !returnVal.Type().Eq(f.GetOutputType()) {
		msg := fmt.Sprintf("Tried returning '%s' of type '%s', expected type '%s'", returnVal, returnVal.Type(), f.GetOutputType())
		return nil, &EvalError{msg}
	}

	return returnVal.Val, nil
}

func (f MSFunction) Arity() int {
	return len(f.unBoundParams)
}

func (f MSFunction) Bind(args []MSVal) (MSVal, error) {

	if !f.initialized() {
		return nil, BindingError{msg: fmt.Sprintf("Cannot bind uninitialized function '%s'", f.fname())}
	}

	// Add all the args to the function bindings
	newf, err := f.bindArgs(args)

	// Binding error
	if err != nil {
		return nil, err
	}

	return *newf, nil
}


// -----------------------------------------------------------
// helpers
// -----------------------------------------------------------

func (f *MSFunction) GetOutputType() mstype.MSType{
	return f.returnType
}

func (f *MSFunction) initialized() bool {
	return f.fbody != nil
}

func (f *MSFunction) fname() string {
	return f.name.VarName()
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

	// Copy bindings
	newBound := f.copyBound()
	newUnbound := f.copyUnBound()

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
		closure: f.closure,
	}

	return &fnew, nil
}
