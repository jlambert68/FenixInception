package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Dispatch Engine for SupportedTestDomains
//

func (gatewayObject *GatewayTowardsPluginObject_struct) initiateDispatchEngineForSupportedTestDomains() {

	// Create return Channel
	returnChannel = make(chan dbResultMessage_struct)

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.dispatchEngineForSupportedTestDomains()
}

// ********************************************************************************************
// Forward getSupportedTestDomains-messages from incoming channel to child gateway or plugin
//

func (gatewayObject *GatewayTowardsPluginObject_struct) dispatchEngineForSupportedTestDomains() {

	var clientAddress gRPCClientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to dispatch engine
		getSupportedTestDomainsToBeForwarded := <-gatewayObject.supportedTestDataDomainsRequestChannel

		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":                                   "6478c599-ac46-45f2-9c06-6575c01854e4",
			"getSupportedTestDomainsToBeForwarded": getSupportedTestDomainsToBeForwarded,
		}).Debug("Received a new getSupportedTestDomains-message from channel that shoud be forwarded")

		// Create the channel that the client address should be sent back on
		returnClientAddressChannel := make(chan dbResultMessage_struct)

		// Get Clients address
		dbMessage := dbMessage_struct{
			DB_READ,
			"Clients",
			getSupportedTestDomainsToBeForwarded.PluginId,
			nil,
			returnClientAddressChannel}

		// Send Read message to database to receive address
		gatewayObject.dbMessageQueue <- dbMessage

		// Wait for address from channel, then close the channel
		clientAddressByteArray := <-returnClientAddressChannel
		close(returnClientAddressChannel)

		// Convert saved json object into Go-struct
		err = json.Unmarshal(clientAddressByteArray.value, &clientAddress)
		if err != nil {
			// Problem with unmarshal the json object
			gatewayObject.logger.WithFields(logrus.Fields{
				"ID": "74bca1b9-2fe2-47cf-a776-6968d16921eb",
				"getSupportedTestDomainsToBeForwarded.PluginId,": getSupportedTestDomainsToBeForwarded.PluginId,
			}).Error("Can't unmarshal Client address object from Database")
			//TODO Send Error information to Fenix
		} else {
			// Send TestInstruction to client using gRPC-call
			addressToDial := clientAddress.clientIp + clientAddress.clientPort

			// Set up connection to Client Gateway or Plugin
			remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "83830dbc-edbe-4087-a5a9-d164fb54c395",
					"addressToDial": addressToDial,
					"error message": err,
				}).Error("Did not connect to Child (Gateway or Plugin) Server!")
				//TODO Send Error information to Fenix
			} else {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "7f8910b0-752a-4e80-8210-8f9f2f19dbbc",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to child-gateway- or Plugin-Server!")

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

				// Do gRPC-call to client gateway or Plugin
				ctx := context.Background()
				returnMessage, err := gatewayClient.GetSupportedTestDataDomains(ctx, &getSupportedTestDomainsToBeForwarded)
				if err != nil {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "f549c867-1250-4276-af6e-901908cd6221",
						"returnMessage": returnMessage,
						"error":         err,
					}).Error("Problem to send getSupportedTestDomains to child-Gateway or Plugin")
					//TODO Send Error information to Fenix
				} else {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "85ee050c-6880-4c33-8f27-4f09380a5a67",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send getSupportedTestDomains to child-Gateway or Plugin")

				}
			}
		}
	}
}
