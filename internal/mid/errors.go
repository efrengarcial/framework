package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"strings"
)

type iErrBadRequest interface {
	GetTitle() string
	GetEntityName() string
}

type iErrCustomParameterized interface {
	GetTitle() string
	GetParams() map[string]string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Middleware Error Handler in server package
func Error(logger *logrus.Logger) gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny, logger)
}

func jsonAppErrorReporterT(errType gin.ErrorType, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {

			err := detectedErrors[0].Err
			var (
				status     int
				title      string
				entityName string
				params     map[string]string
			)

			switch err.(type) {
			case iErrBadRequest:
				status = http.StatusBadRequest
				iError, _ := err.(iErrBadRequest)
				title = iError.GetTitle()
				entityName = iError.GetEntityName()
				c.Header("X-app-alert", err.Error())
				c.Header("x-app-params", entityName)
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status" : status,
					"title" :  title,
					"params" : entityName,
					"entityName" : entityName,
					"errorKey" :  err.Error(),
				})
			case iErrCustomParameterized:
				status = http.StatusBadRequest
				iError, _ := err.(iErrCustomParameterized)
				title = iError.GetTitle()
				params = iError.GetParams()
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status" : status,
					"title" : title,
					"params" : params,
				})
			case validator.ValidationErrors:
				status = http.StatusBadRequest
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status" : status,
				})
			default:
				var errorLog strings.Builder
				errorLog.WriteString(err.Error() + "\n\n")
				if err, ok := err.(stackTracer); ok {
					for _, f := range err.StackTrace() {
						errorLog.WriteString(fmt.Sprintf("%+v \n\n", f))
					}
				}
				logger.Error("error", errorLog.String())
				fmt.Printf("with stack trace => %+v \n\n", err)
				status = http.StatusInternalServerError
				c.AbortWithStatusJSON(status, gin.H{
					"message": "error.internalServerError",
					"status" : status,
					"title" : "Error Interno del Servidor",
				})
			}

			return
		}
	}
}
