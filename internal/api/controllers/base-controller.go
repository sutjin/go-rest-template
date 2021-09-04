package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	variable "sutjin/go-rest-template/internal/pkg/config"
	models "sutjin/go-rest-template/internal/pkg/models"
)

func GetVersion(c *gin.Context) {
	info := models.AppInfo{
		Version:  variable.Version,
		Deployed: variable.BuildTime,
	}
	c.JSON(http.StatusOK, info)
}
