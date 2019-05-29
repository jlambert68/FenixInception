package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Transmit Engine for SendMessageToFenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) initiateSendMessageToFenix() {

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.transmitEngineForSendMessageToFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) transmitEngineForSendMessageToFenix() {

	for {

		// TODO for all traffic back to Fenix processing must support HALTing of sending
		// Wait for data comes from channel to transmit engine
		informationMessageToBeForwarded := <-gatewayObject.informationMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":                              "3103a0a5-ff25-4e85-9939-5e8d72e26ba3",
			"informationMessageToBeForwarded": informationMessageToBeForwarded,
		}).Debug("Received a new informationMessage from channel that shoud be forwarded")

		// Send TestInstruction to client using gRPC-call
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "0c8ba004-1402-46d4-ac3b-c49863ff817e",
				"addressToDial": addressToDial,
				"error message": err,
			}).Warning("Did not connect to Child (Gateway or Plugin) Server!")
			// TODO Send Error information to Fenix
			// TODO Add message to memmory cash for later resend
			// TODO Save message in localDB for later resend
		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "5ae8053c-1aca-4fc5-9d01-f2c5577933cc",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Parent-gateway- or Plugin-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			informationMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
			informationMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

			// Do gRPC-call to client gateway or Fenix
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendMessageToFenix(ctx, informationMessageToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "3becc080-360b-42b4-8fe8-32a01bf820bc",
					"returnMessage": returnMessage,
					"error":         err,
				}).Warning("Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix")
				// TODO Send Error information to Fenix
				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "860b622b-c2fb-442c-b4e1-ac5fdf4f8d36",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send OK of 'informationMessageToBeForwarded' to parent-Gateway or Fenix")

				// TODO Check for messages to Resend (If so then put them on channel)

			}
		}
	}
}
