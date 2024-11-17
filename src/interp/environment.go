package interp

import "fmt"

////////////////////////////////////////
// Environment row
////////////////////////////////////////

type EnvRow struct {
	name 	string			// Name of the variable
	value 	*EvalResult		// Value of the variable (always an EvalResult)
	rtype 	ResultType		// Type of the variable (fixed, this can never change)
}

func (er *EnvRow) String() string {
	return fmt.Sprintf("%v %v: %v", er.rtype, er.name, er.value)
}

////////////////////////////////////////
// Error
////////////////////////////////////////

type EnvironmentError struct {
	message string
}

func (e *EnvironmentError) Error() string {
	return "Environment Error: " + e.message
}

////////////////////////////////////////
// Environment & constructor
////////////////////////////////////////

type Environment struct {
	variables map[string]EnvRow
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{make(map[string]EnvRow), enclosing}
}

////////////////////////////////////////
// Environment methods
////////////////////////////////////////

func (env *Environment) GetVar(name string) (EvalResult, error) {

	// Checks at evaluation time if the variable is defined
	if env.containsVar(name) {
		return *(env.variables[name].value), nil
	}

	// Not found in this scope, check enclosing scope
	if env.enclosing != nil {
		return env.enclosing.GetVar(name)
	}

	// No enclosing scope, return error
	return EvalResult{}, &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined", name)}
}

func (env *Environment) NewVar(name string, value EvalResult, rtype ResultType) error {

	// Check if the variable is empty
	if !env.validVarName(name) {
		return &EnvironmentError{fmt.Sprintf("Invalid variable name: '%v'", name)}
	}

	// Check if already in env. We don't allow re-declaring variables
	if env.containsVar(name) {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is already defined", name)}
	}

	env.variables[name] = EnvRow{name, &value, rtype}

	return nil
}

func (env *Environment) SetVar(name string, value EvalResult) error {

	// Check if the variable is empty
	if !env.validVarName(name) {
		return &EnvironmentError{fmt.Sprintf("Invalid variable name: '%v'", name)}
	}

	// Check if the variable is defined. If it's not
	// check the enclosing scope for the variable.
	if !env.containsVar(name) {

		// Don't have the variable in this scope, but
		// has enclosing scope, set the value there.
		if env.enclosing != nil {
			return env.enclosing.SetVar(name, value)
		}

		// No enclosing scope, and no variable found
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined in scope", name)}
	}

	// Variable is defined, first check for type compatibility.
	if !env.compatibleType(name, value) {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is of type '%v' and cannot be assigned a value of type '%v'", name, env.variables[name].rtype, value.ReturnType())}
	}

	// Set the value, this is safe now
	env.variables[name] = EnvRow{name, &value, env.variables[name].rtype}

	return nil
}

////////////////////////////////////////
// Helper functions
////////////////////////////////////////

func (env *Environment) printEnv() {

	if env == nil {
		return
	}

	// print enclosing scope first
	env.enclosing.printEnv()

	fmt.Println("-----------------------------------")

	table := ""
	for k, v := range env.variables {
		table += fmt.Sprintf("%v: %v\n", k, v.String())
	}
	
	fmt.Println(table)
}

func (env *Environment) validVarName(name string) bool {
	return name != ""
}

func (env *Environment) containsVar(name string) bool {
	_, ok := env.variables[name]
	return ok
}

func (env *Environment) compatibleType(name string, value EvalResult) bool {
	return env.variables[name].rtype == value.ReturnType()
}