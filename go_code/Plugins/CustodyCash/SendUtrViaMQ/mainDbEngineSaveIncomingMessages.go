package SendUtrViaMQ

import (
	"encoding/json"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/jlambert68/FenixInception/go_code/common_code/pluginDB"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	//"jlambert/FenixInception2/go_code/common_code"
	gRPC_DB "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
)

// **********************************************************************************************************
// Save incoming 'SupportedTestDataDomainsRequest' to Main Database for Fenix Inception
//

/*
func saveSupportedTestDataDomainsMessageInDB(supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (messageSavedInDB bool) {
	var messageToBeSavedInDbAsByteArray []byte
	var err error

	messageSavedInDB = true
	supportedTestDataDomainsRequest.
	// Convert message into byteArray
	messageToBeSavedInDbAsByteArray, err = json.Marshal(*supportedTestDataDomainsRequest)

	// Check if there was an error when converting into byteArray
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "08af30bd-0641-44d8-93be-aefd8954d671",

		}).Error("'testInstructionMessage' was saved in Plugin database")
	}


	// Return if the messages was saved or not in database
	return messageSavedInDB

}
*/

// **********************************************************************************************************
// Save incoming 'TestInstruction_RT' to Main Database for Fenix Inception
//
func saveTestInstructionRTMessageInDB(testInstructionRT *gRPC.TestInstruction_RT) (messageSavedInDB bool) {

	var messageTBeSavedInDB gRPC_DB.WriteKeyValueMessage
	var messageToBeSavedInDbAsByteArray []byte
	var err error

	messageSavedInDB = true

	// Convert message into byteArray
	messageToBeSavedInDbAsByteArray, err = json.Marshal(*testInstructionRT)

	// Check if there was an error when converting into byteArray
	if err != nil {
		// Error when converting into byteArray
		logger.WithFields(logrus.Fields{
			"ID": "0a4d828f-c898-4fb1-818e-6e7ef512f527",
		}).Error("'testInstructionRT' could not be converted into an byteArray")

		// Set return code
		messageSavedInDB = false
	} else {
		// OK when converting into byteArray

		// Prepare data message to be sent to KeyValue-DB
		messageTBeSavedInDB = gRPC_DB.WriteKeyValueMessage{
			Key:            testInstructionRT.TestInstructionGuid,
			Bucket:         pluginDB.DbBucketForTestInstructionToBeExecutedMEssages,
			ValueSaveType:  gRPC_DB.ValueSaveTypeEnum_BYTEARRAY,
			Value:          messageToBeSavedInDbAsByteArray,
			ValueString:    "",
			CurrentVersion: gRPC_DB.CurrentVersionEnum(gRPC_DB.CurrentVersionEnum_value[pluginDB.GetHighestPluginGRPCVersion()]),
		}

		// Send message, to be sent to KeyValue-DB, via channel
		messageWasSavedInDB := SaveMessageUsingGRPC(messageTBeSavedInDB)

		// Check if message was saved in KeyValue-DB
		if messageWasSavedInDB == false {
			logger.WithFields(logrus.Fields{
				"ID": "bd589b5c-95da-4884-8bfa-9fe487d57c1f",
			}).Error("Message couldn't be saved in KeyValue-DB")

			// Set return code
			messageSavedInDB = false
		}
	}

	// Return that the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Call KeyValue-DB using gRPC to save message
//
func SaveMessageUsingGRPC(messageTBeSavedInDB gRPC_DB.WriteKeyValueMessage) (messageSavedInDB bool) {

	var (
		addressToDial                  string
		remotePluginDbServerConnection *grpc.ClientConn
		gRpcClientTowardPluginDB       gRPC.GatewayTowayPluginClient

		gRpcContexType context.Context
		dialSuccess    = false
	)
	var gRPCerr error
	var err error

	messageSavedInDB = true

	// Use background call in gRPC-call
	gRpcContexType = context.Background()

	// Get address to next child gateway/plugin
	addressToDial, err = getPluginDBAddressAndPort()

	if err != nil {
	// Couldn't get Address to PluginDB
		messageSavedInDB = false

		logger.WithFields(logrus.Fields{
			"ID": "6b3d5eb4-8bcf-43ce-a3d7-0b6fb8578ea3",
		}).Error("Couldn't get gRPC-address to KeyValue-DB")
	} else {
		// Got gRPC-address to PluginDB

		// Set up connection to Plugin
		remotePluginDbServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())

		// Check if the dial was a success
		if err != nil {
			// Could not dial PluginDB
		} else {
			// Dial was successful towards PluginDB

			// Creates a new gateway gRPC-Client towards Plugin
			gRpcClientTowardPluginDB = gRPC.NewGatewayTowayPluginClient(remotePluginDbServerConnection)

			messageSavedInDB =
			return messageSavedInDB
		}
	}
}
