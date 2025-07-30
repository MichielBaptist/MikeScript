package mstype

import (
	"fmt"
	"mikescript/src/token"
)

////////////////////////////////////////////////////////////

// Result type basically defines how to interpret an EvalResult
//
type ResultType uint8
const (
	RT_INVALID ResultType = iota
	RT_NOTHING

	RT_INT
	RT_FLOAT
	RT_STRING
	RT_BOOL

	RT_TUPLE
	RT_FUNCTION

	RT_BREAK
	RT_CONTINUE
	RT_RETURN
)

func (res ResultType) String() string {
	switch res {

	// Real types
	case RT_INT:		return "int"
	case RT_FLOAT:		return "float"
	case RT_STRING:		return "string"
	case RT_BOOL:		return "bool"

	// composite types
	case RT_TUPLE:		return "tuple"
	case RT_FUNCTION:	return "function"

	// statement types (not a real type, but used by interp)
	case RT_INVALID:	return "invalid"
	case RT_NOTHING:	return ""
	case RT_BREAK:		return "break"
	case RT_CONTINUE:	return "continue"
	default:			return "UNKNOWN -- check expression.go"
	}
}

// Interface between token and ast
func TokenToType(tok *token.Token) (MSType, error) {
	switch tok.Type {
	case token.FLOAT_TYPE: 		return &MSSimpleTypeS{Rt: RT_FLOAT}, nil
	case token.INT_TYPE:		return &MSSimpleTypeS{Rt: RT_INT}, nil
	case token.STRING_TYPE:		return &MSSimpleTypeS{Rt: RT_STRING}, nil
	case token.BOOLEAN_TYPE:	return &MSSimpleTypeS{Rt: RT_BOOL}, nil
	default:					return nil, &MSTypeError{Msg: fmt.Sprintf("Cannot convert token: %s to type", tok) }
	}
}

////////////////////////////////////////////////////////////
type MSType interface {
	Eq(other *MSType) bool	// compare types
	String() string			// show type
	Nullable() bool			// Can values of this type be <nil>?
}

// Non nullable types
var MS_BOOL MSType = &MSSimpleTypeS{Rt: RT_BOOL}
var MS_INT MSType = &MSSimpleTypeS{Rt: RT_INT}
var MS_FLOAT MSType = &MSSimpleTypeS{Rt: RT_FLOAT}
var MS_STRING MSType = &MSSimpleTypeS{Rt: RT_STRING}

// The nothing type contains no elements, there is no
// possible value this type can produce.
var MS_NOTHING MSType = &MSSimpleTypeS{Rt: RT_NOTHING}

// break, continue and return types
var MS_BREAK MSType = &MSSimpleTypeS{Rt: RT_BREAK}
var MS_CONTINUE MSType = &MSSimpleTypeS{Rt: RT_CONTINUE}
var MS_RETURN MSType = &MSSimpleTypeS{Rt: RT_RETURN}

////////////////////////////////////////////////////////////

type MSTypeError struct {
	Msg string
}

func (e *MSTypeError) Error() string {
	return e.Msg
}