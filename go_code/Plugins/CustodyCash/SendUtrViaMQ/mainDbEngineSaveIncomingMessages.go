package SendUtrViaMQ

import (
	"encoding/json"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	//"jlambert/FenixInception2/go_code/common_code"
	gRPC_DB "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
)

// **********************************************************************************************************
// Save incoming 'SupportedTestDataDomainsRequest' to Main Database for Fenix Inception
//
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

// **********************************************************************************************************
// Save incoming 'TestInstruction_RT' to Main Database for Fenix Inception
//
func saveTestInstructionRTMessageInDB(testInstructionRT *gRPC.TestInstruction_RT) (messageSavedInDB bool) {

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

		}).Error("'testInstructionRT' could not be converted into an byeArray")
	} else {
		// OK when converting into byteArray

		// TODO Prepare message to be sent to KeyValue-DB
		messageTBeSavedInDB := & gRPC_DB.WriteKeyValueMessage{
			testInstructionRT.TestInstructionGuid,
			common_code.
			ValueSaveType
			Value
			ValueString
			CurrentVersion       CurrentVersionEnum
		}

		// Send message, to be sent to KeyValue-DB, via channel
		dbSaveMessageChannel <- messageTBeSavedInDB
	}


	// Return if the messages was saved or not in database
	return messageSavedInDB
}
