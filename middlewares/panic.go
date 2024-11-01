package middlewares

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/sirupsen/logrus"
)

func HandlePanic(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(string(debug.Stack()))
			logrus.Error(err)
			errors.GlamrErrorInternalServerError().Respond(c)
		}
	}()
	c.Next()
}
