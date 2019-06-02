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
// Initiate Dispatch Engine for SupportedTestDomains
//

func (gatewayObject *gatewayTowardsPluginObject_struct) initiateDispatchEngineForSupportedTestDomains() {

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.dispatchEngineForSupportedTestDomains()
}

// ********************************************************************************************
// Forward getSupportedTestDomains-messages from incoming channel to child gateway or plugin
//

func (gatewayObject *gatewayTowardsPluginObject_struct) dispatchEngineForSupportedTestDomains() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(SERVICE_SLEEP_TIME * time.Second)

		} else {
			// Run service and process messages

			// ***** Wait for data comes from channel to dispatch engine *****
			getSupportedTestDomainsToBeForwarded := <-gatewayObject.supportedTestDataDomainsRequestChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.testInstructionMessageChannel),
				"supportedTestDataDomainsRequestChannel",
				"3c85d4c4-260d-4612-a153-2eeeff864621")

			logger.WithFields(logrus.Fields{
				"ID":                                   "05f9ae3c-dc0f-4107-ae85-ae57d6dde6c1",
				"testExecutionLogMessageToBeForwarded": getSupportedTestDomainsToBeForwarded,
			}).Debug("Received a new 'getSupportedTestDomainsToBeForwarded' from channel that shoud be forwarded")

			// ***** Send getSupportedTestDomains-message to client gateway/plugin using gRPC-call *****
			addressToDial := getClientAddressAndPort(getSupportedTestDomainsToBeForwarded.PluginId)

			// Set up connection to Parent Gateway or Fenix
			remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"607ff7f6-4452-49f6-a9ec-608118073a3b",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to child (Gateway or Plugin) Server!",
				)

				// Convert getSupportedTestDomainsToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*getSupportedTestDomainsToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"f1c5f3f6-341e-4fb3-96aa-9886da021c8d",
						gRPC.InformationMessage_FATAL,
						"getSupportedTestDomainsToBeForwarded",
						getSupportedTestDomainsToBeForwarded.String(),
						err.Error(),
						"Error when converting 'getSupportedTestDomainsToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						getSupportedTestDomainsToBeForwarded.MessageId,
						testExecutionLogMessageToBeForwardedByteArray,
						BUCKET_RESEND_GET_TESTDATA_DOMAINS_TO_PLUGIN,
						"1abe7614-c951-4186-821f-deb2b4203f64",
					)
				}

			} else {
				//Connection OK
				logger.WithFields(logrus.Fields{
					"ID":            "954306f9-d890-48ad-ac5a-68143166b48d",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to child-gateway/Plugin!")

				// Convert getSupportedTestDomainsToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*getSupportedTestDomainsToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"12c16a5f-ba94-4e16-9607-214b9281da9a",
						gRPC.InformationMessage_FATAL,
						"getSupportedTestDomainsToBeForwarded",
						getSupportedTestDomainsToBeForwarded.String(),
						err.Error(),
						"Error when converting 'getSupportedTestDomainsToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					getSupportedTestDomainsToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
					getSupportedTestDomainsToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

					// Do gRPC-call to child gateway or Plugin
					ctx := context.Background()
					returnMessage, err := gatewayClient.GetSupportedTestDataDomains(ctx, getSupportedTestDomainsToBeForwarded)
					if err != nil {
						// Error when rending gRPC to child
						LogErrorAndSendInfoToFenix(
							"6b77117b-e98c-4451-a91a-d3889438b3d7",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'getSupportedTestDomainsToBeForwarded' to child-Gateway or Plugin",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							getSupportedTestDomainsToBeForwarded.MessageId,
							testExecutionLogMessageToBeForwardedByteArray,
							BUCKET_RESEND_GET_TESTDATA_DOMAINS_TO_PLUGIN,
							"9446b644-66f3-452a-9436-add798b0dad9",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "77eef2f1-116d-4342-a02e-dce88d1e1b52",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'getSupportedTestDomainsToBeForwarded' to child-Gateway or Plugin")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
