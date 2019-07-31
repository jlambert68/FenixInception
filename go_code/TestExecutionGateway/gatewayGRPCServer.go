package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"jlambert/FenixInception2/go_code/common_code"
	"net"
	"strconv"
)

// ***************************************************************
// Start gateways gRPC-server for messages towards Plugins
//
func startGatewayGRPCServerForMessagesTowardsPlugins() {

	var err error

	// Start gRPC serve
	common_code.logger.WithFields(logrus.Fields{
		"ID": "225c2b2e-1891-4175-8950-a1c5b721be17",
	}).Debug("Gateway gRPC server towards plugin tries to start")

	common_code.logger.WithFields(logrus.Fields{
		"gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort": common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort,
		"ID": "e6eacb4e-1ac4-4b75-87bb-e9bb05e728e1",
	}).Info("Trying to listening on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	common_code.gatewayTowardsPluginListener, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))
	if err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "2bb724a4-efbf-4d86-b120-735e82d9b920",
		}).Fatal("failed to listen:")

	} else {
		common_code.logger.WithFields(logrus.Fields{
			"gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort": common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort,
			"ID": "c36912f1-54ce-458a-b2c3-f64209692ede",
		}).Info("Success in listening on port " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))

	}

	// Creates a new registerGatewayTowardsPluginerver gRPC server
	go func() {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "025e810f-90ed-4418-ad99-93734e4d0e3f",
		}).Info("Starting Gateway gRPC Server towards Plugin")

		common_code.registerGatewayTowardsPluginerver = grpc.NewServer()
		gRPC.RegisterGatewayTowayPluginServer(common_code.registerGatewayTowardsPluginerver, &common_code.GRPCServerTowardsPluginStruct{})

		common_code.logger.WithFields(logrus.Fields{
			"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
		}).Debug("registerGatewayTowardsPluginerver for Gateway started")

		err := common_code.registerGatewayTowardsPluginerver.Serve(common_code.gatewayTowardsPluginListener)
		if err != nil {
			common_code.logger.WithFields(logrus.Fields{
				"ID":  "33aff438-cd24-462c-b609-003459b3fecd",
				"err": err,
			}).Fatal("Couldn't serve 'registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)'")
		} else {
			common_code.logger.WithFields(logrus.Fields{
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
	common_code.registerGatewayTowardsPluginerver.GracefulStop()

	common_code.logger.WithFields(logrus.Fields{
		"ID": "a7ee0eea-3bf2-44b7-a5e1-942a15d0d7fd",
	}).Info("Gracefull stop for: 'registerGatewayTowardsPluginerver'")

	// Stop listening towards clients
	err := common_code.gatewayTowardsFenixListener.Close()
	if err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Error("Couldn't Stop listening, from Children, on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	} else {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Info("Stop listening, from Children, on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	}
}

// ***************************************************************
// Start gateways gRPC-server for messages towards Fenix
//
func startGatewayGRPCServerForMessagesTowardsFenix() {

	var err error

	// Start gRPC serve
	common_code.logger.WithFields(logrus.Fields{
		"ID": "57e5c579-975d-46a6-a92b-027bbbefff1e",
	}).Debug("Gateway gRPC server towards Fenix tries to start")

	// If no port from Parent Gateway/Fenix has been received then exit this function
	if common_code.gatewayConfig.ParentgRPCAddress.ConnectionToParentDoneAtLeastOnce == false && common_code.gatewayConfig.IntegrationTest.StartWithOutAnyParent == false {
		// Gateway har never had a successful connection to parent gateway/Fenx
		common_code.logger.WithFields(logrus.Fields{
			"ID": "92da53d1-345b-493c-bf01-023b882bcbc2",
		}).Fatal("This Gateway has never been registrated at perent gateway/Fenix, exiting")
		return
	}

	common_code.logger.WithFields(logrus.Fields{
		"gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort": common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort,
		"ID": "e6eacb4e-1ac4-4b75-87bb-e9bb05e728e1",
	}).Info("Trying to listening on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))

	common_code.gatewayTowardsFenixListener, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))
	if err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID":    "2bb724a4-efbf-4d86-b120-735e82d9b920",
		}).Fatal("failed to listen: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))

	} else {
		common_code.logger.WithFields(logrus.Fields{
			"gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort": common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort,
			"ID": "5f86c32d-f143-4853-b35a-35d58548bb7d",
		}).Info("Success in listening on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort), 10))

	}

	// Creates a new registerGatewayTowardsFenixServer gRPC server
	go func() {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "8b105bd8-3825-40e5-ab5e-47e78b8a6f86",
		}).Info("Starting Gateway gRPC Server towards Fenix")

		common_code.registerGatewayTowardsFenixServer = grpc.NewServer()
		gRPC.RegisterGatewayTowardsFenixServer(common_code.registerGatewayTowardsFenixServer, &common_code.GRPCServerTowardsFenixStruct{})

		common_code.logger.WithFields(logrus.Fields{
			"ID": "0bb9fe6a-aff6-4a4e-827e-63b4dbd8d85d",
		}).Debug("registerGatewayTowardsFenixServer for Gateway started")

		err := common_code.registerGatewayTowardsFenixServer.Serve(common_code.gatewayTowardsFenixListener)
		if err != nil {
			common_code.logger.WithFields(logrus.Fields{
				"ID":  "00c5957f-c2f4-4aa7-a049-159c7fda06ee",
				"err": err,
			}).Fatal("Couldn't serve 'registerGatewayTowardsFenixServer.Serve(gatewayTowardsFenixListener)'")
		} else {
			common_code.logger.WithFields(logrus.Fields{
				"ID": "93dfa984-0b9c-48cf-8fae-a5c5403bdd92",
			}).Debug("Success in serving 'registerGatewayTowardsFenixServer.Serve(gatewayTowardsFenixListener)'")
		}
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Fenix
//
func stopGatewayGRPCServerForMessagesTowardsFenix() {

	// Close gRPC connection to Parent gateway/Fenix
	common_code.registerGatewayTowardsFenixServer.GracefulStop()

	common_code.logger.WithFields(logrus.Fields{
		"ID": "ed85f1c6-98d9-4c08-bf52-2364a42a5bad",
	}).Info("Gracefull stop for: 'registerGatewayTowardsFenixServer'")

	// Stop listening towards parent gateway/Fenix
	err := common_code.gatewayTowardsFenixListener.Close()
	if err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "3d18cec5-6d65-4a71-87ac-e43ca88e459f",
		}).Error("Couldn't Stop listening, from Parent, on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))
	} else {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "87ea079f-b7a2-444e-b1f2-adbfebe3d804",
		}).Info("Stop listening, from Parent, on port: " + strconv.FormatInt(int64(common_code.gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort), 10))
	}
}
