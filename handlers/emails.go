package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/services"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/sirupsen/logrus"
)

func CreateEmailTemplateHandler(c *gin.Context) {
	st := time.Now()
	defer utils.LogTimeTaken("handlers.CreateEmailTemplateHandler", st)

	req := entities.CreateEmailTemplateReq{}
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error(err)
		errors.GlamrErrorBadRequest().Respond(c)
		return
	}

	res, derr := services.CreateEmailTemplate(req)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, res)
}
