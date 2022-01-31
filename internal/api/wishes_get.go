package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (api *API) getAllWishes(c *gin.Context) {
	wishes, err := api.DB.GetAllWishes()

	if err != nil {
		logrus.Errorf("Failed to get all wishes >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to get wishes")
		return
	}

	c.JSON(http.StatusOK, wishes)
}
