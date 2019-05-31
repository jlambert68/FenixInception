package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Dispatch Engine for SendTestExecutionLogTowardsFenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) initiateSendTestExecutionLogTowardsFenix() {

	// Start Transmit Engine, for TestExecutionLogs as a go-routine
	go gatewayObject.transmitEngineForSendTestExecutionLogTowardsFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) transmitEngineForSendTestExecutionLogTowardsFenix() {

	for {
		// TODO HALT forwarding of messages if Fenix says so

		// ***** Wait for data comes from channel to transmit engine ****
		testExecutionLogMessageToBeForwarded := <-gatewayObject.testExecutionLogMessageChannel

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
					"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
				)

			} else {
				// Marshaling to []byte OK

				// Save message to local DB for later processing
				SaveMessageToLocalDB(
					testExecutionLogMessageToBeForwarded.LogMessageId,
					testExecutionLogMessageToBeForwardedByteArray,
					BUCKET_RESEND_LOG_MESSAGES_TO_FENIX,
					"c22d1b76-2f52-4d76-8013-fcc65b94da2f",
				)
			}

		} else {
			//Connection OK
			logger.WithFields(logrus.Fields{
				"ID":            "9f72b39f-34c7-417c-bba0-0e0630a4f7a1",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to parent-gateway/Fenix!")

			// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
			testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*testExecutionLogMessageToBeForwarded)

			if err != nil {
				// Error when Unmarshaling to []byte
				LogErrorAndSendInfoToFenix(
					"20ba7b7d-a27e-4bb1-a846-3f7cde8ecf04",
					gRPC.InformationMessage_FATAL,
					"testExecutionLogMessageToBeForwarded",
					testExecutionLogMessageToBeForwarded.String(),
					err.Error(),
					"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
				)

			} else {
				// Marshaling to []byte OK

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

				// ChangeSenderId to this gatway's SenderId before sending the data forward
				testExecutionLogMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
				testExecutionLogMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

				// Do gRPC-call to client gateway or Fenix
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestExecutionLogTowardsFenix(ctx, testExecutionLogMessageToBeForwarded)
				if err != nil {
					// Error when rending gRPC to parent
					LogErrorAndSendInfoToFenix(
						"59939fa7-22b6-4595-a70b-816416cc228f",
						gRPC.InformationMessage_WARNING,
						"returnMessage",
						returnMessage.String(),
						err.Error(),
						"Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix",
					)

					// Save message to local DB for later processing
					SaveMessageToLocalDB(
						testExecutionLogMessageToBeForwarded.LogMessageId,
						testExecutionLogMessageToBeForwardedByteArray,
						BUCKET_RESEND_LOG_MESSAGES_TO_FENIX,
						"421438e6-a99f-407b-b338-51cb69d061e2",
					)

				} else {
					// gRPC Send message OK
					logger.WithFields(logrus.Fields{
						"ID":            "fed7b18e-5cd5-485c-b14d-0b90ff720feb",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send OK of 'testExecutionLogMessageToBeForwarded' to Parent-Gateway or Fenix")

					// TODO Check for messages to Resend (If so then put them on channel)

				}
			}
		}
	}
}
