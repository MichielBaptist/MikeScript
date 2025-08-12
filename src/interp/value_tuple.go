package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

type MSTuple struct {
	Values []MSVal
}

func (t MSTuple) Type() mstype.MSType {
	types := make([]mstype.MSType, len(t.Values))
	for i, arg := range t.Values {
		types[i] = arg.Type()
	}
	return &mstype.MSCompositeTypeS{Types: types}
}

func (i MSTuple) String() string {
	strs := make([]string, len(i.Values))
	for i, v := range i.Values{
		strs[i] = v.String()
	}
	return "(" + strings.Join(strs, ", ") + ")"
}

func (r MSTuple) Nullable() bool {
	return false
}

// --------------------------------------------------------
// implmeents indexable
// --------------------------------------------------------

func (a MSTuple) Get(at MSVal) (MSVal, error) {

	// early exit on bad index
	if err := a.ValidIndex(at) ; err != nil {
		return MSNothing{}, nil
	}

	// safe cast now
	idxInt := at.(MSInt)

	return a.Values[idxInt.Val], nil
}

func (a MSTuple) Set(at, val MSVal) (MSVal, error) {

	// early exit on bad index
	if err := a.ValidIndex(at) ; err != nil {
		return MSNothing{}, nil
	}

	// early exit on bad value
	if err := a.ValidValue(val) ; err != nil {
		return MSNothing{}, err
	}

	// Check if index is MSInt
	idxInt := at.(MSInt)

	// Check if the target type is valid
	targetVal := a.Values[idxInt.Val]

	if !targetVal.Type().Eq(val.Type()) {
		msg := fmt.Sprintf("Cannot assign '%s' or type '%s' at index '%d' of type '%s'", val, val.Type(), idxInt.Val, targetVal.Type())
		return MSNothing{}, &EvalError{message: msg}
	}

	// Actually set the value
	a.Values[idxInt.Val] = val

	return val, nil
}

func (a MSTuple) ValidIndex(idx MSVal) error {

	if idx == nil {
		msg := fmt.Sprintf("Trying to use invalid index '%s'", idx)
		return &EvalError{message: msg}
	}

	idxInt, ok := idx.(MSInt)

	if !ok {
		msg := fmt.Sprintf("Cannot use '%s' of type '%s' as an index, expected type '%s'.", idx, idx.Type(), mstype.MS_INT)
		return &EvalError{message: msg}
	}

		// Check if the index is in range or not
	if idxInt.Val < 0 || idxInt.Val >= len(a.Values) {
		msg := fmt.Sprintf("Array index out of bounds: '%d', expected value in '[%d, %d]'", idxInt.Val, 0, len(a.Values))
		return &EvalError{message: msg}
	}

	return nil
}

func (a MSTuple) ValidValue(val MSVal) error {
	return nil
}