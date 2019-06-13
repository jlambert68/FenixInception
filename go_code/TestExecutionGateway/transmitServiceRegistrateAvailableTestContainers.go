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

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateRegistrateAvailableTestContainersTowardsFenix() {

	// Start Transmit Engine, for SupportedTestDataDomainsMessage-messages as a go-routine
	logger.WithFields(logrus.Fields{
		"ID": "75042917-b443-402d-ab42-7de628e6939a",
	}).Info("Initiate: 'transmitEngineForRegistrateAvailableTestContainersTowardsFenix'")

	go gatewayObject.transmitEngineForRegistrateAvailableTestContainersTowardsFenix()
}

// ********************************************************************************************
// Forward supportedTestDataDomainsMessageToBeForwarded-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForRegistrateAvailableTestContainersTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// Wait for data comes from channel to transmit engine
			supportedTestContainerMessageToBeForwarded := <-gatewayObject.availbleTestContainersAtPluginMessageTowardsFenixChannel

			logger.WithFields(logrus.Fields{
				"ID": "39908e44-789e-4a8c-ba70-66d6c350b090",
				"supportedTestContainerMessageToBeForwarded": supportedTestContainerMessageToBeForwarded,
			}).Debug("Received a new 'supportedTestContainerMessageToBeForwarded' from channel that shoud be forwarded")

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.availbleTestContainersAtPluginMessageTowardsFenixChannel),
				"availbleTestContainersAtPluginMessageTowardsFenixChannel",
				"ce1103cf-a98f-4228-97a2-ac9a76353464")

			// ***** Send AvailableTestContainera to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"9a026040-3d0f-422d-a52b-5fe665c43b6b",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Child (Gateway or Plugin) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestContainerMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestContainerMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"b555ea59-685e-45fc-ba43-a04e9da6305a",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						supportedTestContainerMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting supportedTestContainerMessageToBeForwarded into a byte array, stopping futher processing of this TestContainersObject",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						supportedTestContainerMessageToBeForwarded.MessageId,
						supportedTestContainerMessageToBeForwardedByteArray,
						BucketForResendOfAvailableTestContainers,
						"922b5520-a374-4762-824e-f2b831b83bc9",
					)
				}

			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestContainerMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestContainerMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"79dede01-1b2f-449e-884d-6c97ccb3bc00",
						gRPC.InformationMessage_FATAL,
						"supportedTestContainerMessageToBeForwarded",
						supportedTestContainerMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting supportedTestContainerMessageToBeForwarded into a byte array, stopping futher processing of this TestContersObject",
					)
				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					supportedTestContainerMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.gatewayId
					supportedTestContainerMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.gatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.RegistrateAailableTestContainers(ctx, supportedTestContainerMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"eab4885f-e018-4416-963e-7f303faf57f8",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'supportedTestContainerMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							supportedTestContainerMessageToBeForwarded.MessageId,
							supportedTestContainerMessageToBeForwardedByteArray,
							BucketForResendOfAvailableTestContainers,
							"85d70e29-a698-4b5c-a9ef-ada7e5ac4695",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "d0476286-0b30-4b65-ba21-1e1dc6a3ee9e",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'supportedTestContainerMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
