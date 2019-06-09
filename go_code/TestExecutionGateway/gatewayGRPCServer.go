package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
	"strconv"
)

// ***************************************************************
// Start gateways gRPC-server for messages towards Plugins
//
func startGatewayGRPCServerForMessagesTowardsPlugins() {

	var err error

	// If no port from Parent Gateway/Fenix has been received then exit this function
	if gatewayConfig.parentgRPCAddress.connectionToParentDoneAtLeastOnce == false {
		// Gateway har never had a successful connection to parent gateway/Fenx
		logger.WithFields(logrus.Fields{
			"ID": "f6398d61-fd34-49f4-a6e5-18b2318428cb",
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
		}).Error("Couldn't Stop listening, from Children, on port: " + strconv.FormatInt(int64(gatewayConfig.gatewayIdentification.gatewayChildrenCallOnThisPort), 10))
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Info("Stop listening, from Children, on port: " + strconv.FormatInt(int64(gatewayConfig.gatewayIdentification.gatewayChildrenCallOnThisPort), 10))
	}
}

// ***************************************************************
// Start gateways gRPC-server for messages towards Fenix
//
func startGatewayGRPCServerForMessagesTowardsFenix() {

	var err error

	// If no port from Parent Gateway/Fenix has been received then exit this function
	if gatewayConfig.parentgRPCAddress.connectionToParentDoneAtLeastOnce == false {
		// Gateway har never had a successful connection to parent gateway/Fenx
		logger.WithFields(logrus.Fields{
			"ID": "92da53d1-345b-493c-bf01-023b882bcbc2",
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

		err := registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":  "00c5957f-c2f4-4aa7-a049-159c7fda06ee",
				"err": err,
			}).Fatal("Couldn't serve 'registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)'")
		} else {
			logger.WithFields(logrus.Fields{
				"ID": "93dfa984-0b9c-48cf-8fae-a5c5403bdd92",
			}).Debug("Success in serving 'registerGatewayTowardsPluginerver.Serve(gatewayTowardsPluginListener)'")
		}
	}()

}

// ***************************************************************
// Stop gateways gRPC-server for messages towards Fenix
//
func stopGatewayGRPCServerForMessagesTowardsFenix() {

	// Close gRPC connection to Parent gateway/Fenix
	registerGatewayTowardsFenixServer.GracefulStop()

	logger.WithFields(logrus.Fields{
		"ID": "ed85f1c6-98d9-4c08-bf52-2364a42a5bad",
	}).Info("Gracefull stop for: 'registerGatewayTowardsFenixServer'")

	// Stop listening towards parent gateway/Fenix
	err := gatewayTowardsFenixListener.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "3d18cec5-6d65-4a71-87ac-e43ca88e459f",
		}).Error("Couldn't Stop listening, from Parent, on port: " + strconv.FormatInt(int64(gatewayConfig.gatewayIdentification.gatewaParentCallOnThisPort), 10))
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "87ea079f-b7a2-444e-b1f2-adbfebe3d804",
		}).Info("Stop listening, from Parent, on port: " + strconv.FormatInt(int64(gatewayConfig.gatewayIdentification.gatewaParentCallOnThisPort), 10))
	}
}
