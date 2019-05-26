package TestExecutionGateway

import (
	"encoding/json"
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
	go gatewayObject.transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix()
}

// ********************************************************************************************
// Forward SupportedTestDataDomainsMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObject_struct) transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix() {

	var parentAddress gRPCClientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to transmit engine
		supportedTestDataDomainsMessageToBeForwarded := <-gatewayObject.supportedTestDataDomainsMessageTowardsFenixChannel

		gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
			"ID": "0259aa0d-a161-45de-ae77-17d317605a0b",
			"supportedTestDataDomainsMessageToBeForwarded": supportedTestDataDomainsMessageToBeForwarded,
		}).Debug("Received a new 'supportedTestDataDomainsMessageToBeForwarded' from channel that shoud be forwarded")

		// Create the channel that the client address should be sent back on
		returnParentAddressChannel := make(chan dbResultMessage_struct)
		//TODO Change Bucket-name and this key into variables in common_object.go
		// Get Clients address
		dbMessage := dbMessage_struct{
			DB_READ,
			"Parent",
			"ParenId",
			nil,
			returnParentAddressChannel}

		// Send Read message to database to receive address
		gatewayObject.gatewayCommonObjects.dbMessageQueue <- dbMessage

		// Wait for address from channel, then close the channel
		parentAddressByteArray := <-returnParentAddressChannel
		close(returnParentAddressChannel)

		// Convert saved json object into Go-struct
		err = json.Unmarshal(parentAddressByteArray.value, &parentAddress)
		if err != nil {
			// Problem with unmarshal the json object
			gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
				"ID":                      "a55a68d0-e6e6-4b0c-b02f-319600e39e96",
				"parentAddressByteArray,": parentAddressByteArray,
			}).Error("Can't unmarshal gRPCClients address object from database")
			//TODO Send Error information to Fenix
		} else {
			// Send SupportedTestDataDomainsMessage to client(parent gateway/Fenix) using gRPC-call
			addressToDial := parentAddress.clientIp + parentAddress.clientPort

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
					"ID":            "9abef370-58b6-44d6-bcf3-a91c9bac1911",
					"addressToDial": addressToDial,
					"error message": err,
				}).Warning("Did not connect to Parent (Gateway or Fenix) Server!")
				// TODO Send Error information to Fenix
				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
					"ID":            "05ff7bdd-435a-42a4-9921-aa5da8bc9d65",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to Parent-gateway- or Fenix-Server!")

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

				// ChangeSenderId to this gatway's SenderId before sending the data forward
				supportedTestDataDomainsMessageToBeForwarded.SenderId = CallingSystemId
				supportedTestDataDomainsMessageToBeForwarded.SenderName = CallingSystemName

				// Do gRPC-call to client gateway or Fenix
				ctx := context.Background()
				returnMessage, err := gatewayClient.RegistrateAvailableTestDataDomains(ctx, &supportedTestDataDomainsMessageToBeForwarded)
				if err != nil {
					gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
						"ID":            "08ad8a9f-ee08-4915-a458-eeb38dc14b40",
						"returnMessage": returnMessage,
						"error":         err,
					}).Warning("Problem to send 'supportedTestDataDomainsMessageToBeForwarded' to parent-Gateway or Fenix")
					// TODO Send Error information to Fenix
					// TODO Add message to memmory cash for later resend
					// TODO Save message in localDB for later resend
				} else {
					gatewayObject.gatewayCommonObjects.logger.WithFields(logrus.Fields{
						"ID":            "266c5829-9635-4cf1-bfe9-9f70470c96f9",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send OK of 'supportedTestDataDomainsMessageToBeForwarded' to Parent-Gateway or Fenix")

					// TODO Check for messages to Resend (If so then put them on channel)

				}
			}
		}
	}
}
