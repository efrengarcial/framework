package user

//https://medium.com/rungo/error-handling-in-go-f0125de052f0
//https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
//https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee

var (
	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrLoginAlreadyUsed      = NewErrBadRequest( "Nombre de usuario ya utilizado!", "error.userexists", "userManagement")
	ErrEmailAlreadyUsed      = NewErrBadRequest( "Email ya utilizado!", "error.emailexists", "userManagement")
	ErrIdExist               = 	NewErrBadRequest("Un nuevo Usuario ya no puede tener un ID","error.idexists", "userManagement")
)

type IErrBadRequest interface {
	error
	GetTitle() string
	GetEntityName() string
}

type IErrCustomParameterized interface {
	error
	GetTitle() string
	GetParams() map[string]string
}

type errBase struct {
	//https://github.com/gin-gonic/gin/issues/274
	Title       string
	Details     string
	Status      int
	Code        string
	Href        string
}

type ErrBadRequest struct {
	errBase
	EntityName		string
}

type ErrCustomParameterized struct {
	errBase
	Params		map[string]string
}

func (e *ErrBadRequest) GetTitle() string {
	return e.Title
}

func (e *ErrBadRequest) GetEntityName() string {
	return e.EntityName
}

func (e *ErrBadRequest) Error() string {
	return e.Code
}


func (e *ErrCustomParameterized) GetTitle() string {
	return e.Title
}

func (e *ErrCustomParameterized) GetParams() map[string]string {
	return e.Params
}

func (e *ErrCustomParameterized) Error() string {
	return e.Code
}

func NewErrBadRequest(title, code, entityName string) *ErrBadRequest {
	return &ErrBadRequest{
		errBase{Title:title, Code: code},
		entityName,
	}
}

