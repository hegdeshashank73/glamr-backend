package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, entities.HealthRes{
		Status:       "ok",
		IsProduction: utils.IsProduction(),
		Region:       utils.GetRegion(),
	})
}

func TestHandler(c *gin.Context) {
	data := map[string]string{}

	if _, err := vendors.DBMono.Exec("SELECT 1;"); err != nil {
		data["db_get"] = err.Error()
	} else {
		data["db_get"] = "1"
	}

	c.JSON(http.StatusOK, data)
}
