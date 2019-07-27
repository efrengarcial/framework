package service

//https://medium.com/rungo/error-handling-in-go-f0125de052f0
//https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
//https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee

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

type ErrEmailAlreadyUsed struct {
	ErrBadRequest
}

func (e *ErrBadRequest) GetErrorKey() string {
	return e.EntityName +"."+ e.ErrorKey
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

func NewErrBadRequest(message,entityName string, errorKey string) *ErrBadRequest {
	return &ErrBadRequest{
		errBase{Message:message, EntityName:entityName, ErrorKey: errorKey,
		},
	}
}
func NewErrLoginAlreadyUsed() *ErrLoginAlreadyUsed {
	return &ErrLoginAlreadyUsed{
		ErrBadRequest{
			errBase{Message: "Nombre de inicio de sesión ya usado!",EntityName: "userManagement" , ErrorKey: "userexists",
			},
		},
	}
}
func NewErrEmailAlreadyUsed() *ErrEmailAlreadyUsed {
	return &ErrEmailAlreadyUsed{
		ErrBadRequest{
			errBase{Message: "Correo electrónico ya está en uso!",EntityName: "userManagement" , ErrorKey: "emailexists",
			},
		},
	}
}
