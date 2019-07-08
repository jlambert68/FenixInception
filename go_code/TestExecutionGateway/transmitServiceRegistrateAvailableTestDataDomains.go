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

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateRegistrateAvailableTestDataDomainsTowardsFenix() {

	// Start Transmit Engine, for SupportedTestDataDomainsMessage-messages as a go-routine
	logger.WithFields(logrus.Fields{
		"ID": "c9b7074f-5529-4ad7-a2b9-19c27b59ddd3",
	}).Info("Initiate: 'transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix'")

	go gatewayObject.transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix()
}

// ********************************************************************************************
// Forward supportedTestDataDomainsMessageToBeForwarded-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// Wait for data comes from channel to transmit engine
			supportedTestDataDomainsMessageToBeForwarded := <-gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix

			logger.WithFields(logrus.Fields{
				"ID": "0259aa0d-a161-45de-ae77-17d317605a0b",
				"supportedTestDataDomainsMessageToBeForwarded": supportedTestDataDomainsMessageToBeForwarded,
			}).Debug("Received a new 'supportedTestDataDomainsMessageToBeForwarded' from channel that shoud be forwarded")

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix),
				"supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix",
				"7f51813f-07f2-4e05-9143-ecc40e59bb17")

			// ***** Send AvailableTestDataDomains to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"b479a416-c783-4b3a-9f8d-5c008d85fa9f",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Child (Gateway or Plugin) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestDataDomainsMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"661b6fb7-d125-4e8f-91b4-3d655d39963a",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						supportedTestDataDomainsMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						supportedTestDataDomainsMessageToBeForwarded.MessageId,
						supportedTestDataDomainsMessageToBeForwardedByteArray,
						BucketForResendOfLogMesagesTowardsFenix,
						"0a71308b-a172-4342-b3da-b45cf923860b",
					)
				}

			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestDataDomainsMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"891a385e-03eb-4cbd-9f2b-2ba1f70655fd",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						supportedTestDataDomainsMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					supportedTestDataDomainsMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					supportedTestDataDomainsMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.RegistrateAvailableTestDataDomains(ctx, supportedTestDataDomainsMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"256bb8ba-3348-436b-8817-aef590d312f1",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							supportedTestDataDomainsMessageToBeForwarded.MessageId,
							supportedTestDataDomainsMessageToBeForwardedByteArray,
							BucketForResendOfLogMesagesTowardsFenix,
							"338075e7-b43f-4cbf-be84-b5096c4658c1",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "8a0daa7f-c9ee-4c59-a982-75c741daeebd",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'testExecutionLogMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
