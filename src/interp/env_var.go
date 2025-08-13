package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

const typecol_size int = 20
const namecol_size int = 20
const defcol_size int = 40

// ////////////////////////////////////////
// // Environment row
// ////////////////////////////////////////

func stringCut(s string, n int) string {
	if (len(s) > n && len(s) > 3){
		return fmt.Sprintf("%s...", s[0:n-3])
	}
	return s
}

////////////////////////////////////////
// Environment & constructor
////////////////////////////////////////

type Environment struct {
	variables map[string]MSVal
	types map[string]mstype.MSType
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		variables: make(map[string]MSVal),
		types: make(map[string]mstype.MSType),
		enclosing: enclosing,
	}
}

////////////////////////////////////////
// Environment methods
////////////////////////////////////////

func (env *Environment) walkBack(depth int) *Environment {
	target := env
	for i := 0 ; i < depth ; i++ {
		target = target.enclosing
	}
	return target
}

func (env *Environment) GetVar(name string, depth int) (MSVal, error) {

	// Walk back depth environments
	targetEnv := env.walkBack(depth)

	// Checks at evaluation time if the variable is defined
	if val, ok := targetEnv.variables[name] ; ok{
		return val, nil
	} else {
		return MSNothing{}, &EnvironmentError{fmt.Sprintf("Variable '%v' is not defined", name)}
	}
}

func (env *Environment) NewVar(name string, value MSVal) error {

	// Check if already in env. We don't allow re-declaring variables
	if val, ok := env.variables[name] ; ok {
		return &EnvironmentError{fmt.Sprintf("Variable '%v' is already defined as '%v'", name, val)}
	}

	// Set variable
	env.variables[name] = value

	return nil
}

func (env *Environment) SetVar(name string, value MSVal, depth int) error {

	targetEnv := env.walkBack(depth)

	// If env does not contain the var, throw error
	if _, ok := targetEnv.variables[name] ; !ok {
		return varNotFound(name)
	}

	// Variable is defined, first check for type compatibility.
	if err := targetEnv.compatibleType(name, value) ; err != nil {
		return err
	}

	// Set the value, this is safe now
	targetEnv.variables[name] = value

	return nil
}

////////////////////////////////////////
// Helper functions
////////////////////////////////////////

func rowRepr(name string, value MSVal) string {
	c1 := stringCut(fmt.Sprintf("%v", value.Type()), typecol_size)
	c2 := stringCut(name, namecol_size)
	c3 := stringCut(value.String(), defcol_size)
	return fmt.Sprintf("| %-*v | %-*v | %-*v |", typecol_size, c1, namecol_size, c2, defcol_size, c3)
}

func (env *Environment) printEnv() int {

	if env == nil {
		return 0
	}

	// print enclosing scope first
	depth := env.enclosing.printEnv()

	rows := []string{}
	for name, value := range env.variables {
		rows = append(rows, rowRepr(name, value))
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


func (env *Environment) compatibleType(name string, newValue MSVal) error {

	// Get relevant row; crash on issue
	oldValue, ok := env.variables[name]

	// ok check, should never trigger normally
	if !ok {
		return varNotFound(name)
	}

	// Crash on dereferencing nil, then there is an issue.
	expectedType := oldValue.Type()
	receivedType := newValue.Type()
	ok = expectedType.Eq(receivedType)

	if ok {
		return nil
	} else {
		return incompatibleTypes(name, expectedType, receivedType)
	}
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

func varNotFound(name string) error {
	return &EnvironmentError{message: fmt.Sprintf("Variable '%s' is not defined", name)}
}

func incompatibleTypes(name string, target, val mstype.MSType) error {
	return &EnvironmentError{fmt.Sprintf("Variable '%v' is of type '%v' and cannot be assigned a value of type '%v'", name, target, val)}
}

///////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////

func (env *Environment) GetType(name string, depth int) (mstype.MSType, error) {

	targetEnv := env.walkBack(depth)

	t, ok := targetEnv.types[name]

	if !ok {
		msg := fmt.Sprintf("Could not resolve type '%s'", name)
		return nil, &EnvironmentError{message: msg}
	}

	return t, nil
}

func (env *Environment) NewType(name string, t mstype.MSType) error {

	// Check if already in env. We don't allow re-declaring variables
	if typ, ok := env.types[name] ; ok {
		return &EnvironmentError{fmt.Sprintf("Type '%v' is already defined as '%v'", name, typ)}
	}

	// Set variable
	env.types[name] = t

	return nil
}
