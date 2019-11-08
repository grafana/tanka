package client

type ErrorNotFound struct {
	errOut string
}

func (e ErrorNotFound) Error() string {
	return e.errOut
}

type ErrorUnknownResource struct {
	errOut string
}

func (e ErrorUnknownResource) Error() string {
	return e.errOut
}
