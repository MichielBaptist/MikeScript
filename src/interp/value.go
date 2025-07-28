package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)


type MSVal interface {
	Type() mstype.MSType
	String() string
}

////////////////////////////////////////////////////////////
// int
////////////////////////////////////////////////////////////

type MSInt struct {
	Val int
}

func (i MSInt) Type() mstype.MSType {
	return mstype.MS_INT
}

func (i MSInt) String() string {
	return fmt.Sprintf("%v", i.Val)
}

////////////////////////////////////////////////////////////
// float
////////////////////////////////////////////////////////////

type MSFloat struct {
	Val float64
}

func (i MSFloat) Type() mstype.MSType {
	return mstype.MS_FLOAT
}

func (i MSFloat) String() string {
	return fmt.Sprintf("%v", i.Val)
}

////////////////////////////////////////////////////////////
// string
////////////////////////////////////////////////////////////

type MSString struct {
	Val string
}

func (i MSString) Type() mstype.MSType {
	return mstype.MS_STRING
}

func (i MSString) String() string {
	return i.Val
}

////////////////////////////////////////////////////////////
// bool
////////////////////////////////////////////////////////////

type MSBool struct {
	Val bool
}

func (b MSBool) Type() mstype.MSType {
	return mstype.MS_BOOL
}

func (i MSBool) String() string {
	return fmt.Sprintf("%v", i.Val)
}

////////////////////////////////////////////////////////////
// tuple
////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////
// return
////////////////////////////////////////////////////////////

type MSReturn struct {
	Val MSVal
}

func (r MSReturn) Type() mstype.MSType {
	return r.Val.Type()
}

func (r MSReturn) String() string {
	return fmt.Sprintf("Return[%s]", r.Val.String())
}

////////////////////////////////////////////////////////////
// break
////////////////////////////////////////////////////////////

type MSBreak struct {}

func (r MSBreak) Type() mstype.MSType {
	return mstype.MS_BREAK
}

func (r MSBreak) String() string {
	return "break"
}


////////////////////////////////////////////////////////////
// continue
////////////////////////////////////////////////////////////

type MSContinue struct {}

func (r MSContinue) Type() mstype.MSType {
	return mstype.MS_BREAK
}

func (r MSContinue) String() string {
	return "continue"
}

////////////////////////////////////////////////////////////
// Nothing
////////////////////////////////////////////////////////////

type MSNothing struct {}
func (n MSNothing) Type() mstype.MSType {
	return mstype.MS_NOTHING
}
func (n MSNothing) String() string {
	return "nothing"
}