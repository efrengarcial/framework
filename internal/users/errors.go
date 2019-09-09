package users

import "errors"

//https://medium.com/rungo/error-handling-in-go-f0125de052f0
//https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
//https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee

var (
	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")
	ErrLoginAlreadyUsed = NewErrBadRequest( "Nombre de inicio de sesión ya usado!", "userManagement" ,  "userexists")
	ErrEmailAlreadyUsed = NewErrBadRequest( "Correo electrónico ya está en uso!", "userManagement" ,  "emailexists")
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

	Status      int         `json:"status"`
	Code        string      `json:"code"`
	Title       string      `json:"title"`
	Details     string      `json:"details"`
	Href        string      `json:"href"`
}

type ErrBadRequest struct {
	errBase
}

func (e *ErrBadRequest) GetErrorKey() string {
	return e.EntityName +"."+ e.ErrorKey
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

func NewErrBadRequest(message, entityName, errorKey string) *ErrBadRequest {
	return &ErrBadRequest{
		errBase{Message:message, EntityName:entityName, ErrorKey: errorKey,
		},
	}
}
