package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (api *API) getAllCategories(c *gin.Context) {
	categories, err := api.DB.GetAllCategories()

	if err != nil {
		logrus.Errorf("Failed to get all categories >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to get categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}
