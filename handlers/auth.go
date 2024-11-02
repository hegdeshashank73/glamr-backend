package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/services"
	"github.com/sirupsen/logrus"
)

func CreateMagiclinkHandler(c *gin.Context) {
	req := entities.CreateMagiclinkReq{}
	err := c.Bind(&req)

	if err != nil {
		logrus.Error(err)
		errors.GlamrErrorBadRequest().Respond(c)
		return
	}

	if derr := req.Validate(); derr != nil {
		derr.Respond(c)
		return
	}

	res, derr := services.CreateMagiclink(req)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, res)
}

func VerifyMagiclinkHandler(c *gin.Context) {
	req := entities.VerifyMagiclinkReq{}
	err := c.Bind(&req)

	if err != nil {
		logrus.Error(err)
		errors.GlamrErrorBadRequest().Respond(c)
		return
	}

	if derr := req.Validate(); derr != nil {
		derr.Respond(c)
		return
	}

	res, derr := services.VerifyMagiclink(req)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, res)
}
