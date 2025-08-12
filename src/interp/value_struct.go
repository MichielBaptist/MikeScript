package interp

import (
	"fmt"
	"mikescript/src/mstype"
)

type MSStruct struct {
	Name string					// struct name 
	SType mstype.MSType			// reference to type (named type normally)
	Fields map[string]MSVal		// mapping
}

func (i MSStruct) Type() mstype.MSType {
	return i.SType
}

func (i MSStruct) String() string {
	if i.IsNil() {
		return "nothing"
	}
	return fmt.Sprintf("%v{%v}", i.Name, i.Fields)
}

func (r MSStruct) Nullable() bool {
	return true
}

func (r MSStruct) IsNil() bool {
	return r.Fields == nil
}

// ----------------------------------------------------------------
// implements MSFieldable
// ----------------------------------------------------------------

func (s MSStruct) Get(field string) (MSVal, error) {

	// Cannot access a 'nothing' struct.
	// This struct serves as a placeholder for nothing
	// But still contains a reference to a type so
	// We can do proper type checking.
	if s.IsNil() {
		msg := fmt.Sprintf("Cannot access field '%s' of 'nothing' of type 'nothing'", field)
		return nil, &EvalError{message: msg}
	}

	if val, ok := s.Fields[field]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("field %q not found", field)
}

func (s MSStruct) Set(field string, val MSVal) (MSVal, error) {

	if s.IsNil() {
		msg := fmt.Sprintf("Cannot set field '%s' of 'nothing' of type 'nothing'", field)
		return nil, &EvalError{message: msg}
	}

	// checks same type
	if err := s.ValidField(field); err != nil {
		return nil, err
	}

	// type check aswell
	if err := s.ValidValue(field, val); err != nil {
		return nil, err
	}

	// set
	s.Fields[field] = val

	return val, nil
}

func (i MSStruct) ValidField(field string) error {
	if _, ok := i.Fields[field]; !ok {
		return fmt.Errorf("field %q is not valid", field)
	}
	return nil
}

func (i MSStruct) ValidValue(field string, val MSVal) error {
	if err := i.ValidField(field); err != nil {
		return err
	}
	if !val.Type().Eq(i.Fields[field].Type()) {
		return fmt.Errorf("field %q expects type %q, got %q",
			field, i.Fields[field].Type(), val.Type())
	}
	return nil
}