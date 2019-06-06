package TestExecutionGateway

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

func processConfigFile() {
	if _, err := toml.DecodeFile("gatewayConfig.toml", &gatewayConfig); err != nil {
		logger.WithFields(logrus.Fields{
			"ID":    "9f0c6e06-5211-4d48-829e-49f2e9382b14",
			"error": err,
		}).Fatal("Couldn't decode 'gatewayConfig.toml', stopping gateway")
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "8c165291-40f2-4c99-a709-be456bbaf96e",
		}).Debug(" 'gatewayConfig.toml' was decoded")
	}
}
