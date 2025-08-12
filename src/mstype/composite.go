package mstype

type MSCompositeTypeS struct {
	Types []MSType
}

func (t *MSCompositeTypeS) Eq(o MSType) bool {

	switch other := o.(type){
	case *MSCompositeTypeS:

		// Check if the same length
		if len(t.Types) != len(other.Types) {
			return false
		}

		// Check all types if they are equal
		for i := range t.Types {
			if !t.Types[i].Eq(other.Types[i]){
				return false
			}
		}
		return true
	}
	return false
}

// func (t *MSCompositeTypeS) String() string {
// 	lefts := []string{}
// 	for _, l := range t.Types {
// 		lefts = append(lefts, l.String())
// 	}
// 	return "(" + strings.Join(lefts, ", ") + ")"
// }

func (t *MSCompositeTypeS) Nullable() bool {
	return false
}