package FenixGatewayServer

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"jlambert/FenixInception2/go_code/common_code"
	gRPC "jlambert/FenixInception2/go_code/common_code/Gateway_gRPC_api"
	"time"
)

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func sendTestInstructionTowardsPlugin(testInstructionId string) (err error) {
	// Get Payload for TestInstruction to be sent to Plugin
	testInstructionPayload, err := getTestInstructionPayloadToPlugin(testInstructionId)

	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		err = dispatchEngineForTestInstructions(testInstructionPayload)

	}

	return err

}

// ********************************************************************************************
// Send TestInstructions to Fenix own Gateway for forwarding to Plugin
//

func dispatchEngineForTestInstructions(testInstructionToBeForwarded *gRPC.TestInstruction_RT) (err error) {

	// Service can be started and Stopped by central control of the Gateway
	if gatewayMustStopProcessing == true {
		// Service should be stopped from processing any messages
		time.Sleep(common_code.ServieSleepTime * time.Second)

	} else {
		// Process messages

		logger.WithFields(logrus.Fields{
			"ID":                           "ffd854c0-f570-4d1a-b945-f479df8d3be1",
			"testInstructionToBeForwarded": testInstructionToBeForwarded,
		}).Debug("Received a new TestInstruction that shoud be forwarded")

		// Send TestInstruction to client using gRPC-call
		addressToDial := getClientAddressAndPort(testInstructionToBeForwarded.PluginId)

		// Set up connection to Client Gateway or Plugin
		remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "73a85fbe-bf9d-4b15-bdf1-d56aca0f7754",
				"addressToDial": addressToDial,
				"error message": err,
			}).Error("Did not connect to Fenix gateway Server!")

		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "0110fa80-c543-4d5c-a02d-ea23fd8969f6",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to child-gateway- or Plugin-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			testInstructionToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
			testInstructionToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

			// Do gRPC-call to client gateway or Plugin
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendTestInstructionTowardsPlugin(ctx, testInstructionToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "69c6d6e9-f15e-4e73-a176-8ba947bf9f92",
					"returnMessage": returnMessage,
					"error":         err,
				}).Error("Problem to send TestInstruction to child-Gateway or Plugin")

			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "466ff0c7-eb34-4188-a1bc-d51bd3cda850",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send TestInstruction to child-Gateway or Plugin")

			}
		}
	}

	return err
}
