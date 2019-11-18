package client

// ErrorNotFound means that the requested object is not found on the server
type ErrorNotFound struct {
	errOut string
}

func (e ErrorNotFound) Error() string {
	return e.errOut
}

// ErrorUnknownResource means that the requested resource type is unknown to the
// server
type ErrorUnknownResource struct {
	errOut string
}

func (e ErrorUnknownResource) Error() string {
	return e.errOut
}
