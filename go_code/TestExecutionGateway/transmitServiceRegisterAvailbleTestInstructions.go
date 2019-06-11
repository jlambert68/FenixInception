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
			availbleTestInstructionAtPluginMessageTowardsFenix:= <-gatewayObject.availbleTestInstructionAtPluginMessageTowardsFenixChannel

			logger.WithFields(logrus.Fields{
				"ID": "bf74025f-3872-46d8-8abd-5ed58ebeec1b",
				"supportedTestDataDomainsMessageToBeForwarded": availbleTestInstructionAtPluginMessageTowardsFenix,
			}).Debug("Received a new 'supportedTestDataDomainsMessageToBeForwarded' from channel that shoud be forwarded")

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel),
				"supportedTestDataDomainsMessageTowardsFenixChannel",
				"ddf637ab-a7ca-4676-b347-aa5a88cf62b57")

			// ***** Send AvailableTestDataDomains to parent gateway Fenix using gRPC-call ****
			addressToDial := getParentAddressAndPort()

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"7ad2be08-2846-4704-ac69-af1a53101f22",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to Child (Gateway or Plugin) Server!",
				)

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*availbleTestInstructionAtPluginMessageTowardsFenix)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"3a90e5d5-1d03-49a5-ad47-3554064206f8",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						availbleTestInstructionAtPluginMessageTowardsFenix.String(),
						err.Error(),
						"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						availbleTestInstructionAtPluginMessageTowardsFenix. MessageId,
						supportedTestDataDomainsMessageToBeForwardedByteArray,
						BucketForResendOfLogMesagesToFenix,
						"2f548354-bb20-454e-904a-36932f341826",
					)
				}

			} else {
				//Connection OK

				// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
				supportedTestDataDomainsMessageToBeForwardedByteArray, err := json.Marshal(*availbleTestInstructionAtPluginMessageTowardsFenix)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"c5e1e433-5c66-4dea-991d-5c5faaade743",
						gRPC.InformationMessage_FATAL,
						"testExecutionLogMessageToBeForwarded",
						availbleTestInstructionAtPluginMessageTowardsFenix.String(),
						err.Error(),
						"Error when converting testExecutionLogMessageToBeForwarded into a byte array, stopping futher processing of this TestInstruction",
					)
				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					availbleTestInstructionAtPluginMessageTowardsFenix.SenderId = gatewayConfig.gatewayIdentification.gatewayId
					availbleTestInstructionAtPluginMessageTowardsFenix.SenderName = gatewayConfig.gatewayIdentification.gatewayName

					// Do gRPC-call to client gateway or Fenix
					ctx := context.Background()
					returnMessage, err := gatewayClient.RegistrateAvailableTestDataDomains(ctx, availbleTestInstructionAtPluginMessageTowardsFenix supportedTestDataDomainsMessageToBeForwarded)
					if err != nil {
						// Error when rending gRPC to parent
						LogErrorAndSendInfoToFenix(
							"9a9698f7-b09f-4ea6-84aa-76ba744916a3",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							availbleTestInstructionAtPluginMessageTowardsFenix.MessageId,
							supportedTestDataDomainsMessageToBeForwardedByteArray,
							BucketForResendOfLogMesagesToFenix,
							"cf0d3d61-3dec-4d56-b0eb-c399ca5d28e6",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "71bb9c49-4a32-436f-8645-6d519e4cad5f",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'testExecutionLogMessageToBeForwarded' to Parent-Gateway or Fenix")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
