package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
)

// ***************************************************************
// Start gateways gRPC-server for messages towards Plugins
//
func (gatewayObject *gatewayTowardsPluginObject_struct) startGatewayGRPCServerForMessagesTowardsPlugins() {

	var err error

	// If no port from Parent Gateway/Fenix has been received then exit this function
	if incomingPortForCallsFromParentGateway == "" {
		// Find first non allocated port from defined start port
		logger.WithFields(logrus.Fields{
			"ID": "737867a7-ecd5-4bea-8249-7c1adc2dc130",
		}).Fatal("This Gateway has never been registrated at perent gateway/Fenix, exiting")
		return
	}

	// Start gRPC serve
	logger.WithFields(logrus.Fields{
		"ID": "225c2b2e-1891-4175-8950-a1c5b721be17",
	}).Debug("Gateway gRPC server tries to start")

	logger.WithFields(logrus.Fields{
		"incomingPortForCallsFromParentGateway": incomingPortForCallsFromParentGateway,
		"ID":                                    "e6eacb4e-1ac4-4b75-87bb-e9bb05e728e1",
	}).Info("Trying to listening on port: " + incomingPortForCallsFromParentGateway)

	gatewayTowardsPluginListener, err = net.Listen("tcp", string(incomingPortForCallsFromParentGateway))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "2bb724a4-efbf-4d86-b120-735e82d9b920",
		}).Fatal("failed to listen:")

	} else {
		logger.WithFields(logrus.Fields{
			"incomingPortForCallsFromParentGateway": incomingPortForCallsFromParentGateway,
			"ID":                                    "77dd611b-4b76-4ed8-a383-8b336a7b29ec",
		}).Info("Success in listening on port " + incomingPortForCallsFromParentGateway)

	}

	// Creates a new registerGatewayTowardsPluginerver gRPC server
	go func() {
		logger.WithFields(logrus.Fields{
			"ID": "3361684d-73f9-49b2-a11a-a2833759363c",
		}).Info("Starting Gateway gRPC Server")

		registerGatewayTowardsPluginerver = grpc.NewServer()
		gRPC.RegisterGatewayTowayPluginServer(registerGatewayTowardsPluginerver, &gatewayTowardsPluginObject_struct{})

		logger.WithFields(logrus.Fields{
			"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
		}).Debug("registerGatewayTowardsPluginerver for Gateway started")

		registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Plugins
//
func stopGatewayGRPCServerForMessagesTowardsPlugins() {
}

// ***************************************************************
// Start gateways gRPC-server for messages towards Fenix
//
func (gatewayObject *gatewayTowardsFenixObject_struct) startGatewayGRPCServerForMessagesTowardsFenix() {

	var err error

	// If no port from Parent Gateway/Fenix has been received then exit this function
	if incomingPortForCallsFromParentGateway == "" {
		// Find first non allocated port from defined start port
		logger.WithFields(logrus.Fields{
			"ID": "737867a7-ecd5-4bea-8249-7c1adc2dc130",
		}).Fatal("This Gateway has never been registrated at perent gateway/Fenix, exiting")
		return
	}

	// Start gRPC serve
	logger.WithFields(logrus.Fields{
		"ID": "225c2b2e-1891-4175-8950-a1c5b721be17",
	}).Debug("Gateway gRPC server tries to start")

	logger.WithFields(logrus.Fields{
		"incomingPortForCallsFromParentGateway": incomingPortForCallsFromParentGateway,
		"ID":                                    "e6eacb4e-1ac4-4b75-87bb-e9bb05e728e1",
	}).Info("Trying to listening on port: " + incomingPortForCallsFromParentGateway)

	gatewayTowardsPluginListener, err = net.Listen("tcp", string(incomingPortForCallsFromParentGateway))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "2bb724a4-efbf-4d86-b120-735e82d9b920",
		}).Fatal("failed to listen:")

	} else {
		logger.WithFields(logrus.Fields{
			"incomingPortForCallsFromParentGateway": incomingPortForCallsFromParentGateway,
			"ID":                                    "77dd611b-4b76-4ed8-a383-8b336a7b29ec",
		}).Info("Success in listening on port " + incomingPortForCallsFromParentGateway)

	}

	// Creates a new registerGatewayTowardsPluginerver gRPC server
	go func() {
		logger.WithFields(logrus.Fields{
			"ID": "3361684d-73f9-49b2-a11a-a2833759363c",
		}).Info("Starting Gateway gRPC Server")

		registerGatewayTowardsPluginerver = grpc.NewServer()
		gRPC.RegisterGatewayTowayPluginServer(registerGatewayTowardsPluginerver, &gatewayTowardsPluginObject_struct{})

		logger.WithFields(logrus.Fields{
			"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
		}).Debug("registerGatewayTowardsPluginerver for Gateway started")

		registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Fenix
//
func stopGatewayGRPCServerForMessagesTowardsFenix() {
}
