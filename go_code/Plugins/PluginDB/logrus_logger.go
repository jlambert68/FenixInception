package PluginKeyValueDBStore

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func initLogger(filename string) {
	logger = logrus.StandardLogger()

	switch keyValueStoreConfig.LoggingLevel.LoggingLevel {

	case logrus.DebugLevel:
		log.Println("'common_config.LoggingLevel': ", keyValueStoreConfig.LoggingLevel.LoggingLevel)

	case logrus.InfoLevel:
		log.Println("'common_config.LoggingLevel': ", keyValueStoreConfig.LoggingLevel.LoggingLevel)

	case logrus.WarnLevel:
		log.Println("'common_config.LoggingLevel': ", keyValueStoreConfig.LoggingLevel.LoggingLevel)

	case logrus.FatalLevel:
		log.Println("'common_config.LoggingLevel': ", keyValueStoreConfig.LoggingLevel.LoggingLevel)

	default:
		log.Println("fea6e90c-c1cd-478a-b423-9008e2463608")
		log.Println("Not correct value for debugging-level, this was used: ", keyValueStoreConfig.LoggingLevel.LoggingLevel)
		log.Println("'keyValueStoreConfig': ", keyValueStoreConfig)
		os.Exit(0)

	}

	logrus.SetLevel(keyValueStoreConfig.LoggingLevel.LoggingLevel)
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
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //os.O_CREATE
		if err == nil {
			logger.Out = file

			logger.Info("********* *************** *********")
			logger.Info("")
			logger.WithFields(logrus.Fields{
				"ID": "d6b7454e-eb99-4c7d-9ec9-84249a7ee848",
			}).Info("********* logger started *********")
			logger.Info("")
			logger.Info("********* *************** *********")
		} else {
			log.Println("Failed to log to file, using default stderr")
		}
	}

	// Should only be done from init functions
	//grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))

}
