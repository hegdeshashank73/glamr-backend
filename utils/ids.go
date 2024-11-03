package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hegdeshashank73/glamr-backend/common"
)

func GenerateAccessToken() string {
	return uuid.New().String()
}

func GetEntityIDFromParams(c *gin.Context, paramName string) common.Snowflake {
	_id, isPresent := c.Params.Get(paramName)
	if isPresent {
		return common.BuildSnowflake(_id)
	}
	return common.Snowflake(0)
}

func GetEntityIDFromQueryParams(c *gin.Context, paramName string) common.Snowflake {
	_id := c.Query(paramName)
	if _id == "" {
		return common.Snowflake(0)
	}
	return common.BuildSnowflake(_id)
}
