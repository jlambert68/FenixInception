package FenixGatewayServer

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func CallBackRegisterAvailbleTestInstructions(availbleTestInstructionAtPluginMessage *gRPC.AvailbleTestInstructionAtPluginMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                                     "1f4b45a2-5620-47a3-a8f0-52a759f424fb",
		"availbleTestInstructionAtPluginMessage": availbleTestInstructionAtPluginMessage,
	}).Debug("Incoming function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Save availbleTestInstructionAtPluginMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveAvailbleTestInstructionAtPluginMessageInDB((availbleTestInstructionAtPluginMessage))
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "827c39fd-e8d6-4ea2-a90d-5a1194f551e2",
		}).Debug("'availbleTestInstructionAtPluginMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'availbleTestInstructionAtPluginMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "e667e9e6-92dc-432e-ad9b-e1c9cb2dbbb5",
		}).Error("'availbleTestInstructionAtPluginMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'availbleTestInstructionAtPluginMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "7d2af962-6f59-41cf-9f26-af08f62413af",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func CallBackRegistrateAailableTestContainers(availbleTestContainersAtPluginMessage *gRPC.AvailbleTestContainersAtPluginMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                                    "63691666-a66d-419b-8328-b49a549070cb",
		"availbleTestContainersAtPluginMessage": availbleTestContainersAtPluginMessage,
	}).Debug("Incoming function CallBack: 'CallBackRegistrateAailableTestContainers'")

	// Save availbleTestContainersAtPluginMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveavailbleTestContainersAtPluginMessageInDB(availbleTestContainersAtPluginMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "a9c9ec6d-fa01-47c9-9f33-f35b136defd8",
		}).Debug("'availbleTestContainersAtPluginMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'availbleTestContainersAtPluginMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "fcf90b0c-6bf9-4bfe-b282-06ae761d832d",
		}).Error("'availbleTestContainersAtPluginMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'availbleTestContainersAtPluginMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "2b19cee5-419a-4281-be52-cd18586541f9",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func CallBackRegistrateAvailableTestDataDomains(supportedTestDataDomainsMessage *gRPC.SupportedTestDataDomainsMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                              "d3672229-2897-4654-83aa-0e56dd1db4b1",
		"supportedTestDataDomainsMessage": supportedTestDataDomainsMessage,
	}).Debug("Incoming function CallBack: 'CallBackRegistrateAvailableTestDataDomains'")

	// Save supportedTestDataDomainsMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveSupportedTestDataDomainsMessageInDB(supportedTestDataDomainsMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "2f013533-71ff-4ad8-9b13-3e9808a4f37b",
		}).Debug("'supportedTestDataDomainsMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'supportedTestDataDomainsMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "20befb0d-cd64-4712-9cfa-a429e83c3e63",
		}).Error("'supportedTestDataDomainsMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'supportedTestDataDomainsMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "77f8516a-7d27-4740-af2f-96a3307950d2",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

//TODO MessageToFenix

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func CallBackSendMessageToFenix(informationMessage *gRPC.InformationMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                 "cd1c1c36-6677-45de-9b73-e6f7e2238d1f",
		"informationMessage": informationMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendMessageToFenix'")

	// Save informationMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveInformationMessageInDB(informationMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "edfc3792-aa4e-4b5b-b048-1f36bcc74a2f",
		}).Debug("'informationMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'informationMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "45d0fe22-a7d7-467b-9cac-be67c7c0cb54",
		}).Error("'informationMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'informationMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "6001dd10-06e4-4998-8ccf-7c9acc75e431",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a InfoMessage toward Fenix
//
func CallBackSendTestInstructionTimeOutTowardsFenix(testInstructionTimeOutMessage *gRPC.TestInstructionTimeOutMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                            "e8eced61-8131-4b4b-ae4c-9d7dad50d48f",
		"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendTestInstructionTimeOutTowardsFenix'")

	// Save testInstructionTimeOutMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveTestInstructionTimeOutMessageInDB(testInstructionTimeOutMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "353a3773-47b5-4eea-baba-59beb44b4320",
		}).Debug("'testInstructionTimeOutMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testInstructionTimeOutMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "f53e96cb-1283-4102-bdd7-a084db7dc953",
		}).Error("'testInstructionTimeOutMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testInstructionTimeOutMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "f8f33004-be43-443f-bf36-5f945a830a60",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a TestExecutionLogMessage toward Fenix
//
func CallBackSendTestExecutionLogTowardsFenix(testExecutionLogMessage *gRPC.TestExecutionLogMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                            "8f8e11cb-8de5-42c4-8c65-ee79c04c9207",
		"testInstructionTimeOutMessage": testExecutionLogMessage,
	}).Debug("Incoming function CallBack: 'CallBackTestExecutionLogMessage'")

	// Save testExecutionLogMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveTestExecutionLogMessageInDB(testExecutionLogMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "410e1022-a56a-4b15-8b60-3468776b400a",
		}).Debug("'testExecutionLogMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testExecutionLogMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "3ba16e81-8304-4b5e-b6fb-1fb837d61307",
		}).Error("'testExecutionLogMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testExecutionLogMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "8a37b8e2-9ad3-425b-be71-515ffa335d54",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a supportedTestDataDomainsWithHeadersMessage toward Fenix
//
func CallBackSupportedTestDataDomains(supportedTestDataDomainsWithHeadersMessage *gRPC.SupportedTestDataDomainsWithHeadersMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID": "218c8382-a6df-4b0e-8bf4-13c2ce327dc0",
		"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
	}).Debug("Incoming function CallBack: 'CallBackSupportedTestDataDomains'")

	// Save testInstructionExecutionResultMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveSupportedTestDataDomainsWithHeadersMessageInDB(supportedTestDataDomainsWithHeadersMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "d5b46485-cfab-42f7-96da-2a90bf1532c3",
		}).Debug("'testInstructionExecutionResultMessage' was saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testInstructionExecutionResultMessage' was saved in Fenix database"
		returnMessage.Acknack = true
	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "24aecb02-b758-44e0-9a5d-498ed772acea",
		}).Error("'testInstructionExecutionResultMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'testInstructionExecutionResultMessage' was Not saved in Fenix database"
		returnMessage.Acknack = true
	}

	logger.WithFields(logrus.Fields{
		"ID": "b09a7d1d-328c-46fe-a83b-9d9e3e073559",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Return message back to Gateway
	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a testInstructionExecutionResultMessage toward Fenix
//
func CallBackSendTestInstructionResultTowardsFenix(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	var positivReturnMesage = "'testInstructionExecutionResultMessage' was saved in Fenix database"
	var negativReturnMesage = "'testInstructionExecutionResultMessage' could not be saved in Fenix database"
	var strangeErrorMessage = "When processing 'testInstructionExecutionResultMessage' an unknown error occured"

	var testInstructionsThatAreStillExecuting []string
	var err error

	logger.WithFields(logrus.Fields{
		"ID":                                    "0c22dd9e-d829-4a52-8468-ef6e02fa57a7",
		"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendTestInstructionResultTowardsFenix'")

	// Save testInstructionExecutionResultMessage in SQL-DB for further processing
	messageSavedWithoutProblem := saveTestInstructionExecutionResultMessageInDB(testInstructionExecutionResultMessage)
	if messageSavedWithoutProblem == true {

		// Message saved OK
		logger.WithFields(logrus.Fields{
			"ID": "d4a595c4-9622-4bc5-82f9-9711a94ba055",
		}).Debug("'testInstructionExecutionResultMessage' was saved in Fenix database")

		// Get all Testinstructions that are a peer to this TestInstruction, and can be run in parallell, and are still executing
		testInstructionsThatAreStillExecuting, err = listPeerTestInstructionPeersWhichIsExecuting(testInstructionExecutionResultMessage.PeerId)

		// Trigger next TestInstructions that is waiting to be executed if all current peers are finished
		if err == nil && len(testInstructionsThatAreStillExecuting) == 0 {
			testInstructionPeersThatShouldBeExecutedNext, err := listNextPeersToBeExecuted(testInstructionExecutionResultMessage.PeerId)
			if err == nil && len(testInstructionsThatAreStillExecuting) > 0 {
				err = triggerSendNextPeersForExecution(testInstructionPeersThatShouldBeExecutedNext)
			}
		}
		if err == nil {
			// OK
			// Create message back to parent Gateway/Plugin
			returnMessage.Comments = positivReturnMesage
			returnMessage.Acknack = true
		} else {
			//Not OK
			returnMessage.Comments = strangeErrorMessage
			returnMessage.Acknack = true
		}

	} else {

		// Message not saved OK
		logger.WithFields(logrus.Fields{
			"ID": "a681220a-4f42-4602-beaa-84954826a809",
		}).Error("'testInstructionExecutionResultMessage' was Not saved in Fenix database")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = negativReturnMesage
		returnMessage.Acknack = false
	}

	logger.WithFields(logrus.Fields{
		"ID": "3b816bbc-e36f-4612-a62f-83751d1d9fe0",
	}).Debug("Leaving function CallBack: 'CallBackSendTestInstructionResultTowardsFenix'")

	// Return message back to Gateway
	return returnMessage, nil
}
