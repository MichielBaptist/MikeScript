package mstype

type MSNamedTypeS struct {
	Name string
	Ref MSType 		// is nil if not resolved
}

func (t *MSNamedTypeS) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSNamedTypeS:		return other.Name == t.Name
	default:				return false
	}
}

// func (t *MSNamedTypeS) String() string {
// 	return fmt.Sprintf("%v{%v}", t.Name, t.Ref)
// }

func (t *MSNamedTypeS) Nullable() bool {
	return false
}
