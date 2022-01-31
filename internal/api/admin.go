package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tktip/wishtree-api/pkg/apiutils"
)

type updateTreeBody struct {
	IsOpen bool `form:"isOpen"`
}

func (body *updateTreeBody) bindAndValidate(c *gin.Context) (isOk bool) {
	err := c.BindJSON(&body)
	if err != nil {
		err = fmt.Errorf("error occurred binding request body >> %v", err)
		apiutils.FailedBindingRequestBody.GinAndLogWarn(err.Error(), c)
		return false
	}

	err = validator.Struct(body)
	if err != nil {
		err = fmt.Errorf("error occurred validating request body >> %v", err)
		apiutils.MissingRequiredFieldsInRequestBody.GinAndLogWarn(err.Error(), c)
		return false
	}

	return true
}

func (api *API) getTreeStatus(c *gin.Context) {
	isOpen, err := api.DB.GetIsTreeOpen()
	if err != nil {
		logrus.Errorf("Failed to get tree status >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to get tree status")
		return
	}

	wishCounts, err := api.DB.GetAllTreeWishCounts()
	if err != nil {
		logrus.Errorf("Failed to get number of wishes >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to get number of wishes")
		return
	}

	c.JSON(http.StatusOK, gin.H{"isOpen": isOpen, "wishCounts": wishCounts})
}

func (api *API) updateTreeStatus(c *gin.Context) {
	var body updateTreeBody
	if !body.bindAndValidate(c) {
		return
	}

	err := api.DB.UpdateTreeStatus(body.IsOpen)

	if err != nil {
		logrus.Errorf("Failed to update tree status >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to update tree status")
		return
	}

	c.JSON(http.StatusOK, gin.H{"isOpen": body.IsOpen})
}
