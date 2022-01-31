package apiutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tktip/wishtree-api/pkg/contextlogging"
)

type errorResponse struct {
	ErrorMsg         string
	ErrorID          int
	HTTPResponseCode int
}

func (errorRes errorResponse) GinAndLogErr(
	err error, c *gin.Context,
) {
	log := contextlogging.GetGinLogContext(c)
	log.Errorf("%v", err)
	c.JSON(errorRes.HTTPResponseCode, gin.H{
		"error":   errorRes.ErrorMsg,
		"errorId": errorRes.ErrorID,
		"context": contextlogging.GetGinLogContext(c).ID,
	})
}

func (errorRes errorResponse) GinAndLogWarn(
	warningMessage string, c *gin.Context,
) {
	log := contextlogging.GetGinLogContext(c)
	log.Warnf(warningMessage)
	c.JSON(errorRes.HTTPResponseCode, gin.H{
		"error":   errorRes.ErrorMsg,
		"errorId": errorRes.ErrorID,
		"context": contextlogging.GetGinLogContext(c).ID,
	})
}

func (errorRes errorResponse) GinWithoutLogging(c *gin.Context) {
	c.JSON(errorRes.HTTPResponseCode, gin.H{
		"error":   errorRes.ErrorMsg,
		"errorId": errorRes.ErrorID,
		"context": contextlogging.GetGinLogContext(c).ID,
	})
}

// FailedBindingRequestBody defines the error response when the body's bind<___> threw an error
var FailedBindingRequestBody = errorResponse{
	ErrorMsg:         "Failed binding request body",
	ErrorID:          400,
	HTTPResponseCode: http.StatusBadRequest,
}

// MissingRequiredFieldsInRequestBody defines the error response when there's one or more
// missing required value from request body
var MissingRequiredFieldsInRequestBody = errorResponse{
	ErrorMsg: "Missing one or more required fields in request body, " +
		"a field has wrong data type, or a value is outside the valid enum values.",
	ErrorID:          4000,
	HTTPResponseCode: http.StatusBadRequest,
}
