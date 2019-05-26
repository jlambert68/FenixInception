package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func (gatewayObject *gatewayTowardsPluginObject_struct) InitLogger(filename string) {
	gatewayObject.logger = logrus.StandardLogger()

	switch LoggingLevel {

	case logrus.DebugLevel:
		log.Println("'common_config.LoggingLevel': ", LoggingLevel)

	case logrus.InfoLevel:
		log.Println("'common_config.LoggingLevel': ", LoggingLevel)

	case logrus.WarnLevel:
		log.Println("'common_config.LoggingLevel': ", LoggingLevel)

	default:
		log.Println("Not correct value for debugging-level, this was used: ", common_config.LoggingLevel)
		os.Exit(0)

	}

	logrus.SetLevel(LoggingLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	//If no file then set standard out

	if filename == "" {
		gatewayObject.gatewayCommonObjects.logger.Out = os.Stdout

	} else {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			gatewayObject.gatewayCommonObjects.logger.Out = file
		} else {
			log.Println("Failed to log to file, using default stderr")
		}
	}

	// Should only be done from init functions
	//grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))

}
