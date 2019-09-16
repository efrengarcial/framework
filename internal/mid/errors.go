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
	GetErrorKey() string
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
			var status int

			switch err.(type) {
			case iErrBadRequest:
				status = http.StatusBadRequest
				iError, _ := err.(iErrBadRequest)
				fmt.Println(iError.GetErrorKey())
			case validator.ValidationErrors:
				status = http.StatusBadRequest
				fmt.Println(err.Error())
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
			}

			// Put the error into response
			//c.IndentedJSON(parsedError.Code, parsedError)
			//c.Abort()
			// or c.AbortWithStatusJSON(parsedError.Code, parsedError)
			c.AbortWithStatusJSON(status, gin.H{
				"message": err.Error(),
				"status" : status,
			})

			return
		}

	}
}
