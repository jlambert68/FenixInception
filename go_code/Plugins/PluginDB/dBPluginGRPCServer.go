package PluginKeyValueDBStore

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

// ***************************************************************
// Start gateways gRPC-server for messages from Plugins to be read or saved from/to DB
//
func startDbPluginGRPCServerForMessagesFromPlugins() {

	var err error

	// Start gRPC serve
	logger.WithFields(logrus.Fields{
		"ID": "5829f8f0-5ceb-47b7-948b-c8b61936806a",
	}).Debug("KeyValueStore gRPC server for read/write requests from Plugins tries to start")

	logger.WithFields(logrus.Fields{
		"keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort": keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort,
		"ID": "200ef98f-503f-4d9c-acc2-6eb1a57477b7",
	}).Info("Trying to listening on port: " + strconv.FormatInt(int64(keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort), 10))

	keyValueStoreListener, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort), 10))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "07af7302-4e94-4bcf-bde3-f3a209b23428",
		}).Fatal("failed to listen:")

	} else {
		logger.WithFields(logrus.Fields{
			"keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort": keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort,
			"ID": "113dfc08-115c-43d3-a946-dddd866822ac",
		}).Info("Success in listening on port " + strconv.FormatInt(int64(keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort), 10))

	}

	// Creates a new registerKeyValueStoreServer gRPC server
	go func() {
		logger.WithFields(logrus.Fields{
			"ID": "9244f3dc-b4dc-4354-9a15-70d97d173c40",
		}).Info("Starting KeyValueStore gRPC Server")

		registerKeyValueStoreServer = grpc.NewServer()
		gRPC.RegisterPluginDBServer(registerKeyValueStoreServer, &gRPCServerForDbKeyValueStoreStruct{})

		logger.WithFields(logrus.Fields{
			"ID": "0a65f2da-bb18-42e6-8a51-a16ae2fd158f",
		}).Debug("registerKeyValueStoreServer for Gateway started")

		err := registerKeyValueStoreServer.Serve(keyValueStoreListener)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":  "014f2a52-8332-4249-b653-8051f6cef7a7",
				"err": err,
			}).Fatal("Couldn't serve 'registerKeyValueStoreServer.Serve(keyValueStoreListener)'")
		} else {
			logger.WithFields(logrus.Fields{
				"ID": "5324b052-c93e-4783-b025-2f4430f34fe2",
			}).Debug("Success in serving 'registerKeyValueStoreServer.Serve(keyValueStoreListener)'")
		}
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Plugins
//
func stopDbPluginGRPCServerForMessagesFromPlugins() {

	// Close gRPC connection to Children
	registerKeyValueStoreServer.GracefulStop()

	logger.WithFields(logrus.Fields{
		"ID": "1e3f84f9-326c-4584-801f-014a4a8d243f",
	}).Info("Gracefull stop for: 'registerKeyValueStoreServer'")

	// Stop listening towards clients
	err := keyValueStoreListener.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "54154530-b7fe-482a-a9fa-ea2fa63f6b5e",
		}).Error("Couldn't Stop listening, from Plugins-messages, on port: " + strconv.FormatInt(int64(keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort), 10))
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "20e82552-cebe-42bc-ba0b-a3689c479798",
		}).Info("Stop listening, from Children, on port: " + strconv.FormatInt(int64(keyValueStoreConfig.KeyValueStoreIdentification.PluginsCallOnThisPort), 10))
	}
}
