package mstype

type MSSimpleTypeS struct {
	Rt ResultType
}

func (t *MSSimpleTypeS) Eq(o MSType) bool {
	switch other := o.(type) {
	case *MSSimpleTypeS:	return t.Rt == other.Rt
	default:				return false
	}
}

func (t *MSSimpleTypeS) String() string {
	return t.Rt.String()
}

func (t *MSSimpleTypeS) Nullable() bool {
	return false
}
