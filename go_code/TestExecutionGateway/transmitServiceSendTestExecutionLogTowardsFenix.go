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

		// ***** Send TestInstruction to client using gRPC-call ****
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			// Connection Not OK
			LogErrorAndSendInfoToFenixFor_SendTestExecutionLogTowardsFenix(
				"512db5b3-9f6f-4981-9520-d53d91e24748",
				gRPC.InformationMessage_WARNING,
				"addressToDial",
				addressToDial,
				err.Error(),
				"Did not connect to Child (Gateway or Plugin) Server!",
			)
		} else {
			//Connection OK

			// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
			testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*testExecutionLogMessageToBeForwarded)

			if err != nil {
				// Error when Unmarshaling to []byte
				LogErrorAndSendInfoToFenixFor_SendTestExecutionLogTowardsFenix(
					"20ba7b7d-a27e-4bb1-a846-3f7cde8ecf04",
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
					"b523d0a6-e3d6-42ea-8746-8c2864a943a5",
				)

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
					LogErrorAndSendInfoToFenixFor_SendTestExecutionLogTowardsFenix(
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

// *********************************************************************************
// Log message to local log and then Send message to Fenix
func LogErrorAndSendInfoToFenixFor_SendTestExecutionLogTowardsFenix(
	id string,
	messageType gRPC.InformationMessage_InformationType,
	infoHeader string,
	info string,
	errorMessage string,
	message string,
) {

	var suffix string = "_A"

	switch messageType {
	case gRPC.InformationMessage_DEBUG:
		// Only logg information and do not send to Fenix
		logger.WithFields(logrus.Fields{
			"ID":       id + suffix,
			infoHeader: info,
			"error":    errorMessage,
		}).Debug(message)

	case gRPC.InformationMessage_INFO:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id + suffix,
			infoHeader: info,
			"error":    errorMessage,
		}).Info(message)

		// Send information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_INFO,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_WARNING:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id + suffix,
			infoHeader: info,
			"error":    errorMessage,
		}).Warning(message)

		// Send Warning information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_WARNING,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_ERROR:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id + suffix,
			infoHeader: info,
			"error":    errorMessage,
		}).Error(message)

		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_ERROR,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_FATAL:
		// Only log and then Terminate Gateway due to problems
		logger.WithFields(logrus.Fields{
			"ID":       id + suffix,
			infoHeader: info,
			"error":    errorMessage,
		}).Fatal(message)

	default:
		logger.WithFields(logrus.Fields{
			"ID":          "d83326f4-1b06-4d13-8010-70b6c829cc88",
			"messageType": messageType,
		}).Fatal("Unknown messageType, stopping gateway")

	}
}
