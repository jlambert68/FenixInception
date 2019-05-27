package TestExecutionGateway

import (
	"encoding/json"
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

	var parentAddress gRPCClientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to transmit engine
		informationMessageToBeForwarded := <-gatewayObject.informationMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":                              "3103a0a5-ff25-4e85-9939-5e8d72e26ba3",
			"informationMessageToBeForwarded": informationMessageToBeForwarded,
		}).Debug("Received a new informationMessage from channel that shoud be forwarded")

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
		dbMessageQueue <- dbMessage

		// Wait for address from channel, then close the channel
		parentAddressByteArray := <-returnParentAddressChannel
		close(returnParentAddressChannel)

		// Convert saved json object into Go-struct
		err = json.Unmarshal(parentAddressByteArray.value, &parentAddress)
		if err != nil {
			// Problem with unmarshal the json object
			logger.WithFields(logrus.Fields{
				"ID":                      "fba82bec-aa10-4f2c-8972-177872bb94c0",
				"parentAddressByteArray,": parentAddressByteArray,
			}).Error("Can't unmarshal gRPCClients address object from database")
			//TODO Send Error information to Fenix
		} else {
			// Send TestInstruction to client using gRPC-call
			addressToDial := parentAddress.clientIp + parentAddress.clientPort

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
				informationMessageToBeForwarded.SenderId = CallingSystemId
				informationMessageToBeForwarded.SenderName = CallingSystemName

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
}
