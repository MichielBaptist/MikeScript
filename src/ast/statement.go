package ast

import (
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
	Vartype 	token.Token			// Type
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

type FuncDeclNodeS struct {
	Fname VariableExpNodeS				// Name
	Params []FuncParamS 				// Parameters
	Rt token.Token						// Return type
	Body *BlockNodeS					// Body of function, may be nil
}

type FuncParamS struct {
	Type token.Token		// Type token
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
