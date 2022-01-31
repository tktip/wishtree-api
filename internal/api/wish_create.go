package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tktip/wishtree-api/internal/db"
	"github.com/tktip/wishtree-api/pkg/apiutils"
)

type createWishBody struct {
	WishID     *int    `form:"wishId"`
	Text       string  `form:"text" validate:"required"`
	CategoryID int     `form:"categoryId" validate:"required"`
	Author     *string `form:"author"`
	ZipCode    string  `form:"zipCode" validate:"required"`
}

//revive:disable-next-line:cyclomatic - ~All the cyclo is err handling
func (api *API) createWish(c *gin.Context) {
	var body createWishBody
	if !body.bindAndValidate(c) {
		return
	}

	isOpen, err := api.DB.GetIsTreeOpen()
	if !isOpen {
		c.String(http.StatusForbidden, "Not accepting new wishes")
		return
	}

	tx, err := api.DB.NewTransaction()
	if err != nil {
		logrus.Errorf("failed to create tx >> %v", err)
		c.String(http.StatusInternalServerError, "Error connecting to database")
		return
	}
	defer db.TxRollbackIfErr(tx, &err)

	var numberOfWishes int
	numberOfWishes, err = db.GetNumberOfTakenWishes(*tx)
	if err != nil {
		logrus.Errorf("failed to get number of taken wishes >> %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if numberOfWishes >= maxNumberOfWishes {
		var archivedWishID int
		archivedWishID, err = db.ArchiveAndClearOldestActiveWish(*tx)
		if err != nil {
			logrus.Errorf("failed to archive oldest active wish >> %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		logrus.Infof("Cleared and archived copy of oldest wish, ID %v", archivedWishID)
	}

	var wishID int
	if body.WishID == nil {
		var freeWish db.Wish
		freeWish, err = db.GetRandomFreeWish(tx)
		if err != nil {
			logrus.Errorf("failed to get random free wish >> %v", err)
			c.String(http.StatusInternalServerError, "Error creating wish")
			return
		}
		wishID = freeWish.ID
	} else {
		wishID = *body.WishID
	}

	var newWish db.Wish
	newWish, err = api.DB.CreateWish(
		tx, wishID, body.Text, body.Author, body.ZipCode, body.CategoryID)

	if err == apiutils.ErrBadWishQuery {
		logrus.Errorf("Could not get wish with given ID >> %v", err)
		c.String(http.StatusBadRequest, "Could not get wish with given ID")
		return
	} else if err == apiutils.ErrWishIDTaken {
		c.String(http.StatusConflict, "Wish already taken")
		return
	} else if err != nil {
		logrus.Errorf("Failed to create wish >> %v", err)
		c.String(http.StatusInternalServerError, "Error creating wish")
		return
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Failed to commit createWish transaction >> %v", err)
		c.String(http.StatusInternalServerError, "Error creating wish")
		return
	}

	c.JSON(http.StatusOK, newWish)
}

func (body *createWishBody) bindAndValidate(c *gin.Context) (isOk bool) {
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
