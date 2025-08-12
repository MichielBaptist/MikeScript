package resolver

// import (
// 	"fmt"
// 	"mikescript/src/ast"
// 	"mikescript/src/mstype"
// )

// type TypeScope map[string]mstype.MSType

// func newTypeScope() TypeScope {
// 	return make(TypeScope)
// }

// func NewMSTypeResolver(ast *ast.Program) MSTypeResolver{
// 	return MSTypeResolver{
// 		Ast: ast,
// 		TypeScopes: make([]TypeScope, 4),
// 	}
// }

// func (r *MSTypeResolver) SetAst(ast *ast.Program) {
// 	r.Ast = ast
// }

// func (r *MSTypeResolver) Reset() {
// 	r.TypeScopes = make([]TypeScope, 0, 10)
// }

// type MSTypeResolver struct {
// 	Ast *ast.Program
// 	TypeScopes []TypeScope
// }

// func (r *MSTypeResolver) currentTypeScope() *TypeScope {
// 	if len(r.TypeScopes) == 0 {
// 		return nil
// 	}
// 	return &r.TypeScopes[len(r.TypeScopes)-1]
// }

// func (r *MSTypeResolver) enterTypeScope() {
// 	r.TypeScopes = append(r.TypeScopes, newTypeScope())
// }

// func (r *MSTypeResolver) leaveTypeScope() {
// 	r.TypeScopes = r.TypeScopes[:len(r.TypeScopes)-1]
// }

// func (r *MSTypeResolver) declare(name string, typ mstype.MSType) {
// 	fmt.Printf(".   TYPE: ++ Declared '%s' --> %p '%+v'\n", name, typ, typ)
// 	(*r.currentTypeScope())[name] = typ
// }

// func (r *MSTypeResolver) resolveLocalType(nt *mstype.MSNamedTypeS) {

// 	// Look for the first occurence of the name in scope stack

// 	var l int = len(r.TypeScopes) - 1
// 	var i int = l
// 	var name string = nt.Name

// 	for ; i >= 0 ; i-- {

// 		if typ, ok := r.TypeScopes[i][name] ; ok {
// 			fmt.Printf(".   TYPE: -- resolved %p '%+v' to %p == '%v'\n", nt, nt, typ, typ)
// 			nt.Ref = typ
// 			return
// 		}
// 	}

// 	fmt.Printf("ERROR: UNABLE TO RESOLVE NAMED TYPE '%s'\n", nt)

// }

// // --------------------------------------------------------
// // resolve
// // --------------------------------------------------------

// func (r *MSTypeResolver) Resolve() {
// 	r.enterTypeScope()
// 	r.resolveStatement(r.Ast)
// 	r.leaveTypeScope()
// }

// func (r *MSTypeResolver) resolveStatement(stm ast.StmtNodeI) {
// 	// fmt.Printf("%p // %#v\n", stm, stm)
// 	switch st := stm.(type) {
// 	case *ast.Program:					r.resolveStatements(st.Statements)
// 	case *ast.BlockNodeS: 				r.resolveBlockNode(st)
// 	case *ast.VarDeclNodeS:				r.resolveVariableDeclaration(st)
// 	case *ast.ExStmtNodeS:				r.resolveExpression(st.Ex)
// 	case *ast.IfNodeS:					r.resolveIfNode(st)
// 	case *ast.WhileNodeS:				r.resolveWhileNode(st)
// 	case *ast.ReturnNodeS:				r.resolveExpression(st.Node)
// 	case *ast.FuncDeclNodeS:			r.resolveFuncDeclaration(st)
// 	case *ast.TypeDefStatementS: 		r.resolveTypeDeclaration(st)
// 	case *ast.StructDeclarationNodeS:	r.resolveStructDeclaration(st)
// 	case *ast.BreakNodeS:				return 	// nothing to resolve
// 	default:							_ = []int{}[0]
// 	}
// }

// func (r *MSTypeResolver) resolveExpression(n ast.ExpNodeI) {
// 	// fmt.Printf("%p // %#v\n", n, n)
// 	switch ex := n.(type){
// 	case *ast.AssignmentNodeS:			r.resolveAssignmentExpression(ex)
// 	case *ast.DeclAssignNodeS:			r.resolveDeclAssignExpression(ex)
// 	case *ast.FuncAppNodeS:				r.resolveFuncAppExpression(ex)
// 	case *ast.FuncCallNodeS:			r.resolveExpression(ex.Fun)
// 	case *ast.BinaryExpNodeS:			r.resolveBinaryExpression(ex)
// 	case *ast.LogicalExpNodeS:			r.resolveLogicalExpression(ex)
// 	case *ast.UnaryExpNodeS:			r.resolveExpression(ex.Node)
// 	case *ast.TupleNodeS:				r.resolveExpressions(ex.Expressions)
// 	case *ast.VariableExpNodeS:			r.resolveVariableExpression(ex)
// 	case *ast.GroupExpNodeS:			r.resolveExpression(ex.Node)
// 	case *ast.ArrayConstructorNodeS:	r.resolveArrayConstructor(ex)
// 	case *ast.ArrayIndexNodeS:			r.resolveArrayIndex(ex)
// 	case *ast.ArrayAssignmentNodeS:		r.resolveArrayAssignment(ex)
// 	case *ast.FieldAccessNodeS:			r.resolveExpression(ex.Target)
// 	case *ast.StructConstructorNodeS:	r.resolveStructConstructor(ex)
// 	case *ast.FieldAssignmentNode:		r.resolveFieldAssignment(ex)
// 	case *ast.LiteralExpNodeS:			return 	// nothing to resolve
// 	default:							fmt.Printf("%v\n", ex) ; _ = []int{}[0]
// 	}
// }

