package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

const typecol_size int = 20
const namecol_size int = 20
const defcol_size int = 40

////////////////////////////////////////
// Environment row
////////////////////////////////////////

func stringCut(s string, n int) string {
	if (len(s) > n && len(s) > 3){
		return fmt.Sprintf("%s...", s[0:n-3])
	}
	return s
}

type EnvRow struct {
	name 	string				// Name of the variable
	value 	*EvalResult			// Value of the variable (always an EvalResult)
	rtype 	mstype.MSType		// Type of the variable (fixed, this can never change)
}

func (er *EnvRow) String() string {
	if er == nil {
		return "ERROR"
	}
	return fmt.Sprintf("%v %v = %v", er.rtype, er.name, er.value)
}

func (er *EnvRow) rowRepr() string {
	c1 := stringCut(er.rtype.String(), typecol_size)
	c2 := stringCut(er.name, namecol_size)
	c3 := stringCut(er.value.String(), defcol_size)
	return fmt.Sprintf("| %-*v | %-*v | %-*v |", typecol_size, c1, namecol_size, c2, defcol_size, c3)
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
	return &Environment{
		variables: make(map[string]EnvRow),
		enclosing: enclosing,
	}
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

func (env *Environment) NewVar(name string, value EvalResult, rtype mstype.MSType) error {

	// Check if the variable is empty
	if !env.validVarName(name) {
		return &EnvironmentError{fmt.Sprintf("Invalid variable name: '%v'", name)}
	}

	if err := env.validValue(name, value) ; err != nil {
		return err
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

	if err := env.validValue(name, value) ; err != nil {
		return err
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
		envVarT := env.variables[name].rtype
		valT := value.Rt
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is of type '%v' and cannot be assigned a value of type '%v'", name, envVarT, valT)}
	}

	// Set the value, this is safe now
	env.variables[name] = EnvRow{name, &value, env.variables[name].rtype}

	return nil
}

func (env *Environment) validValue(name string, val EvalResult) error {

	// Check if we're trying set a value not valid
	if !val.Valid() {
		return &EnvironmentError{fmt.Sprintf("Trying to set '%s' with a value containing an error: %s", name, val)}
	}



	// Check if we're binding a EvalResult containing 'nil' value, this should not happen.
	// if val.val == nil {
	// 	return &EnvironmentError{fmt.Sprintf("Trying to set '%s' with a 'nil' value: %s", name, val)}
	// }

	return nil
}

////////////////////////////////////////
// Helper functions
////////////////////////////////////////

func (env *Environment) printEnv() int {

	if env == nil {
		return 0
	}

	// print enclosing scope first
	depth := env.enclosing.printEnv()

	rows := []string{}
	for _, v := range env.variables {
		rows = append(rows, v.rowRepr())
	}

	// Print table
	fmt.Println(tblbar())
	fmt.Println(strings.Join(rows, "\n"))
	if depth == 0 {
		fmt.Println(tblbar())
	}
	return depth + 1
}

func (env *Environment) validVarName(name string) bool {
	return name != ""
}

func (env *Environment) containsVar(name string) bool {
	_, ok := env.variables[name]
	return ok
}

func (env *Environment) compatibleType(name string, value EvalResult) bool {
	envVarT := env.variables[name].rtype
	valueT := value.Rt
	return envVarT.Eq(&valueT)
}

func tblbar() string {
	return strings.Join([]string{
		"+-",
		strings.Repeat("-", typecol_size),
		"-+-",
		strings.Repeat("-", namecol_size),
		"-+-",
		strings.Repeat("-", defcol_size),
		"-+",
	}, "")
}