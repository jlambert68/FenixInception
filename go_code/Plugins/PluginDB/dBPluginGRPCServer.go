package PluginKeyValueDBStore

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
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
		"keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort": keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort,
		"ID": "200ef98f-503f-4d9c-acc2-6eb1a57477b7",
	}).Info("Trying to listening on port: " + strconv.FormatInt(int64(keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	keyValueStoreListener, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "07af7302-4e94-4bcf-bde3-f3a209b23428",
		}).Fatal("failed to listen:")

	} else {
		logger.WithFields(logrus.Fields{
			"keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort": keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort,
			"ID": "113dfc08-115c-43d3-a946-dddd866822ac",
		}).Info("Success in listening on port " + strconv.FormatInt(int64(keyValueStoreConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	}

	// Creates a new registerKeyValueStoreServer gRPC server
	go func() {
		logger.WithFields(logrus.Fields{
			"ID": "9244f3dc-b4dc-4354-9a15-70d97d173c40",
		}).Info("Starting KeyValueStore gRPC Server")

		registerKeyValueStoreServer = grpc.NewServer()
		gRPC.RegisterGatewayTowayPluginServer(registerKeyValueStoreServer, &gRPCServerForDbKeyValueStoreStruct{})

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
func stopGatewayGRPCServerForMessagesTowardsPlugins() {

	// Close gRPC connection to Children
	registerKeyValueStoreServer.GracefulStop()

	logger.WithFields(logrus.Fields{
		"ID": "a7ee0eea-3bf2-44b7-a5e1-942a15d0d7fd",
	}).Info("Gracefull stop for: 'registerKeyValueStoreServer'")

	// Stop listening towards clients
	err := gatewayTowardsFenixListener.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Error("Couldn't Stop listening, from Children, on port: " + strconv.FormatInt(int64(keyValueStoreConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Info("Stop listening, from Children, on port: " + strconv.FormatInt(int64(keyValueStoreConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	}
}
