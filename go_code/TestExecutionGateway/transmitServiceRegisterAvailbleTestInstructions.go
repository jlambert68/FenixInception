package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"time"
)

// ********************************************************************************************
// Initiate Transmit Engine for SupportedTestDataDomainsMessage-messages
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateRegistrateAvailableTestInstructionsTowardsFenix() {

	// Start Transmit Engine, for SupportedTestDataDomainsMessage-messages as a go-routine
	logger.WithFields(logrus.Fields{
		"ID": "9c1a8e90-fd0a-469d-8c74-ecb8a16b9b61",
	}).Info("Initiate: 'transmitEngineForRegistrateAvailableTestInstructionsTowardsFenix'")

	go gatewayObject.transmitEngineForRegistrateAvailableTestInstructionsTowardsFenix()
}

// ********************************************************************************************
// Forward supportedTestDataDomainsMessageToBeForwarded-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForRegistrateAvailableTestInstructionsTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// Wait for data comes from channel to transmit engine
			availbleTestInstructionAtPluginMessageTowardsFenix := <-gatewayObject.availbleTestInstructionAtPluginMessageTowardsFenixChannel

			logger.WithFields(logrus.Fields{
				"ID": "0c50a6da-f82b-4fde-b1ca-5219f5c3d2ab",
				"supportedTestDataDomainsMessageToBeForwarded": availbleTestInstructionAtPluginMessageTowardsFenix,
			}).Debug("Received a new 'availbleTestInstructionMessageToBeForwarded' from channel that shoud be forwarded")

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel),
				"availbleTestInstructionAtPluginMessageTowardsFenixChannel",
				"3b482f8f-4809-4784-b821-07b384d80a9c")

			// ***** Send AvailableTestDataDomains to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"1a3b70a1-a7a6-46b1-b2f8-30c4c21d5e4d",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to parent (Gateway or Fenix) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*availbleTestInstructionAtPluginMessageTowardsFenix)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"af0ccd4a-95df-47db-90cd-ce3ea5e45cee",
						gRPC.InformationMessage_FATAL,
						"availbleTestInstructionMessageToBeForwarded",
						availbleTestInstructionAtPluginMessageTowardsFenix.String(),
						err.Error(),
						"Error when converting availbleTestInstructionMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						availbleTestInstructionAtPluginMessageTowardsFenix.MessageId,
						supportedTestDataDomainsMessageToBeForwardedByteArray,
						BucketForResendOfAvailableTestInstructionsToFenix,
						"1ecdf91a-b22e-4773-b58a-db5391c1a2e0",
					)
				}

			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*availbleTestInstructionAtPluginMessageTowardsFenix)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"ac32dcd3-68bb-47f3-b832-622acc970aaf",
						gRPC.InformationMessage_FATAL,
						"availbleTestInstructionMessageToBeForwarded",
						availbleTestInstructionAtPluginMessageTowardsFenix.String(),
						err.Error(),
						"Error when converting availbleTestInstructionMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					availbleTestInstructionAtPluginMessageTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					availbleTestInstructionAtPluginMessageTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.RegisterAvailbleTestInstructions(ctx, availbleTestInstructionAtPluginMessageTowardsFenix)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"303ed68b-80d1-4450-bb21-fd5ce407ac27",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'availbleTestInstructionMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							availbleTestInstructionAtPluginMessageTowardsFenix.MessageId,
							supportedTestDataDomainsMessageToBeForwardedByteArray,
							BucketForResendOfAvailableTestInstructionsToFenix,
							"75518180-5c81-47bd-8819-f2988e1ba00e",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "20be8dd7-df8b-46b9-affe-6ce238a7a13a",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'availbleTestInstructionMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
