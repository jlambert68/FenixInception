package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix for incoming TestInstructions that should be sent towards Plugin
//
func (gatewayObject *GatewayTowardsPluginObject_struct) SendTestInstructionTowardsPlugin(ctx context.Context, testInstruction gRPC.TestInstruction_RT) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	gatewayObject.logger.WithFields(logrus.Fields{
		"ID":              "b80b7746-04a8-4ef5-bdc5-2f9e07de754e",
		"testInstruction": testInstruction,
	}).Info("Incoming: 'SendTestInstructionTowardsPlugin'")

	// Set New Database Status on TestInstruction
	testInstruction.Metadata.DbStatusEnum = gRPC.DbStatusEnum_TESTINSTRUCTION_RECEIVED_FROM_PARENT

	// Convert TestInstruction struct into a byte array
	testInstructionByteArray, err := json.Marshal(testInstruction)
	if err != nil {
		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":              "de52efd0-684a-48cf-bebb-d4b6f6c6f517",
			"testInstruction": testInstruction,
			"err":             err,
		}).Error("Error when converting TestInstruction_RT into a byte array, stopping futher processing of this TestInstruction")

		//TODO Send Error information to Fenix

		returnMessage.Comments = "Error when converting incoming TestInstruction_RT into a byte array, stopping futher processing of this TestInstruction"
		returnMessage.Acknack = false
		return returnMessage, nil

	}

	// Save TestInstruction to local database, using local channel
	// Return Channel
	returnChannel := make(chan dbResultMessage_struct)

	dbMessage := dbMessage_struct{
		DB_WRITE,
		"TestInstructions",
		testInstruction.TestInstructionGuid,
		testInstructionByteArray,
		returnChannel}

	// Send message to Database
	gatewayObject.dbMessageQueue <- dbMessage

	// Wait for result on result channel then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.err != nil {
		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":  "446cb366-9d84-4c1c-a628-ef162f7c1747",
			"err": err,
		}).Error("Got an error when Saveing to local DB")
		//TODO Send Error information to Fenix

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Error when saving to local DB, stopping futher processing of this TestInstruction"
		returnMessage.Acknack = false
		return returnMessage, nil
	}

	// Put TestInstruction on queue for further processing
	gatewayObject.testInstructionMessageQueue <- testInstruction

	gatewayObject.logger.WithFields(logrus.Fields{
		"ID":              "73e44541-c793-4ccd-8bc8-c94320f49f29",
		"testInstruction": testInstruction,
	}).Info("Leaving: 'SendTestInstructionTowardsPlugin'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "OK"
	returnMessage.Acknack = true

	return returnMessage, nil
}
