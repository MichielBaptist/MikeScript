package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

type MSArray struct {
	Values []MSVal
	VType mstype.MSType
}

func (n MSArray) Type() mstype.MSType {
	return &mstype.MSArrayType{Type: n.VType}
}

func (n MSArray) String() string {

	strs := make([]string, len(n.Values))
	for i, v := range n.Values {
		strs[i] = v.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ","))
}

// --------------------------------------------------------
// implmeents indexable
// --------------------------------------------------------

func (a MSArray) Get(at MSVal) (MSVal, error) {

	// early exit on bad index
	if err := a.ValidIndex(at) ; err != nil {
		return MSNothing{}, nil
	}

	// safe cast now
	idxInt := at.(MSInt)

	return a.Values[idxInt.Val], nil
}

func (a MSArray) Set(at, val MSVal) (MSVal, error) {

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


	// Actually set the value
	a.Values[idxInt.Val] = val

	return val, nil
}

func (a MSArray) ValidIndex(idx MSVal) error {

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

func (a MSArray) ValidValue(val MSVal) error {
	if !a.VType.Eq(val.Type()) {
		msg := fmt.Sprintf("Cannot assign '%s' of type '%s', expected type '%s'", val, val.Type(), a.VType)
		return &EvalError{message: msg}
	}
	return nil
}