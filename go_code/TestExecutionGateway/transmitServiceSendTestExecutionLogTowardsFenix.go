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

	var parentAddress gRPCClientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to transmit engine
		testExecutionLogMessageToBeForwarded := <-gatewayObject.testExecutionLogMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":                                   "5c9fe63a-fd82-4ccd-8386-2b9c049e51a1",
			"testExecutionLogMessageToBeForwarded": testExecutionLogMessageToBeForwarded,
		}).Debug("Received a new 'testExecutionLogMessageToBeForwarded' from channel that shoud be forwarded")

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
				"ID":                      "58232255-c767-4143-95b3-3f52de741543",
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
					"ID":            "863612af-86b3-431c-8270-9d335f3f0001",
					"addressToDial": addressToDial,
					"error message": err,
				}).Warning("Did not connect to Child (Gateway or Plugin) Server!")
				// TODO Send Error information to Fenix
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
				testExecutionLogMessageToBeForwarded.SenderId = CallingSystemId
				testExecutionLogMessageToBeForwarded.SenderName = CallingSystemName

				// Do gRPC-call to client gateway or Fenix
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestExecutionLogTowardsFenix(ctx, testExecutionLogMessageToBeForwarded)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"ID":            "840312be-e773-4e9d-ae05-0ec8e9dcca5c",
						"returnMessage": returnMessage,
						"error":         err,
					}).Warning("Problem to send 'testExecutionLogMessageToBeForwarded' to parent-Gateway or Fenix")
					// TODO Send Error information to Fenix
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
}
