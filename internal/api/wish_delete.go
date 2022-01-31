package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tktip/wishtree-api/internal/db"
)

func (api *API) deleteWish(c *gin.Context) {
	wishID, ok := validateDeleteQuery(c)
	if !ok {
		return
	}

	wish, err := api.DB.GetWishByID(nil, wishID)

	if err != nil {
		logrus.Errorf("Failed to get all wishes >> %v", err)
		c.String(http.StatusInternalServerError, "Error when getting wish with given ID")
		return
	}
	if wish == nil {
		c.String(http.StatusBadRequest, "Wish with given ID not found")
		return
	}

	tx, err := api.DB.NewTransaction()
	if err != nil {
		logrus.Errorf("failed to create tx >> %v", err)
		c.String(http.StatusInternalServerError, "Error connecting to database")
		return
	}
	defer db.TxRollbackIfErr(tx, &err)

	err = db.ClearWishByID(*tx, wishID)
	if err != nil {
		logrus.Errorf("failed to delete wish >> %v", err)
		c.String(http.StatusInternalServerError, "Failed to delete wish")
		return
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("failed to commit createWish transaction >> %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, wish)
}

func validateDeleteQuery(c *gin.Context) (wishID int, ok bool) {
	var wishIDString string
	wishIDString = c.Param("id")
	if wishIDString == "" {
		c.String(http.StatusBadRequest, "Missing query param for wishID")
		return 0, false
	}

	wishID, err := strconv.Atoi(wishIDString)
	if err != nil {
		c.String(http.StatusBadRequest, "Query param for wish id must be a number")
		return 0, false
	}

	return wishID, true
}
