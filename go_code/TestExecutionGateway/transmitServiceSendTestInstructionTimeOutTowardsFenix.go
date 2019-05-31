package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Transmit Engine for InformationMessage-messages
//

func (gatewayObject *gatewayTowardsFenixObject_struct) initiateSendTestInstructionTimeOutTowardsFenix() {

	// Start Transmit Engine, for InformationMessage-messages as a go-routine
	go gatewayObject.transmitEngineForSendTestInstructionTimeOutTowardsFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) transmitEngineForSendTestInstructionTimeOutTowardsFenix() {

	for {
		// Wait for data comes from channel to transmit engine
		testInstructionTimeOutMessageToBeForwarded := <-gatewayObject.testInstructionTimeOutMessageChannel

		logger.WithFields(logrus.Fields{
			"ID": "d1d4385b-b7c1-473d-8105-2e9b2341ef14",
			"testInstructionTimeOutMessageToBeForwarded": testInstructionTimeOutMessageToBeForwarded,
		}).Debug("Received a new 'testInstructionTimeOutMessageToBeForwarded' from channel that shoud be forwarded")

		// ***** Send Timeout to parent gateway Fenix using gRPC-call ****
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			// Connection Not OK
			LogErrorAndSendInfoToFenix(
				"f4ba70b6-06e0-46cc-9f82-fcd81992d8a3",
				gRPC.InformationMessage_WARNING,
				"addressToDial",
				addressToDial,
				err.Error(),
				"Did not connect to Child (Gateway or Plugin) Server!",
			)
		} else {
			//Connection OK

			// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
			testInstructionTimeOutMessageToBeForwardedByteArray, err := json.Marshal(*testInstructionTimeOutMessageToBeForwarded)

			if err != nil {
				// Error when Unmarshaling to []byte
				LogErrorAndSendInfoToFenix(
					"b0883ffe-33f3-4204-9cb2-e32ef879597b",
					gRPC.InformationMessage_FATAL,
					"testExecutionLogMessageToBeForwarded",
					testInstructionTimeOutMessageToBeForwarded.String(),
					err.Error(),
					"Error when converting 'testInstructionTimeOutMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
				)
			} else {
				// Marshaling to []byte OK

				// Save message to local DB for later processing
				SaveMessageToLocalDB(
					testInstructionTimeOutMessageToBeForwarded.MessageId,
					testInstructionTimeOutMessageToBeForwardedByteArray,
					BUCKET_RESEND_LOG_MESSAGES_TO_FENIX,
					"3b0223ad-7c4a-4fe3-9684-c28214d3f2b5",
				)

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

				// ChangeSenderId to this gatway's SenderId before sending the data forward
				testInstructionTimeOutMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
				testInstructionTimeOutMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

				// Do gRPC-call to client gateway or Fenix
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestInstructionTimeOutTowardsFenix(ctx, testInstructionTimeOutMessageToBeForwarded)
				if err != nil {
					// Error when rending gRPC to parent
					LogErrorAndSendInfoToFenix(
						"6c3ed5a3-9e16-4c19-bbc9-a3b84b921ea0",
						gRPC.InformationMessage_WARNING,
						"returnMessage",
						returnMessage.String(),
						err.Error(),
						"Problem to send 'testInstructionTimeOutMessageToBeForwarded' to parent-Gateway or Fenix",
					)

					// Save message to local DB for later processing
					SaveMessageToLocalDB(
						testInstructionTimeOutMessageToBeForwarded.MessageId,
						testInstructionTimeOutMessageToBeForwardedByteArray,
						BUCKET_RESEND_LOG_MESSAGES_TO_FENIX,
						"ae47bac5-fd22-41dd-80b0-16063d880988",
					)

				} else {
					// gRPC Send message OK
					logger.WithFields(logrus.Fields{
						"ID":            "",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send OK of 'testInstructionTimeOutMessageToBeForwarded' to Parent-Gateway or Fenix")

					// TODO Check for messages to Resend (If so then put them on channel)

				}
			}
		}
	}
}
