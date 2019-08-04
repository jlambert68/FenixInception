package TestExecutionGateway

import (
	"encoding/json"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix for incoming TestInstructions that should be sent towards Plugin
//
func (gRPCServerTowardsPlugin *common_code.GRPCServerTowardsPluginStruct) SendTestInstructionTowardsPlugin(ctx context.Context, testInstruction *gRPC.TestInstruction_RT) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":              "b80b7746-04a8-4ef5-bdc5-2f9e07de754e",
		"testInstruction": *testInstruction,
	}).Debug("Incoming gRPC: 'SendTestInstructionTowardsPlugin'")

	// Set New Database Status and Timestamp on TestInstruction
	testInstruction.Metadata.DbStatusEnum = gRPC.DbStatusEnum_TESTINSTRUCTION_RECEIVED_FROM_PARENT
	testInstruction.Metadata.LastUpdatedDateTimeInDB = time.Now().String()

	// Convert TestInstruction struct into a byte array
	testInstructionByteArray, err := json.Marshal(*testInstruction)
	if err != nil {

		// Error when Unmarshaling to []byte
		LogErrorAndSendInfoToFenix(
			"a1d09308-7766-4f61-8030-46366dae3666",
			gRPC.InformationMessage_ERROR,
			"testInstruction",
			testInstruction.String(),
			err.Error(),
			"Error when converting 'getSupportedTestDomainsToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
		)
		returnMessage.Comments = "Error when converting incoming TestInstruction_RT into a byte array, stopping futher processing of this TestInstruction"
		returnMessage.Acknack = false
		return returnMessage, nil

	}

	// Save TestInstruction to local database, using local channel
	saveOK := SaveMessageToLocalDB(
		testInstruction.TestInstructionGuid,
		testInstructionByteArray,
		common_code.BucketForResendTestInstructionTowardsPlugin,
		"275bebe1-cb7e-4790-b2b5-070acedfff9a",
	)

	if saveOK == false {
		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Error when saving to local DB, stopping futher processing of this TestInstruction"
		returnMessage.Acknack = false
		return returnMessage, nil
	} else {

		logger.WithFields(logrus.Fields{
			"ID":              "5a98f0f1-5de9-4dcc-af76-f2888aaebf76",
			"testInstruction": testInstruction,
		}).Debug("TestInstructions was saved in local database")

		// Put TestInstruction on queue for further processing
		gatewayChannelPackage.TestInstructionMessageChannelTowardsPlugin <- testInstruction

		logger.WithFields(logrus.Fields{
			"ID":              "73e44541-c793-4ccd-8bc8-c94320f49f29",
			"testInstruction": testInstruction,
		}).Debug("Leaving: 'SendTestInstructionTowardsPlugin'")

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "OK"
		returnMessage.Acknack = true

		return returnMessage, nil
	}
}
