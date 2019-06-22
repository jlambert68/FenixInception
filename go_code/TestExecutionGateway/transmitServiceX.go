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
// Initiate Transmit Engine for XTowardsFenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateTransmitEngineForXTowardsFenix() {

	// Start Transmit Engine, for X as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix()
}

// ********************************************************************************************
// Forward X-messages-messages from channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForXTowardsFenix() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// ***** Wait for data comes from channel to transmit engine ****

			testInstructionExecutionResultMessageToBeForwarded := <-gatewayObject.testInstructionExecutionResultMessageTowardsFenixChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.testInstructionExecutionResultMessageTowardsFenixChannel),
				"supportedTestDataDomainsWithHeadersMessageTowardsFenixChannel",
				"084a768d-412f-43cb-bfd6-80fcb777d36c")

			logger.WithFields(logrus.Fields{
				"ID": "1eb79d7e-9a44-4315-bbd4-9e008b0c55fc",
				"testInstructionExecutionResultMessageToBeForwarded": testInstructionExecutionResultMessageToBeForwarded,
			}).Debug("Received a new 'testInstructionExecutionResultMessageToBeForwarded' from channel that shoud be forwarded")

			// ***** Send ExecutionResult to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"8125cf1d-3cb9-4127-9b3a-6f495825a7cb",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Parent (Gateway or Plugin) Server!",
				)

				// Convert testInstructionExecutionResultMessageToBeForwarded-struct into a byte array
				stestInstructionExecutionResultMessageToBeForwardedByteArray, err := json.Marshal(*testInstructionExecutionResultMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"fbf6036a-8b82-4c03-958c-f831bf50946d",
						gRPC.InformationMessage_FATAL,
						"testInstructionExecutionResultMessageToBeForwarded",
						testInstructionExecutionResultMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'testInstructionExecutionResultMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						testInstructionExecutionResultMessageToBeForwarded.MessageId,
						stestInstructionExecutionResultMessageToBeForwardedByteArray,
						BucketForResendOfTestInstructionExecutionResult,
						"77024181-dd7d-4bd2-9fd3-22dce112d065",
					)
				}

			} else {
				//Connection OK

				// Convert testInstructionExecutionResultMessageToBeForwarded-struct into a byte array
				testInstructionExecutionResultMessageToBeForwardedByteArray, err := json.Marshal(*testInstructionExecutionResultMessageToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"7272adc6-58c5-4b0f-86bf-9ee1e929d004",
						gRPC.InformationMessage_FATAL,
						"testInstructionExecutionResultMessageToBeForwarded",
						testInstructionExecutionResultMessageToBeForwarded.String(),
						err.Error(),
						"Error when converting 'testInstructionExecutionResultMessageToBeForwarded' into a byte array, stopping futher processing of this message",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						testInstructionExecutionResultMessageToBeForwarded.MessageId,
						testInstructionExecutionResultMessageToBeForwardedByteArray,
						BucketForResendOfTestInstructionExecutionResult,
						"d86d45ee-cd99-4e82-9f45-2cee0db98ae4",
					)

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					testInstructionExecutionResultMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					testInstructionExecutionResultMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.SendTestInstructionResultTowardsFenix(ctx, testInstructionExecutionResultMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"519cc9e7-5a24-4787-a23e-e58808c22ce5",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							testInstructionExecutionResultMessageToBeForwarded.MessageId,
							testInstructionExecutionResultMessageToBeForwardedByteArray,
							BucketForResendOfTestInstructionExecutionResult,
							"56bd3634-973e-48fc-b142-53a5d49347b2",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "4e9fc5b1-a8d4-4b79-b442-f72291007c70",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'testInstructionExecutionResultMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
