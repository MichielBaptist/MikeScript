package resolver

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
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
	r.vlocals = make(map[*ast.VariableExpNodeS]int)
	r.tlocals = make(map[*mstype.MSNamedTypeS]int)
}

type MSResolver struct {
	Ast *ast.Program
	scopes []scope
	vlocals map[*ast.VariableExpNodeS]int
	tlocals map[*mstype.MSNamedTypeS]int
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

	fmt.Printf(".   Setting: scopes[%d][%s] --> false\n", len(r.scopes) - 1, name)
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

	fmt.Printf(".   Setting: scopes[%d][%s] --> true\n", len(r.scopes) -1, name)
	if current := r.currentScope() ; current != nil {
		(*current)[name] = true
	}
}

func (r *MSResolver) resolveLocalVariable(v *ast.VariableExpNodeS, name string) {

	depth, found, _ := r.findName(name)

	// Found nowhere in scopes
	if !found {
		fmt.Printf(".   VARI: %p %v --> Global\n", v, v)
		return
	}

	if _, ok := r.vlocals[v] ; ok {
		println("Tried assigning already existing expression in the local map...")
	} else {
		fmt.Printf(".   VARI: %p %v --> %d\n", v, v, depth)
	}
	r.vlocals[v] = depth
}

func (r *MSResolver) resolveLocalType(t *mstype.MSNamedTypeS, name string) {

	depth, found, _ := r.findName(name)

	if !found {
		fmt.Printf(".   TYPE: %p %v --> Global\n", t, t)
		return
	}

	if _, ok := r.tlocals[t] ; ok{
		println("Tried assigning already existing expression in the local map...")
	} else {
		fmt.Printf(".   TYPE: %p %v --> %d\n", t, t, depth)
	}
	r.tlocals[t] = depth

	// also add depth to named type
	t.Depth = depth
}

func (r *MSResolver) findName(name string) (int, bool, bool) {
	/* Walk back scope stack to look for name */

	var l int = len(r.scopes) - 1
	var i int = l
	var s scope

	for ; i >= 0 ; i-- {
		s = r.scopes[i]

		// check if it's in this scope
		if defined, ok := s[name] ; ok {
			return l - i, true, defined
		}
	}

	return l - i, false, false
}

// --------------------------------------------------------
// resolve
// --------------------------------------------------------

func (r *MSResolver) Resolve() (map[*ast.VariableExpNodeS]int, map[*mstype.MSNamedTypeS]int) {
	r.resolveStatement(r.Ast)
	return r.vlocals, r.tlocals
}

func (r *MSResolver) resolveStatement(stm ast.StmtNodeI) {
	// fmt.Printf("%p // %#v\n", stm, stm)
	switch st := stm.(type) {
	case *ast.Program:					r.resolveStatements(st.Statements)
	case *ast.BlockNodeS: 				r.resolveBlockNode(st)
	case *ast.VarDeclNodeS:				r.resolveVariableDeclaration(st)
	case *ast.ExStmtNodeS:				r.resolveExpression(st.Ex)
	case *ast.IfNodeS:					r.resolveIfNode(st)
	case *ast.WhileNodeS:				r.resolveWhileNode(st)
	case *ast.ReturnNodeS:				r.resolveExpression(st.Node)
	case *ast.FuncDeclNodeS:			r.resolveFuncDeclaration(st)
	case *ast.TypeDefStatementS: 		r.resolveTypeDeclaration(st)
	case *ast.StructDeclarationNodeS:	r.resolveStructDeclaration(st)
	case *ast.BreakNodeS:				return 	// nothing to resolve
	case *ast.ContinueNodeS:			return 	// nothing to resolve
	default:							fmt.Printf("Resolving: %v\n", st); _ = []int{}[0]
	}
}

