package mstype

import "fmt"

type MSArrayType struct {
	Type MSType
}

func (t *MSArrayType) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSArrayType:	return t.Type.Eq(other.Type)
	default:			return false
	}
}

func (t *MSArrayType) String() string {
	return fmt.Sprintf("%s[]", t.Type.String())
}

func (t *MSArrayType) Nullable() bool {
	return false
}