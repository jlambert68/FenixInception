package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Transmit Engine for SupportedTestDataDomainsMessage-messages
//

func (gatewayObject *gatewayTowardsFenixObject_struct) initiateRegistrateAvailableTestDataDomainsTowardsFenix() {

	// Start Transmit Engine, for SupportedTestDataDomainsMessage-messages as a go-routine
	logger.WithFields(logrus.Fields{
		"ID": "c9b7074f-5529-4ad7-a2b9-19c27b59ddd3",
	}).Info("Initiate: 'transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix'")

	go gatewayObject.transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix()
}

// ********************************************************************************************
// Forward SupportedTestDataDomainsMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix() {

	for {
		// Wait for data comes from channel to transmit engine
		supportedTestDataDomainsMessageToBeForwarded := <-gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel

		logger.WithFields(logrus.Fields{
			"ID": "0259aa0d-a161-45de-ae77-17d317605a0b",
			"supportedTestDataDomainsMessageToBeForwarded": supportedTestDataDomainsMessageToBeForwarded,
		}).Debug("Received a new 'supportedTestDataDomainsMessageToBeForwarded' from channel that shoud be forwarded")

		// Send SupportedTestDataDomainsMessage to client(parent gateway/Fenix) using gRPC-call
		addressToDial := getParentAddressAndPort()

		// Set up connection to Parent Gateway or Fenix
		remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "9abef370-58b6-44d6-bcf3-a91c9bac1911",
				"addressToDial": addressToDial,
				"error message": err,
			}).Warning("Did not connect to Parent (Gateway or Fenix) Server!, will retry")

			// Send Warning information to Fenix
			localInformationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
				OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
				SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
				SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
				MessageId:             generateUUID(),
				MessageType:           gRPC.InformationMessage_WARNING,
				Message:               "Got an error when Saveing to local DB, will retry",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

			// TODO Add message to memmory cash for later resend
			// TODO Save message in localDB for later resend
		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "05ff7bdd-435a-42a4-9921-aa5da8bc9d65",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Parent-gateway- or Fenix-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			supportedTestDataDomainsMessageToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
			supportedTestDataDomainsMessageToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

			// Do gRPC-call to client gateway or Fenix
			ctx := context.Background()
			returnMessage, err := gatewayClient.RegistrateAvailableTestDataDomains(ctx, supportedTestDataDomainsMessageToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "08ad8a9f-ee08-4915-a458-eeb38dc14b40",
					"returnMessage": returnMessage,
					"error":         err,
				}).Warning("Problem to send 'supportedTestDataDomainsMessageToBeForwarded' to parent-Gateway or Fenix")

				// Send Warning information to Fenix
				localInformationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
					OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
					SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
					SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
					MessageId:             generateUUID(),
					MessageType:           gRPC.InformationMessage_WARNING,
					Message:               "Problem to send 'supportedTestDataDomainsMessageToBeForwarded' to parent-Gateway or Fenix",
					OrginalCreateDateTime: generaTimeStampUTC(),
				}

				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "266c5829-9635-4cf1-bfe9-9f70470c96f9",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send OK of 'supportedTestDataDomainsMessageToBeForwarded' to Parent-Gateway or Fenix")

				// TODO Check for messages to Resend (If so then put them on channel)

			}
		}
	}
}
