package TestExecutionGateway

import (
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
		// Wait for data comes from channel to transmit engine
		testExecutionLogMessageToBeForwarded := <-gatewayObject.testExecutionLogMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":                                   "5c9fe63a-fd82-4ccd-8386-2b9c049e51a1",
			"testExecutionLogMessageToBeForwarded": testExecutionLogMessageToBeForwarded,
		}).Debug("Received a new 'testExecutionLogMessageToBeForwarded' from channel that shoud be forwarded")

		// Send TestInstruction to client using gRPC-call
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "863612af-86b3-431c-8270-9d335f3f0001",
				"addressToDial": addressToDial,
				"error message": err,
			}).Warning("Did not connect to Child (Gateway or Plugin) Server!")

			// Send Warning information to Fenix
			localInformationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
				OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
				SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
				SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
				MessageId:             generateUUID(),
				MessageType:           gRPC.InformationMessage_WARNING,
				Message:               "Did not connect to Child (Gateway or Plugin) Server!",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

			// TODO Add message to memmory cash for later resend
			// TODO Save message in localDB for later resend
		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "77167a87-99e8-49cc-ac3f-8966ac109658",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Parent-gateway- or Fenix-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			testExecutionLogMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
			testExecutionLogMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

			// Do gRPC-call to client gateway or Fenix
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendTestExecutionLogTowardsFenix(ctx, testExecutionLogMessageToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "840312be-e773-4e9d-ae05-0ec8e9dcca5c",
					"returnMessage": returnMessage,
					"error":         err,
				}).Warning("Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix")

				// Send Warning information to Fenix
				localInformationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
					OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
					SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
					SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
					MessageId:             generateUUID(),
					MessageType:           gRPC.InformationMessage_WARNING,
					Message:               "Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix",
					OrginalCreateDateTime: generaTimeStampUTC(),
				}

				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "fed7b18e-5cd5-485c-b14d-0b90ff720feb",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send OK of 'testExecutionLogMessageToBeForwarded' to Parent-Gateway or Fenix")

				// TODO Check for messages to Resend (If so then put them on channel)

			}
		}
	}
}
