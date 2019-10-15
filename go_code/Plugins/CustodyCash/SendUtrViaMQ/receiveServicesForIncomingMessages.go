package SendUtrViaMQ

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsRequest towards lugin
//
func CallBackSupportedTestDataDomainsRequestTowardsPluginType(supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                              "02925a48-cd16-48b4-8a95-a33556b85fb7",
		"supportedTestDataDomainsRequest": supportedTestDataDomainsRequest,
	}).Debug("Incoming function CallBack: 'CallBackSupportedTestDataDomainsRequestTowardsPluginType'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "'supportedTestDataDomainsRequest' was processed"
	returnMessage.Acknack = true

	// Start up processing of SupportedTestDataDomainsRequest as a goroutine
	go newIncomingSupportedTestDataDomainsRequest()

	logger.WithFields(logrus.Fields{
		"ID": "41d65996-3ba9-4108-ab62-4482fb2a1989",
	}).Debug("Leaving function CallBack: 'CallBackSupportedTestDataDomainsRequestTowardsPluginType'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a TestInstruction_RT_Message towards lugin
//
func CallBackSendTestInstructionTowardsPluginType(testInstructionMessage *gRPC.TestInstruction_RT) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                     "4d130dc3-761f-4fbe-84b7-6041d9107522",
		"testInstructionMessage": testInstructionMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendTestInstructionTowardsPluginType'")

	// Save testInstructionMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveTestInstructionRTMessageInDB(testInstructionMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "994804b2-642c-429b-966d-1efda98484d8",
		}).Debug("'testInstructionMessage' was saved in Plugin database")

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "'testInstructionMessage' was saved in Plugin database"
		returnMessage.Acknack = true

		// Start up processing of TestInstructionMessage as a goroutine
		go newIncomingTestInstructionMessage()

	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "2b73b378-314a-49ee-8b8f-2e3c001ed199",
		}).Error("'testInstructionMessage' was Not saved in Plugin database")

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "'testInstructionMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "9bd12e90-7607-4e82-9936-663ec2a43663",
	}).Debug("Leaving function CallBack: 'CallBackSendTestInstructionTowardsPluginType'")

	// Return message back to Gateway
	return returnMessage, nil
}
