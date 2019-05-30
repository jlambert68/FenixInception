package TestExecutionGateway

import (
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

		// Send TestInstruction to client using gRPC-call
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "6118ae80-1d86-4530-83e4-549c00d01337",
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
				Message:               "Did not connect to parent (Gateway or Fenix) Server!",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

			// TODO Save message in localDB for later resend
		} else {
			// Connection to parent server OK
			logger.WithFields(logrus.Fields{
				"ID":            "e1c1903e-d72b-4ca2-973e-33b8525cb6ee",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Parent-gateway- or Fenix-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			testInstructionTimeOutMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
			testInstructionTimeOutMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

			// Do gRPC-call to client gateway or Fenix
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendTestInstructionTimeOutTowardsFenix(ctx, testInstructionTimeOutMessageToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "1ae9d406-b8fc-4622-b6a3-8b2c0ce3cdc9",
					"returnMessage": returnMessage,
					"error":         err,
				}).Warning("Problem to send 'testInstructionTimeOutMessageToBeForwarded' to parent-Gateway or Fenix")

				// Send Warning information to Fenix
				localInformationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
					OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
					SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
					SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
					MessageId:             generateUUID(),
					MessageType:           gRPC.InformationMessage_WARNING,
					Message:               "Problem to send 'testInstructionTimeOutMessageToBeForwarded' to parent-Gateway or Fenix",
					OrginalCreateDateTime: generaTimeStampUTC(),
				}

				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				// gRPC messagage send parent server
				logger.WithFields(logrus.Fields{
					"ID":            "985c3a8b-dd01-496a-b3b8-9e1c67b89dd6",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send OK of 'testInstructionTimeOutMessageToBeForwarded' to Parent-Gateway or Fenix")

				// TODO Check for messages to Resend (If so then put them on channel) and remove from DB

			}
		}
	}
}
