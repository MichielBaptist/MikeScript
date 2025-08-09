package interp

import (
	"fmt"
	"mikescript/src/mstype"
)

type MSStruct struct {
	Name string
	Fields map[string]MSVal
}

func (i MSStruct) Type() mstype.MSType {

	types := make(map[string]mstype.MSType)
	for name, val := range i.Fields {
		types[name] = val.Type()
	}

	return &mstype.MSStructTypeS{Name: i.Name, Fields: types}
}

func (i MSStruct) String() string {
	return fmt.Sprintf("%v{%v}", i.Name, i.Fields)
}

// ----------------------------------------------------------------
// implements MSFieldable
// ----------------------------------------------------------------

func (s MSStruct) Get(field string) (MSVal, error) {
	if val, ok := s.Fields[field]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("field %q not found", field)
}

func (s MSStruct) Set(field string, val MSVal) (MSVal, error) {

	// checks same type
	if err := s.ValidField(field); err != nil {
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