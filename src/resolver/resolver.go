package resolver

import (
	"fmt"
	"mikescript/src/ast"
)

type scope map[string]bool

func newScope() scope {
	return make(scope)
}

func NewMSResolver(ast *ast.Program) MSResolver{
	return MSResolver{
		Ast: ast,
		scopes: make([]scope, 4),
	}
}


func (r *MSResolver) SetAst(ast *ast.Program) {
	r.Ast = ast
}

func (r *MSResolver) Reset() {
	r.scopes = make([]scope, 0, 10)
	r.locals = make(map[ast.ExpNodeI]int)
}

type MSResolver struct {
	Ast *ast.Program
	scopes []scope
	locals map[ast.ExpNodeI]int
}

func (r *MSResolver) currentScope() *scope {
	if len(r.scopes) == 0 {
		return nil
	}
	return &r.scopes[len(r.scopes)-1]
}

func (r *MSResolver) enterScope() {
	r.scopes = append(r.scopes, newScope())
}

func (r *MSResolver) leaveScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *MSResolver) declare(name string) {
	// If there is a current scope, we set the
	// name to false (declared but not yet init)

	if current := r.currentScope() ; current != nil {

		if _, ok := (*current)[name] ; ok {
			// TODO: return error
			println(fmt.Sprintf("Re-definition of variable %s", name))
		}
		(*current)[name] = false
	}
}

func (r *MSResolver) define(name string) {
	// If there is a current scope, we set the
	// name to false (declared but not yet init)

	if current := r.currentScope() ; current != nil {
		(*current)[name] = true
	}
}

func (r *MSResolver) resolveLocal(v ast.VariableExpNodeS, name string) {
	
	// traverse scopes from top to bottom (len(scopes) - 1) -> 0
	for i := len(r.scopes) - 1 ; i >= 0 ; i-- {
		current := r.scopes[i]
		if _, ok := current[name] ; ok {
			r.resolveInterp(v, len(r.scopes) - 1 - i)
		}
	}
}

func (r *MSResolver) resolveInterp(v ast.ExpNodeI, depth int) {

	if _, ok := r.locals[v] ; ok{
		println("Tried assigning already existing expression in the local map...")
	} else {
		fmt.Printf(".   %v --> %d\n", v, depth)
	}
	r.locals[v] = depth
}

// --------------------------------------------------------
// resolve
// --------------------------------------------------------

func (r *MSResolver) Resolve() map[ast.ExpNodeI]int {
	r.resolveStatement(*r.Ast)
	return r.locals
}

func (r *MSResolver) resolveStatement(stm ast.StmtNodeI) {
	switch st := stm.(type) {
	case ast.Program:		r.resolveStatements(st.Statements)
	case ast.BlockNodeS: 	r.resolveBlockNode(st)
	case ast.VarDeclNodeS:	r.resolveVariableDeclaration(st)
	case ast.ExStmtNodeS:	r.resolveExpressionStatement(st)
	case ast.IfNodeS:		r.resolveIfNode(st)
	case ast.WhileNodeS:	r.resolveWhileNode(st)
	case ast.ReturnNodeS:	r.resolveExpression(st.Node)
	case ast.FuncDeclNodeS:	r.resolveFuncDeclaration(st)
	}
	
}

func (r *MSResolver) resolveExpression(n ast.ExpNodeI) {
	switch ex := n.(type){
	case ast.AssignmentNodeS:	r.resolveAssignmentExpression(ex)
	case ast.DeclAssignNodeS:	r.resolveDeclAssignExpression(ex)
	case ast.FuncAppNodeS:		r.resolveFuncAppExpression(ex)
	case ast.FuncCallNodeS:		r.resolveExpression(ex.Fun)
	case ast.BinaryExpNodeS:	r.resolveBinaryExpression(ex)
	case ast.LogicalExpNodeS:	r.resolveLogicalExpression(ex)
	case ast.UnaryExpNodeS:		r.resolveExpression(ex.Node)
	case ast.VariableExpNodeS:	r.resolveVariableExpression(ex)
	case ast.GroupExpNodeS:		r.resolveExpression(ex.Node)
	}
}

// --------------------------------------------------------
// statements
// --------------------------------------------------------

func (r *MSResolver) resolveExpressionStatement(n ast.ExStmtNodeS) {
	r.resolveExpression(n.Ex)
}

func (r *MSResolver) resolveStatements(stmts []ast.StmtNodeI) {
	for _, stmt := range stmts{
		r.resolveStatement(stmt)
	}
}

func (r *MSResolver) resolveBlockNode(n ast.BlockNodeS) {

	r.enterScope()
	r.resolveStatements(n.Statements)
	r.leaveScope()

}

func (r *MSResolver) resolveVariableDeclaration(n ast.VarDeclNodeS) {
	r.declare(n.VarName())
	r.define(n.VarName())
}

func (r *MSResolver) resolveFuncDeclaration(n ast.FuncDeclNodeS) {
	
	// Declare and define fname in current scope
	r.declare(n.Fname.VarName())
	r.define(n.Fname.VarName())

	// push scope before declaring params
	r.enterScope()
	for _, p := range n.Params {
		r.declare(p.VarName())
		r.define(p.VarName())
	}
	r.resolveStatements(n.Body.Statements)
	r.leaveScope()
}


func (r *MSResolver) resolveIfNode(n ast.IfNodeS) {
	r.resolveExpression(n.Condition)
	r.resolveStatement(n.ThenStmt)
	if n.ElseStmt != nil {
		r.resolveStatement(n.ElseStmt)
	}
}

func (r *MSResolver) resolveWhileNode(n ast.WhileNodeS) {
	r.resolveExpression(n.Condition)
	r.resolveStatement(n.Body)
}

// --------------------------------------------------------
// expressions
// --------------------------------------------------------


func (r *MSResolver) resolveVariableExpression(v ast.VariableExpNodeS) {

	// // Check if declared but not initialized (when initializer
	// // contains the variable itself). Not used at the moment.
	// current :=r.currentScope()
	// if current != nil && !(*current)[v.Name.Lexeme] {
	// 	msg := fmt.Sprintf("Cannot resolve variable in own initializer: %s", v.Name.Lexeme)
	// 	return ResolveError{msg: msg}
	// }
	r.resolveLocal(v, v.Name.Lexeme)
}

func (r *MSResolver) resolveAssignmentExpression(a ast.AssignmentNodeS) {
	r.resolveExpression(a.Exp)
	r.resolveLocal(a.Identifier, a.Identifier.Name.Lexeme)
}

func (r *MSResolver) resolveBinaryExpression(b ast.BinaryExpNodeS) {
	r.resolveExpression(b.Left)
	r.resolveExpression(b.Right)
}

func (r *MSResolver) resolveLogicalExpression(b ast.LogicalExpNodeS) {
	r.resolveExpression(b.Left)
	r.resolveExpression(b.Right)
}

func (r *MSResolver) resolveDeclAssignExpression(da ast.DeclAssignNodeS) {
	r.resolveExpression(da.Exp)
	r.declare(da.Identifier.VarName())
	r.define(da.Identifier.VarName())
}

func (r *MSResolver) resolveFuncAppExpression(fa ast.FuncAppNodeS) {
	r.resolveExpression(fa.Fun)
	for _, e := range fa.Args{
		r.resolveExpression(e)
	}
}