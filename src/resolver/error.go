package resolver


type ResolveError struct {
	msg string
}

func (e ResolveError) Error() string {
	return "Resolving erroe: " + e.msg
}