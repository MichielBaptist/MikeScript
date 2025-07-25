package mstype

import (
	"fmt"
	"mikescript/src/token"
	"strings"
)

////////////////////////////////////////////////////////////

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
	case RT_NOTHING:	return "nothing"
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
	Eq(other *MSType) bool
	String() string
}

////////////////////////////////////////////////////////////
type MSSimpleTypeS struct {
	Rt ResultType
}

func (t *MSSimpleTypeS) Eq(o *MSType) bool {
	switch other := (*o).(type){
	case *MSSimpleTypeS: 	return t.Rt == other.Rt
	default:				return false
	}
}

func (t *MSSimpleTypeS) String() string {
	return t.Rt.String()
}


////////////////////////////////////////////////////////////
type MSCompositeTypeS struct {
	Left MSType
	Right MSType
}

func (t *MSCompositeTypeS) Eq(o *MSType) bool {
	switch other := (*o).(type) {
	case *MSCompositeTypeS:	return t.Left.Eq(&other.Left) && t.Right.Eq(&other.Right)
	default:				return false
	}
}

func (t *MSCompositeTypeS) String() string {
	return fmt.Sprintf("(%s, %s)", t.Left, t.Right)
}

////////////////////////////////////////////////////////////
type MSOperationTypeS struct {
	Left []MSType
	Right MSType // result of operation
}

func (t *MSOperationTypeS) Eq(o *MSType) bool {
	switch other := (*o).(type) {
	case *MSOperationTypeS:

		// Compare output of op
		e := t.Right.Eq(&other.Right)

		// Compare parameters of op, note this
		// circuit breaks on a single false
		e = e && len(t.Left) == len(other.Left)
		for i := 0; e && i < len(t.Left) ; i++ {
			e = e && t.Left[i].Eq(&other.Left[i])
		}
		return e
	default:
		return false
	}
}

func (t *MSOperationTypeS) String() string{
	lefts := []string{}
	for _, l := range t.Left {
		lefts = append(lefts, l.String())
	}
	return fmt.Sprintf("%s -> %s", strings.Join(lefts, ", "), t.Right.String())
}

////////////////////////////////////////////////////////////

var MS_BOOL MSType = &MSSimpleTypeS{Rt: RT_BOOL}
var MS_INT MSType = &MSSimpleTypeS{Rt: RT_INT}
var MS_FLOAT MSType = &MSSimpleTypeS{Rt: RT_FLOAT}
var MS_STRING MSType = &MSSimpleTypeS{Rt: RT_STRING}
var MS_NOTHING MSType = &MSSimpleTypeS{Rt: RT_NOTHING}
var MS_BREAK MSType = &MSSimpleTypeS{Rt: RT_BREAK}
var MS_CONTINUE MSType = &MSSimpleTypeS{Rt: RT_CONTINUE}

////////////////////////////////////////////////////////////
type MSTypeError struct {
	Msg string
}

func (e *MSTypeError) Error() string {
	return e.Msg
}