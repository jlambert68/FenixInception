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

	// Save availbleTestInstruction-AtPluginMessage in SQL-DB for further processing
	//TODO använda channel till SQL eller functions-anrop?
	gatewayChannelPackage.AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestInstructionAtPluginMessage
	logger.WithFields(logrus.Fields{
		"ID": "827c39fd-e8d6-4ea2-a90d-5a1194f551e2",
	}).Debug("'availbleTestInstructionAtPluginMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "7d2af962-6f59-41cf-9f26-af08f62413af",
	}).Debug("Leaving function CallBack: 'CallBackRegisterAvailbleTestInstructions'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'availbleTestInstructionAtPluginMessage' was saved in Fenix database"
	returnMessage.Acknack = true

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

	// Put testExecutionLogMessage on queue for further processing
	gatewayChannelPackage.AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestContainersAtPluginMessage
	logger.WithFields(logrus.Fields{
		"ID": "d7c8ae8a-713e-4502-b3ce-e3427ee5d7d2",
	}).Debug("'SupportedTestDataDomainsMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "0c721a16-7cb3-415f-8e33-f440d715a65f",
	}).Debug("Leaving function CallBack: 'CallBackRegistrateAailableTestContainers'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'availbleTestContainersAtPluginMessage' was saved in Fenix database"
	returnMessage.Acknack = true

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

	// Put testExecutionLogMessage on queue for further processing
	gatewayChannelPackage.SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsMessage
	logger.WithFields(logrus.Fields{
		"ID": "4d55f441-4e13-41ae-a9eb-e71f288a18e9",
	}).Debug("'SupportedTestDataDomainsMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "fb223a31-88d9-4876-964f-a21d1cdb0a2c",
	}).Debug("Leaving function CallBack: 'CallBackRegistrateAvailableTestDataDomains'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'SupportedTestDataDomainsMessage' was saved in Fenix database"
	returnMessage.Acknack = true

	return returnMessage, nil
}

//TODO MessageToFenix

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a InfoMessage toward Fenix
//
func CallBackSendTestInstructionTimeOutTowardsFenix(testInstructionTimeOutMessage *gRPC.TestInstructionTimeOutMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                            "e8eced61-8131-4b4b-ae4c-9d7dad50d48f",
		"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendTestInstructionTimeOutTowardsFenix'")

	// Put testInstructionTimeOutMessage on queue for further processing
	gatewayChannelPackage.TestInstructionTimeOutMessageChannelTowardsFenix <- testInstructionTimeOutMessage
	logger.WithFields(logrus.Fields{
		"ID": "8b25dfbe-3105-4991-a61b-c5743cb3f49c",
	}).Debug("'testInstructionTimeOutMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "df7ebfce-badc-4a0b-bf29-268e75c3eed6",
	}).Debug("Leaving function CallBack: 'CallBackSendTestInstructionTimeOutTowardsFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "'TestInstructionTimeOutMessage' was saved in Fenix database"
	returnMessage.Acknack = true

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

	// Put testExecutionLogMessage on queue for further processing
	gatewayChannelPackage.TestExecutionLogMessageChannelTowardsFenix <- testExecutionLogMessage
	logger.WithFields(logrus.Fields{
		"ID": "b4925021-02aa-4ff8-b314-7a82fb2b19a0",
	}).Debug("'testExecutionLogMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "f2efe405-498f-42d6-9788-3b4397b322dc",
	}).Debug("Leaving function CallBack: 'CallBackSendTexecutionLogTowardsFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "'TestExecutionLogMessage' was saved in Fenix database"
	returnMessage.Acknack = true

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

	// Put supportedTestDataDomainsWithHeadersMessage on queue for further processing
	gatewayChannelPackage.SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsWithHeadersMessage
	logger.WithFields(logrus.Fields{
		"ID": "58b4643a-600c-48f5-93fe-3638f814863c",
	}).Debug("'supportedTestDataDomainsWithHeadersMessage' was saved in Fenix database")

	logger.WithFields(logrus.Fields{
		"ID": "956d8b7d-7ff4-41d2-b050-539805d1a87a",
	}).Debug("Leaving function CallBack: 'CallBackSupportedTestDataDomains'")

	// Create message back to child Gateway/Plugin
	returnMessage.Comments = "'supportedTestDataDomainsWithHeadersMessage' was saved in Fenix database"
	returnMessage.Acknack = true

	return returnMessage, nil
}

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a testInstructionExecutionResultMessage toward Fenix
//
func CallBackSendTestInstructionResultTowardsFenix(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                                    "d32b63f7-58fb-46f0-a04b-111c6af7b6fd",
		"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
	}).Debug("Incoming function CallBack: 'CallBackSendTestInstructionResultTowardsFenix'")

	// Put testInstructionExecutionResultMessage on queue for further processing
	gatewayChannelPackage.TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix <- testInstructionExecutionResultMessage
	logger.WithFields(logrus.Fields{
		"ID": "6894a11c-ba36-479e-8fdc-2391f7820fc7",
	}).Debug("'SupportedTestDataDomainsMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "1782261e-8bf2-4707-bd69-b2379f02e034",
	}).Debug("Leaving function CallBack: 'CallBackSendTestInstructionResultTowardsFenix'")

	// Create message back to child Gateway/Plugin
	returnMessage.Comments = "'testInstructionExecutionResultMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
