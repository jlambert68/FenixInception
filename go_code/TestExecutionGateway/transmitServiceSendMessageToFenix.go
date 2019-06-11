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
// Initiate Transmit Engine for SendMessageToFenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateSendMessageToFenix() {

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.transmitEngineForSendMessageToFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForSendMessageToFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// Wait for data comes from channel to transmit engine
			informationMessageToBeForwarded := <-gatewayObject.informationMessageChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel),
				"informationMessageChannel",
				"01411586-7dab-4f34-809c-cdc8af8742c5")

			logger.WithFields(logrus.Fields{
				"ID":                              "3103a0a5-ff25-4e85-9939-5e8d72e26ba3",
				"informationMessageToBeForwarded": informationMessageToBeForwarded,
			}).Debug("Received a new informationMessage from channel that shoud be forwarded")

			// ***** Send ExecutionLog to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"6479b122-741d-4ee5-8c32-f9c8f64cd308",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Child (Gateway or Plugin) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*informationMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"f0a53fd1-b7ee-4c27-b018-1e39a4cf9fa0",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						informationMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'informationMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						informationMessageToBeForwarded.MessageId,
						testExecutionLogMessageToBeForwardedByteArray,
						BucketForResendOfInfoMessagesToFenix,
						"42da4af3-a4eb-4498-95e8-e78b13d2365a",
					)
				}
			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				informationMessageToBeForwardedByteArray, err := json.Marshal(*informationMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"e9366c2a-7289-4c13-a765-c9a91469eea2",
						gRPC.InformationMessage_FATAL,
						"informationMessageToBeForwarded",
						informationMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'informationMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						informationMessageToBeForwarded.MessageId,
						informationMessageToBeForwardedByteArray,
						BucketForResendOfInfoMessagesToFenix,
						"8672d03e-4d63-4126-b5ea-8f2c80c44a98",
					)

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					informationMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.gatewayId
					informationMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.gatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.SendMessageToFenix(ctx, informationMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"5a3b02fa-a8b7-4d2c-8eb7-55430c08c1b9",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							informationMessageToBeForwarded.MessageId,
							informationMessageToBeForwardedByteArray,
							BucketForResendOfInfoMessagesToFenix,
							"9782b4bc-afb9-424b-9dc4-43e6f47011b4",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "b18d86ff-5e6d-4f01-b40a-bd67ace547b2",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'informationMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
