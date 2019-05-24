package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Initiate Dispatch Engine for SendMessageToFenix
//

func (gatewayObject *GatewayTowardsPluginObject_struct) initiateSendTestInstructionTimeOutTowardsFenix() {

	// Start Dispatch Engine, for SupportedTestDomains as a go-routine
	go gatewayObject.transmitEngineForSendTestInstructionTimeOutTowardsFenix()
}

// ********************************************************************************************
// Forward InformationMessage-messages from incoming channel towards Fenix
//

func (gatewayObject *GatewayTowardsPluginObject_struct) transmitEngineForSendTestInstructionTimeOutTowardsFenix() {

	var parentAddress gRPCClientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to transmit engine
		testInstructionTimeOutMessageToBeForwarded := <-gatewayObject.testInstructionTimeOutMessageChannel

		gatewayObject.logger.WithFields(logrus.Fields{
			"ID": "d1d4385b-b7c1-473d-8105-2e9b2341ef14",
			"testInstructionTimeOutMessageToBeForwarded": testInstructionTimeOutMessageToBeForwarded,
		}).Debug("Received a new 'testInstructionTimeOutMessageToBeForwarded' from channel that shoud be forwarded")

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
		gatewayObject.dbMessageQueue <- dbMessage

		// Wait for address from channel, then close the channel
		parentAddressByteArray := <-returnParentAddressChannel
		close(returnParentAddressChannel)

		// Convert saved json object into Go-struct
		err = json.Unmarshal(parentAddressByteArray.value, &parentAddress)
		if err != nil {
			// Problem with unmarshal the json object
			gatewayObject.logger.WithFields(logrus.Fields{
				"ID":                      "dd9f143f-9796-4651-a394-b72f5375cb9e",
				"parentAddressByteArray,": parentAddressByteArray,
			}).Error("Can't unmarshal gRPCClients address object from database")
			//TODO Send Error information to Fenix
		} else {
			// Send TestInstruction to client using gRPC-call
			addressToDial := parentAddress.clientIp + parentAddress.clientPort

			// Set up connection to Parent Gateway or Fenix
			remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "6118ae80-1d86-4530-83e4-549c00d01337",
					"addressToDial": addressToDial,
					"error message": err,
				}).Warning("Did not connect to Child (Gateway or Plugin) Server!")
				// TODO Send Error information to Fenix
				// TODO Add message to memmory cash for later resend
				// TODO Save message in localDB for later resend
			} else {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "e1c1903e-d72b-4ca2-973e-33b8525cb6ee",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to Parent-gateway- or Fenix-Server!")

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

				// Do gRPC-call to client gateway or Fenix
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestInstructionTimeOutTowardsFenix(ctx, &testInstructionTimeOutMessageToBeForwarded)
				if err != nil {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "1ae9d406-b8fc-4622-b6a3-8b2c0ce3cdc9",
						"returnMessage": returnMessage,
						"error":         err,
					}).Warning("Problem to send 'testInstructionTimeOutMessageToBeForwarded' to parent-Gateway or Fenix")
					// TODO Send Error information to Fenix
					// TODO Add message to memmory cash for later resend
					// TODO Save message in localDB for later resend
				} else {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "985c3a8b-dd01-496a-b3b8-9e1c67b89dd6",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send OK of 'informationMessageToBeForwarded' to Parent-Gateway or Fenix")

					// TODO Check for messages to Resend (If so then put them on channel)

				}
			}
		}
	}
}
