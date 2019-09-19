package PluginKeyValueDBStore

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

// ***************************************************************
// Start gateways gRPC-server for messages towards Plugins
//
func startDbPluginGRPCServerForMessagesPlugins() {

	var err error

	// Start gRPC serve
	logger.WithFields(logrus.Fields{
		"ID": "225c2b2e-1891-4175-8950-a1c5b721be17",
	}).Debug("Gateway gRPC server towards plugin tries to start")

	logger.WithFields(logrus.Fields{
		"gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort": gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort,
		"ID": "e6eacb4e-1ac4-4b75-87bb-e9bb05e728e1",
	}).Info("Trying to listening on port: " + strconv.FormatInt(int64(gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	gatewayTowardsPluginListener, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "2bb724a4-efbf-4d86-b120-735e82d9b920",
		}).Fatal("failed to listen:")

	} else {
		logger.WithFields(logrus.Fields{
			"gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort": gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort,
			"ID": "c36912f1-54ce-458a-b2c3-f64209692ede",
		}).Info("Success in listening on port " + strconv.FormatInt(int64(gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	}

	// Creates a new registerGatewayTowardsPluginerver gRPC server
	go func() {
		logger.WithFields(logrus.Fields{
			"ID": "025e810f-90ed-4418-ad99-93734e4d0e3f",
		}).Info("Starting Gateway gRPC Server towards Plugin")

		registerGatewayTowardsPluginerver = grpc.NewServer()
		gRPC.RegisterGatewayTowayPluginServer(registerGatewayTowardsPluginerver, &common_code.GRPCServerTowardsPluginStruct{})

		logger.WithFields(logrus.Fields{
			"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
		}).Debug("registerGatewayTowardsPluginerver for Gateway started")

		err := registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":  "33aff438-cd24-462c-b609-003459b3fecd",
				"err": err,
			}).Fatal("Couldn't serve 'registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)'")
		} else {
			logger.WithFields(logrus.Fields{
				"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
			}).Debug("Success in serving 'registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)'")
		}
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Plugins
//
func stopGatewayGRPCServerForMessagesTowardsPlugins() {

	// Close gRPC connection to Children
	registerGatewayTowardsPluginerver.GracefulStop()

	logger.WithFields(logrus.Fields{
		"ID": "a7ee0eea-3bf2-44b7-a5e1-942a15d0d7fd",
	}).Info("Gracefull stop for: 'registerGatewayTowardsPluginerver'")

	// Stop listening towards clients
	err := gatewayTowardsFenixListener.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Error("Couldn't Stop listening, from Children, on port: " + strconv.FormatInt(int64(gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Info("Stop listening, from Children, on port: " + strconv.FormatInt(int64(gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	}
}
