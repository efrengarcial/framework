package service

import "errors"

//https://medium.com/rungo/error-handling-in-go-f0125de052f0
//https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

var (
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
)

type IErrBadRequest interface {
	error
	GetErrorKey() string
}

type IErrInternalServerError interface {
	error
	GetErrorKey() string
}

type errBase struct {
	Message string
	ErrorKey string
	EntityName string
}

type ErrBadRequest struct {
	errBase
}

type ErrLoginAlreadyUsed struct {
	ErrBadRequest
}

func (e *ErrBadRequest) GetErrorKey() string {
	return e.EntityName +"."+ e.ErrorKey
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

func NewErrBadRequest(message, errorKey string) *ErrBadRequest {
	return &ErrBadRequest{
		errBase{Message:message, ErrorKey: errorKey,
		},
	}
}
func NewErrLoginAlreadyUsed(message, errorKey string) *ErrLoginAlreadyUsed {
	return &ErrLoginAlreadyUsed{
		ErrBadRequest{
			errBase{Message:message, ErrorKey: errorKey,
			},
		},
	}
}

