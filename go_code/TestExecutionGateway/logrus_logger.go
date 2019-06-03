package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func initLogger(filename string) {
	logger = logrus.StandardLogger()

	switch gatewayConfig.loggingLevel.loggingLevel {

	case logrus.DebugLevel:
		log.Println("'common_config.LoggingLevel': ", gatewayConfig.loggingLevel.loggingLevel)

	case logrus.InfoLevel:
		log.Println("'common_config.LoggingLevel': ", gatewayConfig.loggingLevel.loggingLevel)

	case logrus.WarnLevel:
		log.Println("'common_config.LoggingLevel': ", gatewayConfig.loggingLevel.loggingLevel)

	case logrus.FatalLevel:
		log.Println("'common_config.LoggingLevel': ", gatewayConfig.loggingLevel.loggingLevel)

	default:
		log.Println("Not correct value for debugging-level, this was used: ", gatewayConfig.loggingLevel.loggingLevel)
		os.Exit(0)

	}

	logrus.SetLevel(gatewayConfig.loggingLevel.loggingLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	//If no file then set standard out

	if filename == "" {
		logger.Out = os.Stdout

	} else {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logger.Out = file
		} else {
			log.Println("Failed to log to file, using default stderr")
		}
	}

	// Should only be done from init functions
	//grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))

}
