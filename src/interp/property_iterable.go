package interp

type MSIterable interface {
	Elems() ([]MSVal, error)
	Len() (MSVal, error)
	From([]MSVal) (MSVal, error)
}