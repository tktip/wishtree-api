package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// LoginBody is the body for a login request
type LoginBody struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (api *API) adminRoutesMiddleware(c *gin.Context) {
	cookie, err := c.Cookie(api.TokenConfig.TokenName)

	if err == http.ErrNoCookie {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if err != nil {
		logrus.Errorf("Failed to get cookie >> %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return api.TokenConfig.PublicKey, nil
	})

	if err != nil {
		logrus.Errorf("error parsing token >> %v", err)
		c.AbortWithStatusJSON(http.StatusForbidden, "Invalid token signing method")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenExpUnix, ok := claims["exp"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, "Invalid token expiry")
			return
		}
		tokenExpTime := time.Unix(int64(tokenExpUnix), 0)

		if tokenExpTime.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusForbidden, "Token expired")
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, "Invalid token")
		return
	}

	c.Set("admin", true)
	c.Next()
}

func (api *API) login(c *gin.Context) {
	var body LoginBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Warnf("login error occurred binding request body >> %v", err)
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}

	isValidCredentials := false
	for _, user := range api.AdminUsers {
		if strings.ToLower(user.Username) == strings.ToLower(body.Username) &&
			user.Password == body.Password {
			isValidCredentials = true
			break
		}
	}

	if !isValidCredentials {
		c.String(http.StatusForbidden, "Feil brukernavn/passord")
		return
	}

	tokenDurationSeconds := api.TokenConfig.DurationHours * 3600
	token, err := api.generateToken(body.Username, tokenDurationSeconds)

	if err != nil {
		logrus.Errorf("error generating token >> %v", err)
		c.String(http.StatusInternalServerError, "Error generating token")
		return
	}

	c.SetCookie(
		api.TokenConfig.TokenName,
		token,
		api.TokenConfig.DurationHours*3600,
		api.TokenConfig.Path,
		api.TokenConfig.Domain,
		api.TokenConfig.Secure,
		api.TokenConfig.HTTPOnly,
	)
	c.JSON(http.StatusOK, gin.H{
		"username":   body.Username,
		"ttlSeconds": tokenDurationSeconds,
	})
}

func (api *API) generateToken(username string, tokenDurationSeconds int) (
	signedToken string, err error,
) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Duration(tokenDurationSeconds) * time.Second).Unix(),
	})

	signedToken, err = token.SignedString(api.TokenConfig.PrivateKey)

	return
}

func (api *API) logout(c *gin.Context) {
	c.SetCookie(
		api.TokenConfig.TokenName,
		"",
		-10,
		api.TokenConfig.Path,
		api.TokenConfig.Domain,
		api.TokenConfig.Secure,
		api.TokenConfig.HTTPOnly,
	)
	c.Status(http.StatusNoContent)
}

func isAdmin(c *gin.Context) bool {
	return c.GetBool("admin")
}
