package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Error string `json:"error"`
}

type glamrError struct {
	reason     string
	httpStatus int
}

type GlamrError interface {
	Respond(c *gin.Context)
	IsNotFound() bool
	ToResponseError() ResponseError
}

func (d glamrError) ToResponseError() ResponseError {
	return ResponseError{Error: d.reason}
}

func (d glamrError) Respond(c *gin.Context) {
	c.JSON(d.httpStatus, d.ToResponseError())
}

func (d glamrError) IsNotFound() bool {
	return d.httpStatus == http.StatusNotFound
}

func Error400(message string) ResponseError {
	return ResponseError{Error: message}
}

func Error400UnsupportedValue(field string) ResponseError {
	return ResponseError{Error: fmt.Sprintf("unsupported value in %s", field)}
}

func Error404EntityDNE(entity string) ResponseError {
	return ResponseError{Error: fmt.Sprintf("%s does not exist", entity)}
}

func Error5xx() ResponseError {
	return ResponseError{Error: "Something went wrong!"}
}

func Error500InternalServer() ResponseError {
	return ResponseError{Error: "Internal Server Error"}
}
