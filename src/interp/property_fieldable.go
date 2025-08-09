package interp

type MSFieldable interface {
	Get(field string) (MSVal, error)
	Set(field string, val MSVal) (MSVal, error)
	ValidField(field string) error
	ValidValue(field string, val MSVal) error
}