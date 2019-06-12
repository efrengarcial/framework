package service


type ErrBadRequest struct {
	Message string
	Key string
}

func (e *ErrBadRequest) BadRequest () string {
	return e.Key
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