func (r *MSResolver) resolveExpression(n ast.ExpNodeI) {
	// fmt.Printf("%p // %#v\n", n, n)
	switch ex := n.(type){
	case *ast.AssignmentNodeS:				r.resolveAssignmentExpression(ex)
	case *ast.DeclAssignNodeS:				r.resolveDeclAssignExpression(ex)
	case *ast.FuncAppNodeS:					r.resolveFuncAppExpression(ex)
	case *ast.FuncCallNodeS:				r.resolveExpression(ex.Fun)
	case *ast.BinaryExpNodeS:				r.resolveBinaryExpression(ex)
	case *ast.LogicalExpNodeS:				r.resolveLogicalExpression(ex)
	case *ast.UnaryExpNodeS:				r.resolveExpression(ex.Node)
	case *ast.TupleNodeS:					r.resolveExpressions(ex.Expressions)
	case *ast.VariableExpNodeS:				r.resolveVariableExpression(ex)
	case *ast.GroupExpNodeS:				r.resolveExpression(ex.Node)
	case *ast.ArrayConstructorNodeS:		r.resolveArrayConstructor(ex)
	case *ast.ArrayIndexNodeS:				r.resolveArrayIndex(ex)
	case *ast.ArrayAssignmentNodeS:			r.resolveArrayAssignment(ex)
	case *ast.FieldAccessNodeS:				r.resolveExpression(ex.Target)
	case *ast.StructConstructorNodeS:		r.resolveStructConstructor(ex)
	case *ast.FieldAssignmentNode:			r.resolveFieldAssignment(ex)
	case *ast.LiteralExpNodeS:				return 	// nothing to resolve
	case *ast.IterableFuncCallNodeS:		r.resolveIterableFuncCallNode(ex)
	case *ast.IterableFuncAppNodeS:			r.resolveIterableFuncApplication(ex)
	case *ast.IterableFuncAppAndCallNodeS:	r.resolveIterableFuncAppAndCall(ex)
	case *ast.RangeConstructorNodeS:		r.resolveRangeConstructor(ex)
	case *ast.StarredExpNodeS:				r.resolveExpression(ex.Node)
	default:								fmt.Printf("%v\n", ex) ; _ = []int{}[0]
	}
}

func (r *MSResolver) resolveRangeConstructor(n *ast.RangeConstructorNodeS) {
	if n.From != nil {
		r.resolveExpression(n.From)
	}
	if n.To != nil {
		r.resolveExpression(n.To)
	}
}

func (r *MSResolver) resolveIterableFuncAppAndCall(n *ast.IterableFuncAppAndCallNodeS) {
	r.resolveExpression(n.Fun)
	r.resolveExpression(n.Args)
}

func (r *MSResolver) resolveIterableFuncApplication(n *ast.IterableFuncAppNodeS) {
	r.resolveExpression(n.Fun)
	r.resolveExpression(n.Args)
}

func (r *MSResolver) resolveIterableFuncCallNode(n *ast.IterableFuncCallNodeS) {
	r.resolveExpression(n.Fun)
}

func (r *MSResolver) resolveType(n mstype.MSType) {
	// fmt.Printf("%p // %#v\n", n, n)
	switch t := n.(type){
	case *mstype.MSSimpleTypeS:		return
	case *mstype.MSCompositeTypeS:	r.resolveTypes(t.Types)
	case *mstype.MSArrayType:		r.resolveType(t.Type)
	case *mstype.MSOperationTypeS:	r.resolveOperationType(t)
	case *mstype.MSStructTypeS:		r.resolveStructType(t)
	case *mstype.MSNamedTypeS:		r.resolveNamedType(t)
	default:						_ = []int{}[0]
	}
}

// --------------------------------------------------------
// statements
// --------------------------------------------------------

func (r *MSResolver) resolveStructDeclaration(sd *ast.StructDeclarationNodeS) {
	r.declare(sd.Name.VarName())
	for _, field := range sd.Fields {
		r.resolveType(field)
	}
	r.define(sd.Name.VarName())
}

func (r *MSResolver) resolveTypeDeclaration(td *ast.TypeDefStatementS) {
	// Note: declare before resolve to detect recursive type defs
	r.declare(td.Tname.VarName())
	r.resolveType(td.Type)
	r.define(td.Tname.VarName())
}


func (r *MSResolver) resolveStatements(stmts []ast.StmtNodeI) {
	for _, stmt := range stmts{
		r.resolveStatement(stmt)
	}
}

func (r *MSResolver) resolveBlockNode(n *ast.BlockNodeS) {
	r.enterScope()
	r.resolveStatements(n.Statements)
	r.leaveScope()
}

func (r *MSResolver) resolveVariableDeclaration(n *ast.VarDeclNodeS) {
	r.declare(n.VarName())
	r.define(n.VarName())
	r.resolveType(n.Vartype)
}

