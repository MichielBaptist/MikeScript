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
	name 	string			// Name of the variable
	value 	*MSVal	// Value of the variable (always an EvalResult)
}

func (er *EnvRow) String() string {
	if er == nil {
		return "ERROR"
	}
	return fmt.Sprintf("%v %v = %v", (*er.value).Type().String(), er.name, er.value)
}

func (er *EnvRow) rowRepr() string {
	val := *er.value

	c1 := stringCut(val.Type().String(), typecol_size)
	c2 := stringCut(er.name, namecol_size)
	c3 := stringCut(val.String(), defcol_size)
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

func (env *Environment) GetVar(name string) (MSVal, error) {

	// Checks at evaluation time if the variable is defined
	if env.containsVar(name) {
		return *env.variables[name].value, nil
	}

	// Not found in this scope, check enclosing scope
	if env.enclosing != nil {
		return env.enclosing.GetVar(name)
	}

	// No enclosing scope, return error
	return MSNothing{}, &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined", name)}
}

func (env *Environment) NewVar(name string, value MSVal) error {

	// Check if already in env. We don't allow re-declaring variables
	if env.containsVar(name) {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is already defined", name)}
	}

	// Set variable
	env.variables[name] = EnvRow{name, &value}

	return nil
}

func (env *Environment) SetVar(name string, value MSVal) error {

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
		expectedType := env.varType(name)
		receivedType := value.Type()
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is of type '%v' and cannot be assigned a value of type '%v'", name, expectedType, receivedType)}
	}

	// Set the value, this is safe now
	env.variables[name] = EnvRow{name, &value}

	return nil
}

////////////////////////////////////////
// Helper functions
////////////////////////////////////////

func (env *Environment) varType(name string) mstype.MSType {
	row := env.variables[name].value
	typ := (*row).Type()
	return typ
}

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
	if depth == 0 {
		fmt.Println(tblbar(depth))
	}
	if len(rows) > 0 {
		fmt.Println(strings.Join(rows, "\n"))
	}
	fmt.Println(tblbar(depth + 1))
	
	return depth + 1
}

func (env *Environment) depth() int {
	if env == nil {
		return 0
	}
	return 1 + env.enclosing.depth()
}

func (env *Environment) containsVar(name string) bool {
	_, ok := env.variables[name]
	return ok
}

func (env *Environment) compatibleType(name string, value MSVal) bool {
	row := env.variables[name]
	expectedType := (*row.value).Type()
	receivedType := value.Type()
	return expectedType.Eq(&receivedType)
}

func tblbar(depth int) string {

	middle := (namecol_size - 2) / 2

	return strings.Join([]string{
		"+-",
		strings.Repeat("-", typecol_size),
		"-+-",
		strings.Repeat("-", middle),
		fmt.Sprintf("%02d", depth),
		strings.Repeat("-", middle),
		"-+-",
		strings.Repeat("-", defcol_size),
		"-+",
	}, "")
}