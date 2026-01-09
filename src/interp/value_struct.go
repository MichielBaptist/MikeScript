package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
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

	fieldss := make([]string, 0, len(i.Fields))
	for name, val := range i.Fields {
		fieldss = append(fieldss, fmt.Sprintf("%s: %s", name, val.String()))
	}
	// join
	fieldsStr := strings.Join(fieldss, ", ")
	return fmt.Sprintf("%v{%v}", i.Name, fieldsStr)
}

func (r MSStruct) Nullable() bool {
	return true
}

func (s MSStruct) NullVal() MSVal {
	return MSStruct{Name: s.Name, Fields: nil, SType: s.SType}
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

	if err := s.ValidField(field); err != nil {
		return nil, err
	}

	if err := s.ValidValue(field, val); err != nil {
		return nil, err
	}

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

	// Check if the field is nullable
	// if the value is MSNothing, we allow it
	_, isNothing := val.(MSNothing)
	
	if i.Fields[field].Nullable() && isNothing {
		return nil
	}

	if !i.Fields[field].Nullable() && isNothing {
		return fmt.Errorf("field %q of type %q is not nullable.", field, i.Fields[field].Type())
	}

	if !val.Type().Eq(i.Fields[field].Type()) {
		return fmt.Errorf("field %q expects type %q, got %q",
			field, i.Fields[field].Type(), val.Type())
	}
	return nil
}

