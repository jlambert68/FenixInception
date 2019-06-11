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

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateSendTestExecutionLogTowardsFenix() {

	// Start Transmit Engine, for TestExecutionLogs as a go-routine
	go gatewayObject.transmitEngineForSendTestExecutionLogTowardsFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForSendTestExecutionLogTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// ***** Wait for data comes from channel to transmit engine ****
			testExecutionLogMessageToBeForwarded := <-gatewayObject.testExecutionLogMessageChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel),
				"testExecutionLogMessageChannel",
				"975c747f-1d0f-4fe4-a851-8fe2e17d7561")

			logger.WithFields(logrus.Fields{
				"ID":                                   "5c9fe63a-fd82-4ccd-8386-2b9c049e51a1",
				"testExecutionLogMessageToBeForwarded": testExecutionLogMessageToBeForwarded,
			}).Debug("Received a new 'testExecutionLogMessageToBeForwarded' from channel that shoud be forwarded")

			// ***** Send ExecutionLog to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"512db5b3-9f6f-4981-9520-d53d91e24748",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Parent (Gateway or Plugin) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*testExecutionLogMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"e395a57d-0d5d-4388-9579-41315a267533",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						testExecutionLogMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'testExecutionLogMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						testExecutionLogMessageToBeForwarded.LogMessageId,
						testExecutionLogMessageToBeForwardedByteArray,
						BucketForResendOfLogMesagesToFenix,
						"c22d1b76-2f52-4d76-8013-fcc65b94da2f",
					)
				}

			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				informationMessageToBeForwardedByteArray, err := json.Marshal(*testExecutionLogMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"aa15d9e4-0dc5-4e04-9c21-7a982e1030bd",
						gRPC.InformationMessage_FATAL,
						"informationMessageToBeForwarded",
						testExecutionLogMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'informationMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						testExecutionLogMessageToBeForwarded.LogMessageId,
						informationMessageToBeForwardedByteArray,
						BucketForResendOfLogMesagesToFenix,
						"ec1b7180-3501-42bb-a138-9653c53e1673",
					)

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					testExecutionLogMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.gatewayId
					testExecutionLogMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.gatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.SendTestExecutionLogTowardsFenix(ctx, testExecutionLogMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"5c967203-e5ef-45b1-bab5-521c34555d07",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							testExecutionLogMessageToBeForwarded.LogMessageId,
							informationMessageToBeForwardedByteArray,
							BucketForResendOfLogMesagesToFenix,
							"76c35bba-6655-4381-81ff-fd55dbab57ab",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "9f1b1a43-98f0-4d65-b91a-928cb0d83a4c",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'informationMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
