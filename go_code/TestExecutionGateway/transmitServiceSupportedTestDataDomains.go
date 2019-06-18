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
// Initiate Dispatch Engine for SendTestExecutionLogTowardsFenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateSupportedTestDataDomainsTowardsFenix() {

	// Start Transmit Engine, for SupportedTestDataDomains as a go-routine
	go gatewayObject.transmitEngineForSupportedTestDataDomainsTowardsFenix()
}

// ********************************************************************************************
// Forward SupportedTestDataDomainsWithHeadersMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForSupportedTestDataDomainsTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// ***** Wait for data comes from channel to transmit engine ****
			supportedTestDataDomainsWithHeadersMessageToBeForwarded := <-gatewayObject.supportedTestDataDomainsWithHeadersMessageTowardsFenixChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsWithHeadersMessageTowardsFenixChannel),
				"supportedTestDataDomainsWithHeadersMessageTowardsFenixChannel",
				"415c0a6b-c742-4a72-8ecd-47d16031f09a")

			logger.WithFields(logrus.Fields{
				"ID": "040072fd-2064-4955-bf31-8cf322e06eb0",
				"supportedTestDataDomainsWithHeadersMessageToBeForwarded": supportedTestDataDomainsWithHeadersMessageToBeForwarded,
			}).Debug("Received a new 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' from channel that shoud be forwarded")

			// ***** Send ExecutionLog to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"64820363-c600-467f-af32-507bba2b608e",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Parent (Gateway or Plugin) Server!",
				)

				// Convert supportedTestDataDomainsWithHeadersMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsWithHeadersMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestDataDomainsWithHeadersMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"ca4526f3-caea-4eb7-88d7-a56f5f85e272",
						gRPC.InformationMessage_FATAL,
						"supportedTestDataDomainsWithHeadersMessageToBeForwarded",
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.MessageId,
						supportedTestDataDomainsWithHeadersMessageToBeForwardedByteArray,
						BucketForResendOfSupportedTestDataDomains,
						"a3ac2811-ce67-4bc5-92d6-f4111034a78e",
					)
				}

			} else {
				//Connection OK

				// Convert supportedTestDataDomainsWithHeadersMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsWithHeadersMessageToBeForwardedByteArray, err := json.Marshal(*supportedTestDataDomainsWithHeadersMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"fe099b1d-1a78-40bb-9796-b43a57b796ea",
						gRPC.InformationMessage_FATAL,
						"supportedTestDataDomainsWithHeadersMessageToBeForwarded",
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' into a byte array, stopping futher processing of this message",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.MessageId,
						supportedTestDataDomainsWithHeadersMessageToBeForwardedByteArray,
						BucketForResendOfSupportedTestDataDomains,
						"f62bdc28-4bb1-4fcc-bc13-377b7f40de12",
					)

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					supportedTestDataDomainsWithHeadersMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					supportedTestDataDomainsWithHeadersMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.SupportedTestDataDomains(ctx, supportedTestDataDomainsWithHeadersMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"ee882c63-abbe-4d7e-9029-948538d14fc9",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							supportedTestDataDomainsWithHeadersMessageToBeForwarded.MessageId,
							supportedTestDataDomainsWithHeadersMessageToBeForwardedByteArray,
							BucketForResendOfSupportedTestDataDomains,
							"1ea4aabb-5277-4910-a0e6-de050b147528",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "d4f9af1d-0b09-4d0a-bc36-0120471b3bad",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
