//go:build windows
// +build windows

package main

import (
	"os"
	"strings"

	"github.com/haraldfw/cfger"
	"github.com/sirupsen/logrus"
	"github.com/tktip/wishtree-api/internal/api"
	"github.com/tktip/wishtree-api/internal/version"
)

func main() {
	cfgFile := ""
	if len(os.Args) > 1 {
		for _, v := range os.Args {
			vals := strings.Split(v, "=")
			if len(vals) == 1 {
				continue
			}
			if vals[0] == "CONFIG" {
				cfgFile = vals[1]
			}
		}
	}
	if cfgFile == "" {
		logrus.Fatal("Missing parameter 'CONFIG'")
	}

	logrus.Infof("Running version %s", version.VERSION)

	apiConfig := api.API{}
	_, err := cfger.ReadStructuredCfg("file::"+cfgFile, &apiConfig)
	if err != nil {
		logrus.Fatalf("Could not read config: %v", err)
	}

	logrus.Fatalf("Api failed: %v", apiConfig.Run())
}
