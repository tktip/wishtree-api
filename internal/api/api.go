package api

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/tktip/cfger"
	"github.com/tktip/wishtree-api/internal/db"
	v "gopkg.in/go-playground/validator.v9"
)

var (
	validator = v.New()
	// There are 1500 wishes in total. Max num 1400 = always keep 100 free.
	maxNumberOfWishes = 1400
)

// AdminUser is a username and password for an admin user
type AdminUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// TokenConfig contains information needed about the JWT token generated
type TokenConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	TokenName     string `yaml:"tokenName"`
	DurationHours int    `yaml:"durationHours"`
	Path          string `yaml:"path"`
	Domain        string `yaml:"domain"`
	Secure        bool   `yaml:"secure"`
	HTTPOnly      bool   `yaml:"httpOnly"`
}

// API - the api
type API struct {
	LogLevel            string        `yaml:"logLevel"`
	DB                  *db.Connector `yaml:"db"`
	AdminUsers          []AdminUser   `yaml:"adminUsers"`
	TokenConfig         TokenConfig   `yaml:"tokenConfig"`
	TokenPrivateKeyPath string        `yaml:"tokenPrivatekeyPath"`
	TokenPublicKeyPath  string        `yaml:"tokenPublickeyPath"`
	HTTPClient          *http.Client
}

// Run - start the server
func (api *API) Run() error {
	err := api.init()
	if err != nil {
		return err
	}

	g := gin.New()
	g.Use(gin.Logger())

	openGroup := g.Group("/")

	openGroup.GET("/wishes", api.getAllWishes)
	openGroup.GET("/categories", api.getAllCategories)
	openGroup.POST("/wishes", api.createWish)
	openGroup.POST("/login", api.login)
	openGroup.GET("/tree-status", api.getTreeStatus)

	adminGroup := g.Group("/")
	adminGroup.Use(api.adminRoutesMiddleware)

	adminGroup.DELETE("/wishes/:id", api.deleteWish)
	adminGroup.GET("/logout", api.logout)
	adminGroup.POST("/tree-status", api.updateTreeStatus)

	return g.Run()
}

func (api *API) init() error {
	if api.LogLevel != "" {
		level, err := logrus.ParseLevel(api.LogLevel)
		if err != nil {
			return err
		}

		logrus.SetLevel(level)
		logrus.Infof("Running with log level '%s'", api.LogLevel)
	}

	_, err := cfger.ReadStructuredCfg("env::CONFIG", &api)
	if err != nil {
		logrus.Fatalf("Could not read config:  %v", err)
	}

	err = validator.Struct(api)
	if err != nil {
		return err
	}

	api.HTTPClient = &http.Client{
		Timeout: time.Second * 15,
	}

	err = api.setupJWTKeys()
	if err != nil {
		return err
	}

	return nil
}

func (api *API) setupJWTKeys() (err error) {
	var privateKeyRaw, publicKeyRaw []byte

	privateKeyRaw, err = ioutil.ReadFile(api.TokenPrivateKeyPath)
	if err != nil {
		err = fmt.Errorf("error reading private key file >> %w", err)
		return
	}
	api.TokenConfig.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		err = fmt.Errorf("error parsing private key >> %w", err)
		return
	}

	publicKeyRaw, err = ioutil.ReadFile(api.TokenPublicKeyPath)
	if err != nil {
		err = fmt.Errorf("error reading public key file >> %w", err)
		return
	}
	api.TokenConfig.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyRaw)
	if err != nil {
		err = fmt.Errorf("error parsing public key >> %w", err)
		return
	}

	return
}
