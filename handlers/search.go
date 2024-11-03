package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/services"
	"github.com/hegdeshashank73/glamr-backend/utils"
)

func GetSearchOptionsHandler(c *gin.Context) {
	st := time.Now()
	defer utils.LogTimeTaken("handlers.SearchOptionsHandler", st)

	val, _ := c.Get("person")
	person := val.(entities.Person)
	s3Key := c.Query("s3_key")
	countryCode := c.Query("country")
	req := entities.SearchOptionsReq{}
	req.S3Key = s3Key
	if countryCode == "" {
		countryCode = "US"
	}
	req.CountryCode = countryCode

	res, derr := services.GetSerpAPISearchResults(&person, req)
	if derr != nil {
		derr.Respond(c)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetSearchHistoryHandler(c *gin.Context) {
	st := time.Now()
	defer utils.LogTimeTaken("handlers.GetSearchHistoryHandler", st)

	val, _ := c.Get("person")
	person := val.(entities.Person)

	res, derr := services.GetSearchHistory(&person)
	if derr != nil {
		derr.Respond(c)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetSearchHistoryOptionsHandler(c *gin.Context) {
	st := time.Now()
	defer utils.LogTimeTaken("handlers.GetSearchHistoryOptionsHandler", st)

	val, _ := c.Get("person")
	person := val.(entities.Person)
	searchID := utils.GetEntityIDFromParams(c, "search_id")
	req := entities.GetSearchHistoryOptionsReq{}
	req.SearchID = searchID
	res, derr := services.GetSearchHistoryOptions(&person, req)
	if derr != nil {
		derr.Respond(c)
		return
	}
	c.JSON(http.StatusOK, res)
}