// func (r *MSTypeResolver) resolveType(n mstype.MSType) {
// 	// fmt.Printf("%p // %#v\n", n, n)
// 	switch t := n.(type){
// 	case *mstype.MSSimpleTypeS:		return
// 	case *mstype.MSCompositeTypeS:	r.resolveTypes(t.Types)
// 	case *mstype.MSArrayType:		r.resolveType(t.Type)
// 	case *mstype.MSOperationTypeS:	r.resolveOperationType(t)
// 	case *mstype.MSStructTypeS:		r.resolveStructType(t)
// 	case *mstype.MSNamedTypeS:		r.resolveNamedType(t)
// 	default:						_ = []int{}[0]
// 	}
// }

// // --------------------------------------------------------
// // statements
// // --------------------------------------------------------

// func (r *MSTypeResolver) resolveStructDeclaration(sd *ast.StructDeclarationNodeS) {

// 	// Convert StructDecl to type
// 	fields := make(map[string]mstype.MSType)
// 	for name, field := range sd.Fields {
// 		fields[name.VarName()] = field
// 	}
// 	sdt := &mstype.MSStructTypeS{Name: sd.Name.VarName(), Fields: fields}

// 	// declare the struct as a type
// 	r.declare(sd.Name.VarName(), sdt)

// 	// resolve all types in struct
// 	for _, field := range sd.Fields {
// 		r.resolveType(field)
// 	}
// }

// func (r *MSTypeResolver) resolveTypeDeclaration(td *ast.TypeDefStatementS) {
// 	r.declare(td.Tname.VarName(), td.Type)
// 	r.resolveType(td.Type)
// }

// func (r *MSTypeResolver) resolveStatements(stmts []ast.StmtNodeI) {
// 	for _, stmt := range stmts{
// 		r.resolveStatement(stmt)
// 	}
// }

// func (r *MSTypeResolver) resolveBlockNode(n *ast.BlockNodeS) {
// 	r.enterTypeScope()
// 	r.resolveStatements(n.Statements)
// 	r.leaveTypeScope()
// }

// func (r *MSTypeResolver) resolveVariableDeclaration(n *ast.VarDeclNodeS) {
// 	r.resolveType(n.Vartype)
// }

// func (r *MSTypeResolver) resolveFuncDeclaration(n *ast.FuncDeclNodeS) {

// 	// Resolve argument types, they are
// 	// evaluated outside the function
// 	for _, p := range n.Params {
// 		r.resolveType(p.Type)
// 	}

// 	r.enterTypeScope()
// 	r.resolveStatements(n.Body.Statements)
// 	r.leaveTypeScope()
// }

// func (r *MSTypeResolver) resolveIfNode(n *ast.IfNodeS) {
// 	r.resolveExpression(n.Condition)
// 	r.resolveStatement(n.ThenStmt)
// 	if n.ElseStmt != nil {
// 		r.resolveStatement(n.ElseStmt)
// 	}
// }

// func (r *MSTypeResolver) resolveWhileNode(n *ast.WhileNodeS) {
// 	r.resolveExpression(n.Condition)
// 	r.resolveStatement(n.Body)
// }

// func (r *MSTypeResolver) resolveExpressions(es []ast.ExpNodeI) {
// 	for _, e := range es {
// 		r.resolveExpression(e)
// 	}
// }

// // --------------------------------------------------------
// // expressions
// // --------------------------------------------------------

// func (r *MSTypeResolver) resolveFieldAssignment(n *ast.FieldAssignmentNode) {
// 	r.resolveExpression(n.Target)
// 	r.resolveExpression(n.Value)
// }

// func (r *MSTypeResolver) resolveStructConstructor(n *ast.StructConstructorNodeS) {
// 	r.resolveType(n.Name)
// 	for _, exp := range n.Fields {
// 		r.resolveExpression(exp)
// 	}
// }

// func (r *MSTypeResolver) resolveArrayAssignment(n *ast.ArrayAssignmentNodeS) {
// 	r.resolveExpression(n.Index)
// 	r.resolveExpression(n.Target)
// 	r.resolveExpression(n.Value)
// }

// func (r *MSTypeResolver) resolveArrayConstructor(n *ast.ArrayConstructorNodeS) {
// 	r.resolveType(n.Type)
// 	if n.N != nil {
// 		r.resolveExpression(n.N)
// 	}
// 	r.resolveExpressions(n.Vals)
// }

// func (r *MSTypeResolver) resolveArrayIndex(n *ast.ArrayIndexNodeS) {
// 	r.resolveExpression(n.Target)
// 	r.resolveExpression(n.Index)
// }

// func (r *MSTypeResolver) resolveVariableExpression(v *ast.VariableExpNodeS) {
// }

// func (r *MSTypeResolver) resolveAssignmentExpression(a *ast.AssignmentNodeS) {
// 	r.resolveExpression(a.Exp)
// }

// func (r *MSTypeResolver) resolveBinaryExpression(b *ast.BinaryExpNodeS) {
// 	r.resolveExpression(b.Left)
// 	r.resolveExpression(b.Right)
// }

// func (r *MSTypeResolver) resolveLogicalExpression(b *ast.LogicalExpNodeS) {
// 	r.resolveExpression(b.Left)
// 	r.resolveExpression(b.Right)
// }

// func (r *MSTypeResolver) resolveDeclAssignExpression(da *ast.DeclAssignNodeS) {
// 	r.resolveExpression(da.Exp)
// }

// func (r *MSTypeResolver) resolveFuncAppExpression(fa *ast.FuncAppNodeS) {
// 	r.resolveExpression(fa.Fun)
// 	r.resolveExpressions(fa.Args)
// }

// // --------------------------------------------------------
// // types
// // --------------------------------------------------------
