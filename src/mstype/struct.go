package mstype

import "fmt"

type MSStructTypeS struct {
	Name string
	Fields map[string]MSType
}

func (t *MSStructTypeS) Eq(o MSType) bool {

	other, ok := o.(*MSStructTypeS)

	if !ok {
		return false
	}

	if len(other.Fields) != len(t.Fields) {
		return false
	}

	for name, fieldT := range t.Fields {
		
		// get other
		ofieldT, ok := other.Fields[name]

		if !ok || !fieldT.Eq(ofieldT) {
			return false
		}
	}

	// all checks ok
	return true
}

func (t *MSStructTypeS) String() string {
	return fmt.Sprintf("%v{%v}", t.Name, t.Fields)
}

func (t *MSStructTypeS) Nullable() bool {
	return true
}
