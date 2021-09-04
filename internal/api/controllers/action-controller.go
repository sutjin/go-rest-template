package controllers

import (
	"net/http"

	pg "sutjin/go-rest-template/internal/pkg/db"
	models "sutjin/go-rest-template/internal/pkg/models"
	http_err "sutjin/go-rest-template/pkg/http-err"
	"sutjin/go-rest-template/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetActions godoc
// @Summary Get Available Actions
// @Description Get available action from db
// @Tags sample
// @ID get-action
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Action
// @Failure 400 {object} http_err.HTTPError
// @Router /actions [get]
func GetActions(c *gin.Context) {
	var actions = []models.Action{}

	err := pg.DB.Model(&actions).Select()
	if err != nil {
		http_err.NewError(c, http.StatusBadRequest, err)
		logger.InlineLog(c, err.Error(), "database error")
		return
	}

	c.JSON(http.StatusOK, actions)
}

// PostAction godoc
// @Summary Post new action
// @Description Post new action to DB
// @Tags sample
// @ID post-action
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Action
// @Failure 400 {object} http_err.HTTPError
// @Router /actions [post]
func PostAction(c *gin.Context) {
	input := &models.Action{}
	err := c.BindJSON(input)

	if err != nil {
		http_err.NewError(c, http.StatusBadRequest, err)
		logger.InlineLog(c, err.Error(), input)
		return
	}

	_, err = pg.DB.Model(input).Insert()

	if err != nil {
		http_err.NewError(c, http.StatusBadRequest, err)
		logger.InlineLog(c, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, input)
}
