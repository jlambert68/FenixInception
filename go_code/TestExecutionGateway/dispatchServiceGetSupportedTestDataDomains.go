package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Dispatch Engine for SupportedTestDomains
//

func (gatewayObject *gatewayTowardsPluginObject_struct) initiateDispatchEngineForSupportedTestDomains() {

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.dispatchEngineForSupportedTestDomains()
}

// ********************************************************************************************
// Forward getSupportedTestDomains-messages from incoming channel to child gateway or plugin
//

func (gatewayObject *gatewayTowardsPluginObject_struct) dispatchEngineForSupportedTestDomains() {

	for {
		// Wait for data comes from channel to dispatch engine
		getSupportedTestDomainsToBeForwarded := <-gatewayObject.supportedTestDataDomainsRequestChannel

		logger.WithFields(logrus.Fields{
			"ID":                                   "6478c599-ac46-45f2-9c06-6575c01854e4",
			"getSupportedTestDomainsToBeForwarded": getSupportedTestDomainsToBeForwarded,
		}).Debug("Received a new getSupportedTestDomains-message from channel that shoud be forwarded")

		// Send TestInstruction to client using gRPC-call
		addressToDial := getClientAddressAndPort(getSupportedTestDomainsToBeForwarded.PluginId)

		// Set up connection to Client Gateway or Plugin
		remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "83830dbc-edbe-4087-a5a9-d164fb54c395",
				"addressToDial": addressToDial,
				"error message": err,
			}).Error("Did not connect to Child (Gateway or Plugin) Server!")

			// Send Error information to Fenix
			//Send Error information to Fenix
			localInformationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
				OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
				SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
				SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
				MessageId:             generateUUID(),
				MessageType:           gRPC.InformationMessage_ERROR,
				Message:               "Did not connect to Child (Gateway or Plugin) Server!",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "7f8910b0-752a-4e80-8210-8f9f2f19dbbc",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to child-gateway- or Plugin-Server!")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			getSupportedTestDomainsToBeForwarded.SenderId = gatewayConfig.gatewayIdentification.callingSystemId
			getSupportedTestDomainsToBeForwarded.SenderName = gatewayConfig.gatewayIdentification.callingSystemName

			// Do gRPC-call to client gateway or Plugin
			ctx := context.Background()
			returnMessage, err := gatewayClient.GetSupportedTestDataDomains(ctx, getSupportedTestDomainsToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "f549c867-1250-4276-af6e-901908cd6221",
					"returnMessage": returnMessage,
					"error":         err,
				}).Error("Problem to send getSupportedTestDomains to child-Gateway or Plugin")

				// Send Error information to Fenix
				localInformationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
					OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
					SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
					SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
					MessageId:             generateUUID(),
					MessageType:           gRPC.InformationMessage_ERROR,
					Message:               "Problem to send getSupportedTestDomains to child-Gateway or Plugin",
					OrginalCreateDateTime: generaTimeStampUTC(),
				}

			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "85ee050c-6880-4c33-8f27-4f09380a5a67",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send getSupportedTestDomains to child-Gateway or Plugin")

			}
		}
	}
}
