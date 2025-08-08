package interp

type MSIndexable interface {
	Get(at MSVal) (MSVal, error)
	Set(at MSVal, val MSVal) (MSVal, error)
	ValidIndex(idx MSVal) error
	ValidValue(val MSVal) error
}