package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/services"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/spf13/viper"
)

func setContextAndProceed(c *gin.Context, authHeader string, shouldAbort bool) {
	token := authHeader
	if strings.Contains(authHeader, " ") {
		token = strings.Split(authHeader, " ")[1]
	}
	pRet, derr := services.GetPersonFromToken(token)
	if derr != nil && shouldAbort {
		errors.GlamrUnauthenticated().Respond(c)
		c.Abort()
		return
	}

	c.Set("person", pRet.Person)
	c.Set("token", token)
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		defer utils.LogTimeTaken("middlewares.Auth", st)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		setContextAndProceed(c, authHeader, true)
		c.Next()
	}
}

func Private(shouldAbort bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		defer utils.LogTimeTaken("middlewares.Private", st)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" && shouldAbort {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		isPresent := false
		for _, password := range viper.GetStringSlice("PRIVATE_PASSWORD") {
			if authHeader == "Bearer "+password || authHeader == password {
				isPresent = true
			}
		}
		if !isPresent && shouldAbort {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("is_private", isPresent)
		c.Next()
	}
}

func AuthNoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		defer utils.LogTimeTaken("middlewares.AuthNoAuth", st)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		setContextAndProceed(c, authHeader, true)
		c.Next()
	}
}

func AuthOrPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		defer utils.LogTimeTaken("middlewares.AuthOrPrivate", st)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		setContextAndProceed(c, authHeader, false)
		Private(false)(c)
		c.Next()
	}
}