func (r *MSResolver) resolveFuncDeclaration(n *ast.FuncDeclNodeS) {
	
	// Declare and define fname in current scope
	r.declare(n.Fname.VarName())
	r.define(n.Fname.VarName())

	// Resolve function types
	for _, t := range n.Params {
		r.resolveType(t.Type)
	}
	r.resolveType(n.Rt)

	// push scope before declaring params
	r.enterScope()
	for _, p := range n.Params {
		r.declare(p.VarName())
		r.define(p.VarName())
	}
	r.resolveStatements(n.Body.Statements)
	r.leaveScope()
}


func (r *MSResolver) resolveIfNode(n *ast.IfNodeS) {
	r.resolveExpression(n.Condition)
	r.resolveStatement(n.ThenStmt)
	if n.ElseStmt != nil {
		r.resolveStatement(n.ElseStmt)
	}
}

func (r *MSResolver) resolveWhileNode(n *ast.WhileNodeS) {
	r.resolveExpression(n.Condition)
	r.resolveStatement(n.Body)
}

func (r *MSResolver) resolveExpressions(es []ast.ExpNodeI) {
	for _, e := range es {
		r.resolveExpression(e)
	}
}

// --------------------------------------------------------
// expressions
// --------------------------------------------------------

func (r *MSResolver) resolveFieldAssignment(n *ast.FieldAssignmentNode) {
	r.resolveExpression(n.Target)
	r.resolveExpression(n.Value)
}

func (r *MSResolver) resolveStructConstructor(n *ast.StructConstructorNodeS) {
	for _, exp := range n.Fields {
		r.resolveExpression(exp)
	}
}

func (r *MSResolver) resolveArrayAssignment(n *ast.ArrayAssignmentNodeS) {
	r.resolveExpression(n.Index)
	r.resolveExpression(n.Target)
	r.resolveExpression(n.Value)
}

func (r *MSResolver) resolveArrayConstructor(n *ast.ArrayConstructorNodeS) {
	if n.N != nil {
		r.resolveExpression(n.N)
	}
	r.resolveExpressions(n.Vals)
	r.resolveType(n.Type)
}

func (r *MSResolver) resolveArrayIndex(n *ast.ArrayIndexNodeS) {
	r.resolveExpression(n.Target)
	r.resolveExpression(n.Index)
}

func (r *MSResolver) resolveVariableExpression(v *ast.VariableExpNodeS) {

	// // Check if declared but not initialized (when initializer
	// // contains the variable itself). Not used at the moment.
	// current :=r.currentScope()
	// if current != nil && !(*current)[v.Name.Lexeme] {
	// 	msg := fmt.Sprintf("Cannot resolve variable in own initializer: %s", v.Name.Lexeme)
	// 	return ResolveError{msg: msg}
	// }
	r.resolveLocalVariable(v, v.Name.Lexeme)
}

func (r *MSResolver) resolveAssignmentExpression(a *ast.AssignmentNodeS) {
	r.resolveExpression(a.Exp)
	r.resolveLocalVariable(a.Identifier, a.Identifier.Name.Lexeme)
}

func (r *MSResolver) resolveBinaryExpression(b *ast.BinaryExpNodeS) {
	r.resolveExpression(b.Left)
	r.resolveExpression(b.Right)
}

func (r *MSResolver) resolveLogicalExpression(b *ast.LogicalExpNodeS) {
	r.resolveExpression(b.Left)
	r.resolveExpression(b.Right)
}

func (r *MSResolver) resolveDeclAssignExpression(da *ast.DeclAssignNodeS) {
	r.resolveExpression(da.Exp)
	r.declare(da.Identifier.VarName())
	r.define(da.Identifier.VarName())
}

func (r *MSResolver) resolveFuncAppExpression(fa *ast.FuncAppNodeS) {
	r.resolveExpression(fa.Fun)
	r.resolveExpressions(fa.Args)
}

func (r *MSResolver) resolveTypes(ts []mstype.MSType) {
	for _, t := range ts {
		r.resolveType(t)
	}
}

func (r *MSResolver) resolveOperationType(ot *mstype.MSOperationTypeS) {
	r.resolveTypes(ot.Left)
	r.resolveType(ot.Right)
}

func (r *MSResolver) resolveStructType(st *mstype.MSStructTypeS) {
	for _, t := range st.Fields {
		r.resolveType(t)
	}
}

func (r *MSResolver) resolveNamedType(nt *mstype.MSNamedTypeS) {
	r.resolveLocalType(nt, nt.Name)
}