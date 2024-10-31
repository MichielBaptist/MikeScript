package main

import "fmt"

type EnvRow struct {
	name 	string			// Name of the variable
	value 	*EvalResult		// Value of the variable
	rtype 	ResultType		// Type of the variable (fixed, this can never change)
}

func (er *EnvRow) String() string {
	return fmt.Sprintf("%v %v: %v", er.rtype, er.name, er.value)
}

type Environment struct {
	// Map of all variables to their values
	variables map[string]EnvRow
}

type EnvironmentError struct {
	message string
}

func NewEnvironment() *Environment {
	return &Environment{make(map[string]EnvRow)}
}

func (e *EnvironmentError) Error() string {
	return "Environment Error: " + e.message
}

func (env *Environment) Get(name string) (EvalResult, error) {

	// Checks at evaluation time if the variable is defined
	if val, ok := env.variables[name]; ok {
		return *val.value, nil
	}
	return EvalResult{}, &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined", name)}
}

func (env *Environment) NewVar(name string, value EvalResult, rtype ResultType) error {

	// Check if the variable is empty
	if name == "" {
		return &EnvironmentError{fmt.Sprintf("Invalid variable name: '%v'", name)}
	}

	// Check if already in env. We don't allow re-declaring variables
	if _, ok := env.variables[name]; ok {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is already defined", name)}
	}

	env.variables[name] = EnvRow{name, &value, rtype}

	return nil
}

func (env *Environment) Set(name string, value EvalResult) error {

	// Check if the variable is empty
	if name == "" {
		return &EnvironmentError{fmt.Sprintf("Invalid variable name: '%v'", name)}
	}

	// Check if the variable is defined
	if _, ok := env.variables[name]; !ok {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined", name)}
	}

	// Check if the given EvalResult type is the same as
	// the one in the environment
	if env.variables[name].rtype != value.ReturnType() {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is of type '%v' and cannot be assigned a value of type '%v'", name, env.variables[name].rtype, value.ReturnType())}
	}

	env.variables[name] = EnvRow{name, &value, env.variables[name].rtype}

	return nil
}

func (env *Environment) printEnv() {
	table := ""
	for k, v := range env.variables {
		table += fmt.Sprintf("%v: %v\n", k, v.String())
	}
	
	fmt.Println(table)
}