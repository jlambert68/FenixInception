package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"time"
)

var returnChannel chan dbResultMessage_struct

// ********************************************************************************************
// Initiate Dispatch Engine for TestInstructions to be executed
//

func (gatewayObject *GatewayTowardsPluginObject_struct) dispatchTestInstruction() {

	// Create return Channel
	returnChannel = make(chan dbResultMessage_struct)

	// Start Dispatch Engine, for TestInstruction to be executed, as a go-routine
	go gatewayObject.dispatchTestInstructionEngineReceive()
}

// ********************************************************************************************
// Forward TestInstructions from incoming channel to child gateway or plugin
//

func (gatewayObject *GatewayTowardsPluginObject_struct) dispatchTestInstructionEngineReceive() {

	var clientAddress clientAddress_struct
	var err error

	for {
		// Wait for data comes from channel to dispatch engine
		testInstructionToBeForwarded := <-gatewayObject.testInstructionMessageQueue

		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":                           "f015de87-d51f-4025-8e27-cf818bbe255d",
			"testInstructionToBeForwarded": testInstructionToBeForwarded,
		}).Info("Received a new TestInstruction from channel that shoud be forwarded")

		// Get Clients address
		dbMessage := dbMessage_struct{
			DB_READ,
			"Clients",
			testInstructionToBeForwarded.PluginId,
			nil,
			returnChannel}

		// Send Read message to database to receive address
		gatewayObject.dbMessageQueue <- dbMessage

		// Wait for address from channel, then close the channel
		returnClientAddressChannel := make(chan dbResultMessage_struct)
		clientAddressByteArray := <-returnClientAddressChannel
		close(returnClientAddressChannel)

		// Convert saved json object into Go-struct
		err = json.Unmarshal(clientAddressByteArray.value, &clientAddress)
		if err != nil {
			// Problem with unmarshal the json object
			gatewayObject.logger.WithFields(logrus.Fields{
				"ID":                                     "cfa5245e-1c30-42b2-b8b6-94f0633fda02",
				"testInstructionToBeForwarded.PluginId,": testInstructionToBeForwarded.PluginId,
			}).Error("Can't unmarshal Client address object from Database")
			//TODO Send Error information to Fenix
		} else {
			// Send TestInstruction to client using gRPC-call
			addressToDial := clientAddress.clientIp + clientAddress.clientPort

			// Set up connection to Client Gateway or Plugin
			remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "3ad75b4a-4c7d-44eb-b781-c8cfd18c4fcf",
					"addressToDial": addressToDial,
					"error message": err,
				}).Error("Did not connect to Child (Gateway or Plugin) Server!")
				//TODO Send Error information to Fenix
			} else {
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":            "4510ae97-1753-4d4c-a7f5-3987054bd969",
					"addressToDial": addressToDial,
				}).Info("gRPC connection OK to Worker Server!")

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

				// Do gRPC-call to client gateway or Plugin
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestInstructionTowardsPlugin(ctx, &testInstructionToBeForwarded)
				if err != nil {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "4a0f0753-6394-4115-aa56-cf26c9354f83",
						"returnMessage": returnMessage,
						"error":         err,
					}).Error("Problem to send TestInstruction to child-Gateway or Plugin")
					//TODO Send Error information to Fenix
				} else {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":            "ef7f99b4-dcc0-4583-b16a-de3bf093880f",
						"addressToDial": addressToDial,
					}).Info("gRPC-send TestInstruction to child-Gateway or Plugin")

					// Set New Database Status and Timestamp on TestInstruction
					testInstructionToBeForwarded.Metadata.DbStatusEnum = gRPC.DbStatusEnum_TESTINSTRUCTION_SENT_TO_CHILD
					testInstructionToBeForwarded.Metadata.LastUpdatedDateTimeInDB = time.Now().String()

					// Convert TestInstruction struct into a byte array
					testInstructionByteArray, err := json.Marshal(testInstructionToBeForwarded)
					if err != nil {
						gatewayObject.logger.WithFields(logrus.Fields{
							"ID":              "b6cbe063-65a1-4d79-90f1-9849fb830e16",
							"testInstruction": testInstructionToBeForwarded,
							"err":             err,
						}).Error("Error when converting TestInstruction into a byte array, stopping futher processing of this TestInstruction")

						//TODO Send Error information to Fenix

					} else {

						// Save TestInstruction to local database, using local channel
						// Return Channel
						returnChannel := make(chan dbResultMessage_struct)

						dbMessage := dbMessage_struct{
							DB_WRITE,
							"TestInstructions",
							testInstructionToBeForwarded.TestInstructionGuid,
							testInstructionByteArray,
							returnChannel}

						// Send message to Database
						gatewayObject.dbMessageQueue <- dbMessage

						// Wait for result on result channel and then close returnChannel
						returnDBMessage := <-returnChannel
						close(returnChannel)

						if returnDBMessage.err != nil {
							gatewayObject.logger.WithFields(logrus.Fields{
								"ID":  "69e2f6bd-5ca3-4292-893d-075ff5ac79e0",
								"err": err,
							}).Error("Got an error when Saveing to local DB")

							//TODO Send Error information to Fenix

						}

					}
				}
			}
		}
	}
}
