package interp

import (
	"fmt"
	"mikescript/src/mstype"
)

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

func (i MSInt) Nullable() bool {
	return false
}

func (i MSInt) NullVal() MSVal {
	return nil
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

func (i MSFloat) Nullable() bool {
	return false
}

func (i MSFloat) NullVal() MSVal {
	return nil
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

func (i MSString) Nullable() bool {
	return false
}

func (i MSString) NullVal() MSVal {
	return nil
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

func (i MSBool) Nullable() bool {
	return false
}

func (i MSBool) NullVal() MSVal {
	return nil
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

func (r MSReturn) Nullable() bool {
	return false
}

func (i MSReturn) NullVal() MSVal {
	return nil
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

func (r MSBreak) Nullable() bool {
	return false
}

func (i MSBreak) NullVal() MSVal {
	return nil
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

func (r MSContinue) Nullable() bool {
	return false
}

func (i MSContinue) NullVal() MSVal {
	return nil
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

func (r MSNothing) Nullable() bool {
	return true
}

func (i MSNothing) NullVal() MSVal {
	return i
}