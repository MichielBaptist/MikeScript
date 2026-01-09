package mstype

import "fmt"

type MSNamedTypeS struct {
	Name string
	Depth int			// scope depth where defined, used to compare named types
}

func (t *MSNamedTypeS) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSNamedTypeS:		return other.Name == t.Name && other.Depth == t.Depth
	default:				return false
	}
}

func (t *MSNamedTypeS) String() string {
	return fmt.Sprintf("%v", t.Name)
}

func (t *MSNamedTypeS) Nullable() bool {
	return false
}
