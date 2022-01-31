//go:build linux
// +build linux

package main

import (
	"github.com/haraldfw/cfger"
	"github.com/sirupsen/logrus"
	"github.com/tktip/wishtree-api/internal/api"
	"github.com/tktip/wishtree-api/internal/version"
)

func main() {
	logrus.Infof("Running version %s", version.VERSION)

	apiConfig := api.API{}
	_, err := cfger.ReadStructuredCfg("env::CONFIG", &apiConfig)
	if err != nil {
		logrus.Fatalf("Could not read config: %v", err)
	}

	logrus.Fatalf("Api failed: %v", apiConfig.Run())
}
