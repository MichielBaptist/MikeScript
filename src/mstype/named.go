package mstype

type MSNamedTypeS struct {
	Name string
}

func (t *MSNamedTypeS) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSNamedTypeS:		return other.Name == t.Name
	default:				return false
	}
}

func (t *MSNamedTypeS) String() string {
	return t.Name
}

func (t *MSNamedTypeS) Nullable() bool {
	return false
}
