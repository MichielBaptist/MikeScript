package ast

import (
	"mikescript/src/mstype"
	"mikescript/src/token"
)

type Program struct {
	Statements []StmtNodeI
}

type BlockNodeS struct {
	Statements []StmtNodeI
}

type VarDeclNodeS struct {
	Identifier 	VariableExpNodeS	// Name
	Vartype 	mstype.MSType		// Type 
}

type ExStmtNodeS struct {
	Ex ExpNodeI
}

type IfNodeS struct {
	Condition 	ExpNodeI
	ThenStmt 	StmtNodeI
	ElseStmt 	StmtNodeI
}

type WhileNodeS struct {
	Condition 	ExpNodeI
	Body 		BlockNodeS
}

type ContinueNodeS struct {
	Tk token.Token
}

type BreakNodeS struct {
	Tk token.Token
}

type ReturnNodeS struct {
	Node ExpNodeI
}

type FuncDeclNodeS struct {
	Fname VariableExpNodeS				// Name
	Params []FuncParamS 				// Parameters
	Rt mstype.MSType					// Return type
	Body *BlockNodeS					// Body of function, may be nil
}

type FuncParamS struct {
	Type mstype.MSType		// Var type
	Iden VariableExpNodeS 	// Var name
}

// forces possible structs for StmtNode
func (Program) statmentPlaceholder() {}
func (BlockNodeS) statmentPlaceholder() {}
func (VarDeclNodeS) statmentPlaceholder() {}
func (ExStmtNodeS) statmentPlaceholder() {}
func (IfNodeS) statmentPlaceholder() {}
func (WhileNodeS) statmentPlaceholder() {}
func (ContinueNodeS) statmentPlaceholder() {}
func (BreakNodeS) statmentPlaceholder() {}
func (FuncDeclNodeS) statmentPlaceholder() {}
func (ReturnNodeS) statmentPlaceholder() {}


////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

func (fd *FuncDeclNodeS) GetFuncType() *mstype.MSOperationTypeS {
	typelist := []mstype.MSType{}
	for _, par := range fd.Params {
		typelist =	append(typelist, par.Type)
	}
	return &mstype.MSOperationTypeS{Left: typelist, Right: fd.Rt}
}

func (rs *ReturnNodeS) HasReturnValue() bool {
	return rs.Node != nil
}