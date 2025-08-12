package interp

// import (
// 	"fmt"
// 	"mikescript/src/mstype"
// )

// ////////////////////////////////////////
// // Error
// ////////////////////////////////////////

// type TypeEnvError struct {
// 	message string
// }

// func (e *TypeEnvError) Error() string {
// 	return "TypeEnv Error: " + e.message
// }

// ////////////////////////////////////////
// // TypeEnv & constructor
// ////////////////////////////////////////

// type TypeEnv struct {
// 	types map[string]mstype.MSType
// }

// func NewTypeEnv() *TypeEnv {
// 	return &TypeEnv{types: make(map[string]mstype.MSType)}
// }

// ////////////////////////////////////////
// // TypeEnv methods
// ////////////////////////////////////////

// func (env *TypeEnv) GetType(name string) (mstype.MSType, error) {

// 	t, ok := env.types[name]

// 	if !ok {
// 		msg := fmt.Sprintf("Could not resolve type '%s'", name)
// 		return nil, &EnvironmentError{message: msg}
// 	}

// 	return t, nil
// }

// func (env *TypeEnv) NewType(name string, t mstype.MSType) error {

// 	// Check if already in env. We don't allow re-declaring variables
// 	if env.containsName(name) {
// 		return &TypeEnvError{fmt.Sprintf("Type '%v' is already defined", name)}
// 	}

// 	// Set variable
// 	env.types[name] = t

// 	return nil
// }

// ////////////////////////////////////////
// // Helper functions
// ////////////////////////////////////////

// func (env *TypeEnv) printEnv() {
// 	for i, s := range env.types {
// 		fmt.Printf("%s --> %v\n", i, s)
// 	}
// }

// func (env *TypeEnv) containsName(name string) bool {
// 	_, ok := env.types[name]
// 	return ok
// }
