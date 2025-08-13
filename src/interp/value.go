package interp

import (
	"mikescript/src/mstype"
)


type MSVal interface {

	// ask type for type checking
	Type() mstype.MSType

	// stringer
	String() string

	// Nullable and get null value
	Nullable() bool
	NullVal() MSVal
}
