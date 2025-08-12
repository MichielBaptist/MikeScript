package mstype

type MSOperationTypeS struct {
	Left []MSType	// param types
	Right MSType 	// output types
}

func (t *MSOperationTypeS) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSOperationTypeS:

		// Compare output of op
		e := t.Right.Eq(other.Right)

		// Compare parameters of op, note this
		// circuit breaks on a single false
		e = e && len(t.Left) == len(other.Left)
		for i := 0; e && i < len(t.Left) ; i++ {
			e = e && t.Left[i].Eq(other.Left[i])
		}
		return e
	}
	return false
}

// func (t *MSOperationTypeS) String() string{
// 	lefts := []string{}
// 	for _, l := range t.Left {
// 		lefts = append(lefts, l.String())
// 	}
// 	return fmt.Sprintf("(%s -> %s)", strings.Join(lefts, ", "), t.Right.String())
// }

func (t *MSOperationTypeS) Nullable() bool {
	return true
}