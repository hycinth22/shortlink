package server

type BadRequestError struct {
	Msg string
}

func (c BadRequestError) Error() string {
	return c.Msg
}

type ServerError struct {
	Err error
}

func (c ServerError) Error() string {
	return c.Err.Error()
}
