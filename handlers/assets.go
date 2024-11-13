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

func AssetUploadHandler(c *gin.Context) {
	st := time.Now()
	defer utils.LogTimeTaken("handlers.AssetImageUploadHandler", st)

	req := entities.AssetUploadHandlerReq{}
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

	res, derr := services.GeneratePresignedURL(req)
	if derr != nil {
		derr.Respond(c)
		return
	}
	c.JSON(http.StatusOK, res)
}
