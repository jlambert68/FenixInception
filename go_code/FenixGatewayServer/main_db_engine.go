package FenixGatewayServer

import (
	"database/sql"
	"errors"
	"fmt"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"jlambert/FenixInception2/go_code/common_code"
)

var mainDB *sql.DB

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "fenixinception"
)

// **********************************************************************************************************
// Initiate connection to Main Database for Fenix Inception
//
func initiateMainDB() {

	logger.WithFields(logrus.Fields{
		"ID": "8be4a508-9584-4ff3-8d97-028591ece5d4",
	}).Debug("Initiate connection to Main Database for Fenix Inception")

	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	mainDB, err = sql.Open("postgres", dbinfo)

	// If some error occured then close the application
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "8be4a508-9584-4ff3-8d97-028591ece5d4",
			"err": err,
		}).Fatal("Error when initiateting connection to Main Database for Fenix Inception, closing Fenix Inception")
	}
}

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func saveTestInstructionExecutionResultMessageInDB(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (messageSavedInDB bool) {
	// TODO change from Update to Insert new copy row

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "UPDATE testcaseexecutions.Executing_TestInstructions "
	sqlToBeExecuted = sqlToBeExecuted + "SET TestInstructionExecutionFinished=$1, TestInstructionExecutionResult=$2, updatedDateTime=$3, ResponsePayload=$4 "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE TestInstructionId=$5"
	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                                    "a4c440ca-e7aa-431d-960c-b8d63e0fd15b",
			"err":                                   err,
			"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from testInstructionExecutionResultMessage")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		sqlResult, err := sqlStatement.Exec(true, testInstructionExecutionResultMessage.TestExecutionResult.String(), common_code.GeneraTimeStampUTC(), testInstructionExecutionResultMessage, testInstructionExecutionResultMessage.TestInstructionGuid)
		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                                    "b0e9a525-a296-4cf7-b8a8-ce5c4049f5b4",
				"err":                                   err,
				"sqlResult":                             sqlResult,
				"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
			}).Error("Error when updating Main Database with data from testInstructionExecutionResultMessage")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                                    "b0e9a525-a296-4cf7-b8a8-ce5c4049f5b4",
				"err":                                   err,
				"sqlResult":                             sqlResult,
				"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
			}).Debug("Fenix main Database was updeted with  data from testInstructionExecutionResultMessage")
		}

	}

	// Return if the messages was saved or not
	return messageSavedInDB
}

// **********************************************************************************************************
// List all peers to that are still execution
//
func listPeerTestInstructionPeersWhichIsExecuting(testInstructionPeerGuid string) (testInstructionPeersThatAreStillExecuting []string, err error) {

	// Prepare SQL
	var sqlToBeExecuted = "SELECT TestInstructionId "
	sqlToBeExecuted = sqlToBeExecuted + "FROM testcaseexecutions.Executing_TestInstructions "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE TestInstructionExecutionFinished isnull AND "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionPeerId = $1 "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Error while preparing SQL
		logger.WithFields(logrus.Fields{
			"ID":                      "f9ba926c-8701-45bc-9e62-694a48188325",
			"err":                     err,
			"testInstructionPeerGuid": testInstructionPeerGuid,
			"sqlToBeExecuted":         sqlToBeExecuted,
		}).Error("Error when Praparing SQL for listing peers to newly finished TestInstruction, which is still executing")

	} else {
		// Preparing SQL was OK
		// Now Execute SQL in DB
		sqlResult, err := sqlStatement.Query(testInstructionPeerGuid)
		if err != nil {
			// Error while executing the SQL
			logger.WithFields(logrus.Fields{
				"ID":        "b0e9a525-a296-4cf7-b8a8-ce5c4049f5b4",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when looking for Test Instructions-peers which has not been finished executing")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "42b5135c-876c-4e0b-9d42-d8c02f612933",
				"sqlResult": sqlResult,
			}).Debug("Found Test Instruction-peers in Fenix main Database")

			// Convert SQL-result into correct return message for function
			var testInstructionId string
			for sqlResult.Next() {
				err = sqlResult.Scan(&testInstructionId)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"ID":        "98220887-e01c-4b78-9861-7de75f9a4675",
						"err":       err,
						"sqlResult": sqlResult,
					}).Error("Error when processing SQL-result for Test Instructions-peers which has not been finished executing")
				} else {
					//Add PeerId to array if peers
					testInstructionPeersThatAreStillExecuting = append(testInstructionPeersThatAreStillExecuting, testInstructionId)
				}

			}
		}
	}

	// Return if the messages was saved or not
	return testInstructionPeersThatAreStillExecuting, err
}

// **********************************************************************************************************
// List all peers which are the next ones to be executed
//
func listNextPeersToBeExecuted(testInstructionPeerGuid string) (testInstructionPeersThatShouldBeExecutedNext []string, err error) {

	// Prepare SQL
	var sqlToBeExecuted = "SELECT TestInstructionId "
	sqlToBeExecuted = sqlToBeExecuted + "FROM testcaseexecutions.Executing_TestInstructions "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE TestInstructionExecutionFinished isnull AND "
	sqlToBeExecuted = sqlToBeExecuted + "PreviousTestInstructionPeerId = $1 "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Error while preparing SQL
		logger.WithFields(logrus.Fields{
			"ID":                      "497d1bc6-c76c-4300-8f63-8533112a13f8",
			"err":                     err,
			"testInstructionPeerGuid": testInstructionPeerGuid,
			"sqlToBeExecuted":         sqlToBeExecuted,
		}).Error("Error when Praparing SQL for listing next peers of TestInstructions to be executed")

	} else {
		// Preparing SQL was OK
		// Now Execute SQL in DB
		sqlResult, err := sqlStatement.Query(testInstructionPeerGuid)
		if err != nil {
			// Error while executing SQL
			logger.WithFields(logrus.Fields{
				"ID":        "1a6dfa6a-de5c-4955-ae7e-084266377169",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when looking for next TestInstructions-peers which should be executed")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "0bc01a85-617a-4d59-a8cd-ffe1078831b7",
				"sqlResult": sqlResult,
			}).Debug("sqlResult for Test Instruction-peers to be executed in Fenix main Database")

			// Convert SQL-result into correct return message for function
			var testInstructionId string
			for sqlResult.Next() {
				err = sqlResult.Scan(&testInstructionId)
				if err != nil {
					// Error while looping through result
					logger.WithFields(logrus.Fields{
						"ID":        "a9ed6228-2bd6-4ed8-bd06-789740dd1665",
						"err":       err,
						"sqlResult": sqlResult,
					}).Error("Error when processing SQL-result for Test Instructions-peers which has not been finished executing")
				} else {
					//Add PeerId to array if peers
					testInstructionPeersThatShouldBeExecutedNext = append(testInstructionPeersThatShouldBeExecutedNext, testInstructionId)
				}

			}
		}
	}

	// Return if the messages was saved or not
	return testInstructionPeersThatShouldBeExecutedNext, err
}

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func saveInDbThatTestInstructionHasBeenSentToPlugin(testInstructionId string) (messageSavedInDB bool) {
	// TODO change from Update to Insert new copy row

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "UPDATE testcaseexecutions.executing_testcasesteps "
	sqlToBeExecuted = sqlToBeExecuted + "SET TestInstructionSentDateTime=$1, updatedDateTime=$2 "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE TestInstructionId=$3"
	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Error while preparing SQL
		logger.WithFields(logrus.Fields{
			"ID":                "af48db0e-8351-4f6c-99c0-92273aae725c",
			"err":               err,
			"testInstructionId": testInstructionId,
		}).Error("Error when Praparing SQL for updating Main Database with data that TestInstruction was sent to Plugin")

		messageSavedInDB = false
	} else {
		// Preparing SQL was OK
		// Now Execute SQL in DB
		sqlResult, err := sqlStatement.Exec(common_code.GeneraTimeStampUTC(), testInstructionId)
		if err != nil {
			// Error while executing SQL
			logger.WithFields(logrus.Fields{
				"ID":                "800b9988-e550-4807-9109-10da5c1ef9ca",
				"err":               err,
				"sqlResult":         sqlResult,
				"testInstructionId": testInstructionId,
			}).Error("Error when updating Main Database with data that TestInstruction was sent to Plugin")

			messageSavedInDB = false
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":                "47ab1591-8e23-4b69-8e0d-fa1637261b82",
				"err":               err,
				"sqlResult":         sqlResult,
				"testInstructionId": testInstructionId,
			}).Debug("Fenix Main Database was updeted with data that TestInstruction was sent to Plugin")
		}
	}

	// Return if the messages was saved or not
	return messageSavedInDB
}

// **********************************************************************************************************
// Get TestInstruction-Payload to be sent to Plugin for Execution
//
func getTestInstructionPayloadToPlugin(testInstructionGuid string) (testInstructionPayloadToPlugin *gRPC.TestInstruction_RT, err error) {

	// Prepare SQL
	var sqlToBeExecuted = "SELECT Payload "
	sqlToBeExecuted = sqlToBeExecuted + "FROM testcaseexecutions.Executing_TestInstructions "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE TestInstructionExecutionFinished isnull AND "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionId = $1 "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Error while preparing SQL
		logger.WithFields(logrus.Fields{
			"ID":                  "dfbf3c87-f756-4f97-82d2-9a94e5a87ad2",
			"err":                 err,
			"testInstructionGuid": testInstructionGuid,
			"sqlToBeExecuted":     sqlToBeExecuted,
		}).Error("Error when Praparing SQL for getting TestInstructions-Payload to be sent for execution")

	} else {
		// Preparing SQL was OK
		// Now Execute SQL in DB
		sqlResult, err := sqlStatement.Query(testInstructionGuid)
		if err != nil {
			// Error while executing SQL
			logger.WithFields(logrus.Fields{
				"ID":        "",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when executing SQL for getting TestInstructions-Payload to be sent for execution")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "8102b4d5-c2ed-493e-918a-3365669d9cba",
				"sqlResult": sqlResult,
			}).Debug("sqlResult for TestInstruction-Payload for getting TestInstructions-Payload to be sent for execution")

			// Convert SQL-result into correct return message for function
			//var testInstructionId string
			var numberOfRowsInResult int
			numberOfRowsInResult = 0
			for sqlResult.Next() {
				err = sqlResult.Scan(&testInstructionPayloadToPlugin)
				if err != nil {
					// Error while looping through result
					logger.WithFields(logrus.Fields{
						"ID":        "f03c9fc8-672c-45d2-9874-b7f70f42687a",
						"err":       err,
						"sqlResult": sqlResult,
					}).Error("Error when processing SQL-result for for TestInstruction-Payload for getting TestInstructions-Payload to be sent for execution")
				} else {
					numberOfRowsInResult++
				}
			}

			// If the SQL gave a result that differs from 1 row then something is wrong
			if numberOfRowsInResult != 1 {
				// Error while looping through result
				logger.WithFields(logrus.Fields{
					"ID":                   "bf187ebf-ddbd-49f5-ac07-0be1e3251985",
					"numberOfRowsInResult": numberOfRowsInResult,
				}).Error("Error, expected exactly one row in result set")

				// Create an error
				err = errors.New("Error, expected exactly one row in result set but got " + string(numberOfRowsInResult))
			}
		}
	}

	// Return if the messages was saved or not
	return testInstructionPayloadToPlugin, err
}

// **********************************************************************************************************
// Save incoming 'SupportedTestDataDomainsMessage' to Main Database for Fenix Inception
//
func saveSupportedTestDataDomainsMessageInDB(supportedTestDataDomains *gRPC.SupportedTestDataDomainsMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	/*
	   create table testdatadomains.Supported_TestData_Domains
	   (
	       OriginalSenderId          uuid                     not null, -- The Id of the gateway/plugin that created the message 1
	       OriginalSenderName        varchar default null,              -- he name of the gateway/plugin that created the message 2
	       MessageId                 uuid                     not null, -- A unique id generated when sent from Plugin 3
	       TestDataDomainId          uuid                     not null, -- The unique id of the testdata domain 4
	       TestDataDomainName        varchar default null,              -- The name of the testdata domain 5
	       TestDataDomainDescription varchar default null,              -- A description of the testdata domain 6
	       TestDataDomainMouseOver   varchar default null,              -- A mouse over description of the testdata domain 7
	       OrginalCreateDateTime     timestamp with time zone not null, -- The timestamp when the orignal message was created 8
	       OriginalSystemDomainId    uuid                     not null, -- he Domain/system's Id where the Sender operates 9
	       OriginalSystemDomainName  varchar default null,              -- The Domain/system's Name where the Sender operates 10
	       updatedDateTime           timestamp with time zone not null  -- The Datetime when the row was created/updates 11
	   );
	*/

	var sqlToBeExecuted = "INSERT INTO testdatadomains.Supported_TestData_Domains "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, MessageId, TestDataDomainId, TestDataDomainName, TestDataDomainDescription, "
	sqlToBeExecuted = sqlToBeExecuted + "TestDataDomainMouseOver, OrginalCreateDateTime, OriginalSystemDomainName, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                       "7056bde2-a796-4dab-baa4-1b93dcc4c261",
			"err":                      err,
			"supportedTestDataDomains": supportedTestDataDomains,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'supportedTestDataDomains'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK

		// Get number of Test Data Domain-rows
		numberOfDomains := len(supportedTestDataDomains.TestDataDomains)
		if numberOfDomains > 0 {
			// Loop over all  TestData-domains
			for currentDomain := 0; currentDomain < numberOfDomains; currentDomain++ {

				// Values to insert into database
				sqlResult, err := sqlStatement.Exec(
					supportedTestDataDomains.OriginalSenderId,
					supportedTestDataDomains.OriginalSenderName,
					supportedTestDataDomains.MessageId,
					supportedTestDataDomains.TestDataDomains[currentDomain].TestDataDomainId,
					supportedTestDataDomains.TestDataDomains[currentDomain].TestDataDomainName,
					supportedTestDataDomains.TestDataDomains[currentDomain].TestDataDomainDescription,
					supportedTestDataDomains.TestDataDomains[currentDomain].TestDataDomainMouseOver,
					supportedTestDataDomains.OrginalCreateDateTime,
					supportedTestDataDomains.OriginalSystemDomainId,
					supportedTestDataDomains.OriginalSystemDomainName,
					common_code.GeneraTimeStampUTC())

				if err != nil {
					// Error while executing
					logger.WithFields(logrus.Fields{
						"ID":                       "8047c4c3-461f-42d8-87c6-5c2cb5f8e75b",
						"err":                      err,
						"sqlResult":                sqlResult,
						"supportedTestDataDomains": supportedTestDataDomains,
					}).Error("Error when updating Main Database with data from 'supportedTestDataDomains'")

					messageSavedInDB = false
				} else {
					//SQL executed OK
					logger.WithFields(logrus.Fields{
						"ID":                       "d456499c-2ad1-4677-8e1d-909a7ecab560",
						"err":                      err,
						"sqlResult":                sqlResult,
						"supportedTestDataDomains": supportedTestDataDomains,
					}).Debug("Fenix main Database was updated with data from 'supportedTestDataDomains'")
				}
			}
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'InformationMessage' to Main Database for Fenix Inception
//
func saveInformationMessageInDB(informationMessage *gRPC.InformationMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO messages.InformationMessages "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, OriginalSystemDomainId, OriginalSystemDomainName, MessageId, "
	sqlToBeExecuted = sqlToBeExecuted + "InformationMessageTypeId, InformationMessageTypeName, Message, OrginalCreateDateTime, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                 "626d1aae-2193-4aae-ae58-7df83b2aac3f",
			"err":                err,
			"informationMessage": informationMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'informationMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		/*
			create table messages.InformationMessages
			(
			    OriginalSenderId           uuid                     not null, -- The Id of the gateway/plugin that created the message 1
			    OriginalSenderName         varchar default null,              -- Thn name of the gateway/plugin that created the message 2
			    OriginalSystemDomainId     uuid                     not null, -- The Domain/system's Id where the Sender operates 3
			    OriginalSystemDomainName   varchar default null,              -- The Domain/system's Name where the Sender operates 4
			    MessageId                  uuid                     not null, -- A unique id generated when sent from Gateway or Plugin 5
			    InformationMessageTypeId   int                      not null, -- The Id of the information message type 6
			    InformationMessageTypeName varchar                  not null, -- The name of the information message type 7
			    Message                    varchar default null,              -- The message from Gateway or Plugin to Fenix 8
			    OrginalCreateDateTime      timestamp with time zone not null, -- The timestamp when the orignal message was created 9
			    updatedDateTime            timestamp with time zone not null  -- The Datetime when the row was created/updates 10
			);
		*/
		// Values to insert into database
		sqlResult, err := sqlStatement.Exec(
			informationMessage.OriginalSenderId,
			informationMessage.OriginalSenderName,
			informationMessage.OriginalSystemDomainId,
			informationMessage.OriginalSystemDomainName,
			informationMessage.MessageId,
			informationMessage.MessageType,
			gRPC.InformationMessage_InformationType_name[int32(informationMessage.MessageType)],
			informationMessage.Message,
			informationMessage.OrginalCreateDateTime,
			common_code.GeneraTimeStampUTC())

		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                 "a1aa006d-599d-4e23-97a0-5cf33ca79e52",
				"err":                err,
				"sqlResult":          sqlResult,
				"informationMessage": informationMessage,
			}).Error("Error when updating Main Database with data from 'informationMessage'")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                 "6321f830-e9d3-4185-b80f-501ec053de99",
				"err":                err,
				"sqlResult":          sqlResult,
				"informationMessage": informationMessage,
			}).Debug("Fenix main Database was updated with data from 'informationMessage'")
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'TestExecutionLogMessage' to Main Database for Fenix Inception
//
func saveTestExecutionLogMessageInDB(testExecutionLogMessage *gRPC.TestExecutionLogMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO messages.TestExecutionLogMessage "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, OriginalSystemDomainId, OriginalSystemDomainName, LogMessageId, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionGuid, TestInstructionLogCreatedDateTime, TestInstructionLogCreatedDateTime, TestInstructionLogSentDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "LogMessage, LogMessageTypeId, logmessageTypeName, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                      "9c711fac-20a4-4694-92c1-7d75591e11d2",
			"err":                     err,
			"testExecutionLogMessage": testExecutionLogMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'testExecutionLogMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		/*
			create table messages.logmessages
			(
			    OriginalSenderId                  uuid                     not null, -- The Id of the gateway/plugin that created the message 1
			    OriginalSenderName                varchar default null,              -- The name of the gateway/plugin that created the message 2
			    OriginalSystemDomainId            uuid                     not null, -- The Domain/system's Id where the Sender operates 3
			    OriginalSystemDomainName          varchar default null,              -- he Domain/system's Name where the Sender operates 4
			    LogMessageId                      uuid                     not null, -- A unique id created by plugin 5
			    TestInstructionGuid               uuid                     not null, -- TestInstructionGuid is a unique id created when TestInstruction is created in TestCase 6
			    TestInstructionLogCreatedDateTime timestamp with time zone not null, -- The DateTime when the message was created 7
			    TestInstructionLogSentDateTime    timestamp with time zone not null, -- he DateTime when the message was sent 8
			    LogMessage                        varchar default null,              -- The log message 9
			    LogMessageTypeId                  int                      not null, -- The Id of the log message type, text or json 10
			    logmessageTypeName                varcha                   not null, -- The name of the log message type, text or json 11
			    updatedDateTime                   timestamp with time zone not null  -- The Datetime when the row was created/updates 12
			);
		*/
		// Values to insert into database
		sqlResult, err := sqlStatement.Exec(
			testExecutionLogMessage.OriginalSenderId,
			testExecutionLogMessage.OriginalSenderName,
			testExecutionLogMessage.OriginalSystemDomainId,
			testExecutionLogMessage.OriginalSystemDomainName,
			testExecutionLogMessage.LogMessageId,
			testExecutionLogMessage.TestInstructionGuid,
			testExecutionLogMessage.TestInstructionResultCreatedDateTime,
			testExecutionLogMessage.TestInstructionResultSentDateTime,
			testExecutionLogMessage.LogMessage,
			testExecutionLogMessage.LogMessageType,
			gRPC.LogMessageTypeEnum_name[int32(testExecutionLogMessage.LogMessageType)],
			common_code.GeneraTimeStampUTC())

		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                      "6d56d36d-e953-4894-bb22-b82a0bca96cd",
				"err":                     err,
				"sqlResult":               sqlResult,
				"testExecutionLogMessage": testExecutionLogMessage,
			}).Error("Error when updating Main Database with data from 'testExecutionLogMessage'")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                      "d0621424-33ce-4e97-8519-80e7847a929c",
				"err":                     err,
				"sqlResult":               sqlResult,
				"testExecutionLogMessage": testExecutionLogMessage,
			}).Debug("Fenix main Database was updated with data from 'testExecutionLogMessage'")
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'TestInstructionTimeOutMessage' to Main Database for Fenix Inception
//
func saveTestInstructionTimeOutMessageInDB(testInstructionTimeOutMessage *gRPC.TestInstructionTimeOutMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO messages.TestExecutionLogMessage "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, OriginalSystemDomainId, OriginalSystemDomainName, MessageId, "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalMessageId, TimeOut, OrginalCreateDateTime, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                            "fbe9433f-6875-4a0b-9e34-495425c60242",
			"err":                           err,
			"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'testInstructionTimeOutMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		/*
			create table messages.TestInstructionTimeOutMessage
			(
			    OriginalSenderId         uuid                     not null, -- The Id of the gateway/plugin that created the message 1
			    OriginalSenderName       varchar default null,              -- The name of the gateway/plugin that created the message 2
			    OriginalSystemDomainId   uuid                     not null, -- The Domain/system's Id where the Sender operates 3
			    OriginalSystemDomainName varchar default null,              -- The Domain/system's Name where the Sender operates 4
			    MessageId                uuid                     not null, -- A unique id generated when sent from Plugi 5
			    OriginalMessageId        uuid                     not null, -- A unique id from the TestInstruction request sent from Fenix 6
			    TimeOut                  timestamp with time zone not null, -- The Timeout time when TestInstruction should Timeout if no answer has come back 7
			    OrginalCreateDateTime    timestamp with time zone not null, -- The timestamp when the orignal message was created 8
			    updatedDateTime          timestamp with time zone not null  -- The Datetime when the row was created/updates 9
			);
		*/
		// Values to insert into database
		sqlResult, err := sqlStatement.Exec(
			testInstructionTimeOutMessage.OriginalSenderId,
			testInstructionTimeOutMessage.OriginalSenderName,
			testInstructionTimeOutMessage.OriginalSystemDomainId,
			testInstructionTimeOutMessage.OriginalSystemDomainName,
			testInstructionTimeOutMessage.MessageId,
			testInstructionTimeOutMessage.OriginalMessageId,
			testInstructionTimeOutMessage.OrginalCreateDateTime,
			common_code.GeneraTimeStampUTC())

		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                            "53ce7ec9-8278-4fb9-99df-b04499552893",
				"err":                           err,
				"sqlResult":                     sqlResult,
				"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
			}).Error("Error when updating Main Database with data from 'testInstructionTimeOutMessage'")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                            "37dac8d8-5942-4727-8412-ed3582c28653",
				"err":                           err,
				"sqlResult":                     sqlResult,
				"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
			}).Debug("Fenix main Database was updated with data from 'testInstructionTimeOutMessage'")
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'SupportedTestDataDomainsWithHeadersMessage' to Main Database for Fenix Inception
//
func saveSupportedTestDataDomainsWithHeadersMessageInDB(supportedTestDataDomainsWithHeadersMessage *gRPC.SupportedTestDataDomainsWithHeadersMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO messages.TestExecutionLogMessage "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, MessageId, OriginalMessageId, "
	sqlToBeExecuted = sqlToBeExecuted + "TestDataDomainId, TestDataDomainName, TestDataDomainDescription, TestDataDomainMouseOver "
	sqlToBeExecuted = sqlToBeExecuted + "TestDataFilterHeaderId, TestDataFilterHeaderName, TestDataFilterHeaderDescription, TestDataFilterHeaderMouseOver "
	sqlToBeExecuted = sqlToBeExecuted + "TestDataFilterHash, AllowMultipleChoices, AllowNoChoice, "
	sqlToBeExecuted = sqlToBeExecuted + "TestDataFilterHeaderValueId, TestDataFilterHeaderValueName, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":  "bed75c92-d92c-4f66-915b-863f10140436",
			"err": err,
			"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'supportedTestDataDomainsWithHeadersMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK

		// Get number of Test Data Domain-rows
		numberOfDomains := len(supportedTestDataDomainsWithHeadersMessage.TestDataDomains)
		if numberOfDomains > 0 {
			// Loop over all  TestData-domains
			for currentDomainCounter := 0; currentDomainCounter < numberOfDomains; currentDomainCounter++ {

				// Get number of TestDataFilterHeader-rows
				numberOfTestDataFilterHeaders := len(supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders)
				if numberOfTestDataFilterHeaders > 0 {
					// Loop over all  TestDataFilterHeaders
					for currentDestDataFilterHeaderCounter := 0; currentDestDataFilterHeaderCounter < numberOfTestDataFilterHeaders; currentDestDataFilterHeaderCounter++ {

						// Get number of TestDataFilterHeaderValue-rows
						numberOfTestDataFilterHeaderValues := len(supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderValues)
						if numberOfTestDataFilterHeaderValues > 0 {
							// Loop over all TestDataFilterHeaderValues
							for currentTestDataFilterHeaderValueCounter := 0; currentTestDataFilterHeaderValueCounter < numberOfTestDataFilterHeaderValues; currentTestDataFilterHeaderValueCounter++ {

								/*

										create table testdatadomains.SupportedTestDataDomainsWithHeaders
										(
										    OriginalSenderId                uuid                     not null, -- The Id of the gateway/plugin that created the message 1
										    OriginalSenderName              varchar default null,              -- The name of the gateway/plugin that created the message 2
										    MessageId                       uuid                     not null, -- A unique id generated when sent from Plugin 3
										    OriginalMessageId               uuid                     not null, -- A unique id from the request sent from Fenix 4
										    TestDataDomainId                uuid                     not null, -- he unique id of the testdata domain 5
										    TestDataDomainName              varchar default null,              -- The name of the testdata domain 6
										    TestDataDomainDescription       varchar default null,              -- A description of the testdata domain 7
										    TestDataDomainMouseOver         varchar default null,              -- A mouse over description of the testdata domain 8
										    TestDataFilterHeaderId          uuid                     not null,-- A unique id for the header 9
										    TestDataFilterHeaderName        varchar                  not null, -- A name for the keader 10
										    TestDataFilterHeaderDescription varchar                  not null, -- A description of the header 11
										    TestDataFilterHeaderMouseOver   varchar                  not null, --  A mouse over description of the header 12
										    TestDataFilterHash              varchar                  not null, --  ????????????????????????????????????????????????????? 13
										    AllowMultipleChoices            boolean                  not null, -- Should multiple chocies be allowed for this header 14
										    AllowNoChoice                   boolean                  not null, -- Should no choice be allowed for this header 15
										    TestDataFilterHeaderValueId     uuid                     not null, -- A unique id for the filter value 16
										    TestDataFilterHeaderValueName   varchar                  not null, -- The name for the filter value 17
										    updatedDateTime                 timestamp with time zone not null  -- The Datetime when the row was created/updates 18
										);

									SupportedTestDataDomainsWithHeadersMessage_TestDataDomainWithHeaders
									SupportedTestDataDomainsWithHeadersMessage_TestDataDomainWithHeaders_TestDataFilterHeader
									SupportedTestDataDomainsWithHeadersMessage_TestDataDomainWithHeaders_TestDataFilterHeader_TestDataFilterHeaderValue
								*/

								// Values to insert into database
								sqlResult, err := sqlStatement.Exec(
									supportedTestDataDomainsWithHeadersMessage.OriginalSenderId,
									supportedTestDataDomainsWithHeadersMessage.OriginalSenderName,
									supportedTestDataDomainsWithHeadersMessage.MessageId,
									supportedTestDataDomainsWithHeadersMessage.OriginalMessageId,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataDomainId,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataDomainName,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataDomainDescription,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataDomainMouseOver,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderId,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderName,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderDescription,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderMouseOver,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHash,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].AllowMultipleChoices,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].AllowNoChoice,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderValues[currentTestDataFilterHeaderValueCounter].TestDataFilterHeaderValueId,
									supportedTestDataDomainsWithHeadersMessage.TestDataDomains[currentDomainCounter].TestDataFilterHeaders[currentDestDataFilterHeaderCounter].TestDataFilterHeaderValues[currentTestDataFilterHeaderValueCounter].TestDataFilterHeaderValueName,
									common_code.GeneraTimeStampUTC())

								if err != nil {
									// Error while executing
									logger.WithFields(logrus.Fields{
										"ID":        "b5852c95-a717-4fc6-8d06-0b80e8ca4868",
										"err":       err,
										"sqlResult": sqlResult,
										"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
									}).Error("Error when updating Main Database with data from 'supportedTestDataDomainsWithHeadersMessage'")

									messageSavedInDB = false
								} else {
									//SQL executed OK
									logger.WithFields(logrus.Fields{
										"ID":        "fe45a965-0e03-4d03-be2e-3b9eefffd5fb",
										"err":       err,
										"sqlResult": sqlResult,
										"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
									}).Debug("Fenix main Database was updated with data from 'supportedTestDataDomainsWithHeadersMessage'")
								}
							}
						}
					}
				}
			}
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'AvailbleTestInstructionAtPluginMessage' to Main Database for Fenix Inception
//
func saveAvailbleTestInstructionAtPluginMessageInDB(availbleTestInstructionAtPluginMessage *gRPC.AvailbleTestInstructionAtPluginMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO testInstructions.AvailbleTestInstructions "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, , MessageId, OrginalCreateDateTime, OriginalSystemDomainId, OriginalSystemDomainName,  "
	sqlToBeExecuted = sqlToBeExecuted + "PluginGuid, PluginName, SystemDomainId, SystemDomainName, TestInstructionTypeGuid, TestInstructionTypeName, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionGuid, TestInstructionName, TestInstructionDescription, TestInstructionMouseOver, TestInstructionVisible, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionEnable, TestInstructionColor, DropIDs, LandingZoneGuid, ChoosenLandingZone, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeGuid, TestI, ructionAttributeTypeGuid, TestInstructionAttributeName, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeDescription, TestInstructionAttributeMouseOver, TestInstructionAttributeVisible, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeEnable, TestInstructionAttributeMandatory, TestInstructionAttributeVisibleInBlockArea, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeTypeId, TestInstructionAttributeTypeName, TestInstructionAttributeInputTextBoxGuid, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeInputTextBoxName, TextBoxEditable, TextBoxInputMask, TextBoxAttributeTypeGuid, "
	sqlToBeExecuted = sqlToBeExecuted + "TextBoxAttributeTypeName, TextBoxAttributeValue,, TestInstructionAttributeComboBoxGuid, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeComboBoxName, ComboBoxEditable, ComboBoxInputMask, ComboBoxAttributeTypeGuid, "
	sqlToBeExecuted = sqlToBeExecuted + "ComboBoxAttributeTypeName, ComboBoxAttributeValueGuid, ComboBoxAttributeValue, updatedDateTime "

	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20, "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,38,$39,$40, "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($41,$42,$43,$44,$45,$46,$47,$48,$49);"

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                                     "ce25ce71-9437-457c-a5d4-d9f67b34db3e",
			"err":                                    err,
			"availbleTestInstructionAtPluginMessage": availbleTestInstructionAtPluginMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'availbleTestInstructionAtPluginMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK

		// Get number of TestInstructions-rows
		numberOfTestInstructions := len(availbleTestInstructionAtPluginMessage.TestInstructions)
		if numberOfTestInstructions > 0 {
			// Loop over all TestInstructions
			for currentTestInstruction := 0; currentTestInstruction < numberOfTestInstructions; currentTestInstruction++ {

				// Get number of Attribute-rows
				numberOfAttributes := len(availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes)
				if numberOfAttributes > 0 {
					// Loop over all  Attributes
					for currentAttributeCounter := 0; currentAttributeCounter < numberOfTestDataFilterHeaders; currentAttributeCounter++ {


						availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeDescription,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeMouseOver,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeVisible,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeEnable,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeMandatory,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeVisibleInBlockArea,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].type

						/*

													--A specific Attribute for the Test Instruction
							    TestInstructionAttributeGuid               uuid                     not null, -- 'The unique guid for the TestInstructionAttribute, set plugin'; --23
							    TestInstructionAttributeTypeGuid           uuid                     not null, -- 'The unique guid for the TestInstructionAttribute-type, set by plugin'; --24
							    TestInstructionAttributeName               varchar default null,              -- 'The name of the TestInstructionAttribute'; --25
							    TestInstructionAttributeDescription        varchar default null,              -- 'The description of the TestInstructionAttribute'; --26
							    TestInstructionAttributeMouseOver          varchar default null,              -- 'The mouse over text for the TestInstructionAttribute'; --27
							    TestInstructionAttributeVisible            boolean                  not null, -- 'Should the TestInstructionAttribute be visible in attributes list in GUI or not'; --29
							    TestInstructionAttributeEnable             boolean                  not null, -- 'Should the TestInstructionAttribute be enabled or not'; --29
							    TestInstructionAttributeMandatory          boolean                  not null, -- 'Should the TestInstructionAttribute be mandatory or not'; --30
							    TestInstructionAttributeVisibleInBlockArea boolean                  not null, -- 'Should the TestInstructionAttribute be visible in TestInstructionBlockAreain GUI or not'; --31
							    TestInstructionAttributeTypeId             uuid                     not null, -- 'The Id for what type the instruction attribute is'; --32
							    TestInstructionAttributeTypeName           varchar default null,              -- 'The Name for what type the instruction attribute is. Can one of the folowing types "TextBox", "ComboBox", "FileSelector", "FunctionSelector"'; --33
							-- Properties for TextBox attribute
							    TestInstructionAttributeInputTextBoxGuid   uuid    default null,              -- 'The unique guid for the TestInstructionAttributeInputTextBoxProperties, set by plugin'; --34
							    TestInstructionAttributeInputTextBoxName   varchar default null,              -- 'The name of the TestInstructionAttributeInputTextBoxProperties'; --35
							    TextBoxEditable                            boolean default null,              -- 'Should the the TextBox be editable or not'; --36
							    TextBoxInputMask                           varchar default null,              -- 'Inputmask for the TextBox'; --37
							    TextBoxAttributeTypeGuid                   uuid    default null,              -- 'The unique guid for the Type of the TextBox. Used for datamanupulation'; --38
							    TextBoxAttributeTypeName                   varchar default null,              -- 'The Name for the Type of the TextBox.'; --39
							    TextBoxAttributeValue
							create table testInstructions.AvailbleTestInstructions
							(
							    -- Base information
							    OriginalSenderId                           uuid                     not null, -- 'The Id of the gateway/plugin that created the message'; --1
							    OriginalSenderName                         varchar default null,              -- 'Then name of the plugin that created the message'; --2
							    MessageId                                  uuid                     not null, -- 'A unique id generated when sent from Plugin'; --3
							    OrginalCreateDateTime                      timestamp with time zone not null, -- 'The timestamp when the orignal message was created'; --4
							    OriginalSystemDomainId                     uuid                     not null, -- 'The Domain/system''s Id where the Sender operates'; --5
							    OriginalSystemDomainName                   varchar default null,              -- 'The Domain/system''s Name where the Sender operates';--6
							    -- A specific Test Instruction
							    PluginGuid                                 uuid                     not null, -- Used as unique id for the plugin'; --7
							    PluginName                                 varchar default null,              -- 'Used as unique name for the plugin'; --8
							    SystemDomainId                             uuid                     not null, -- 'The Domain/system''s Id where the Sender operates is'; --9
							    SystemDomainName                           varchar default null,              --  'The Domain/system''s Name where the Sender operates'; --10
							    TestInstructionTypeGuid                    uuid                     not null, --  'The unique guid for the Type of TestInstruction. Set by Client system'; --11
							    TestInstructionTypeName                    varchar default null,              --  'The name for the Type of TestInstruction. Set by Client system'; --12
							    TestInstructionGuid                        uuid                     not null, --  'The unique guid for the TestInstruction. Set when used in Editor'; --13
							    TestInstructionName                        varchar default null,              --  'The name of the TestInstruction'; --14
							    TestInstructionDescription                 varchar default null,              --  'The description of the TestInstruction'; --15
							    TestInstructionMouseOver                   varchar default null,              --  'The mouse over test for the TestInstruction'; --16
							    TestInstructionVisible                     boolean                  not null, --  'Should the TestInstruction be visible in GUI or not'; --17
							    TestInstructionEnable                      boolean                  not null, --  'Should the TestInstruction be enabled or not'; --18
							    TestInstructionColor                       varchar default null,              --  'The color used for presenting the TestInstructionBlock, e.g. #FAF437'; --19
							    DropIDs                                    varchar default null,              --  'A Drop-IDs deciding if receiver accepts obejct when drag n drop. Can have many Drop-IDs. ***For now saved in same column ***'; --20
							    LandingZoneGuid                            varchar default null,              --  'List with IDs of all available LandingZones for TestInstruction Can have many LandingZone-IDs. ***For now saved in same column ***'; --21
							    ChoosenLandingZone                         varchar default null,              --  'The choosen LandingZone for this TestInstructionBlock'; --22

							--A specific Attribute for the Test Instruction
							    TestInstructionAttributeGuid               uuid                     not null, -- 'The unique guid for the TestInstructionAttribute, set plugin'; --23
							    TestInstructionAttributeTypeGuid           uuid                     not null, -- 'The unique guid for the TestInstructionAttribute-type, set by plugin'; --24
							    TestInstructionAttributeName               varchar default null,              -- 'The name of the TestInstructionAttribute'; --25
							    TestInstructionAttributeDescription        varchar default null,              -- 'The description of the TestInstructionAttribute'; --26
							    TestInstructionAttributeMouseOver          varchar default null,              -- 'The mouse over text for the TestInstructionAttribute'; --27
							    TestInstructionAttributeVisible            boolean                  not null, -- 'Should the TestInstructionAttribute be visible in attributes list in GUI or not'; --29
							    TestInstructionAttributeEnable             boolean                  not null, -- 'Should the TestInstructionAttribute be enabled or not'; --29
							    TestInstructionAttributeMandatory          boolean                  not null, -- 'Should the TestInstructionAttribute be mandatory or not'; --30
							    TestInstructionAttributeVisibleInBlockArea boolean                  not null, -- 'Should the TestInstructionAttribute be visible in TestInstructionBlockAreain GUI or not'; --31
							    TestInstructionAttributeTypeId             uuid                     not null, -- 'The Id for what type the instruction attribute is'; --32
							    TestInstructionAttributeTypeName           varchar default null,              -- 'The Name for what type the instruction attribute is. Can one of the folowing types "TextBox", "ComboBox", "FileSelector", "FunctionSelector"'; --33
							-- Properties for TextBox attribute
							    TestInstructionAttributeInputTextBoxGuid   uuid    default null,              -- 'The unique guid for the TestInstructionAttributeInputTextBoxProperties, set by plugin'; --34
							    TestInstructionAttributeInputTextBoxName   varchar default null,              -- 'The name of the TestInstructionAttributeInputTextBoxProperties'; --35
							    TextBoxEditable                            boolean default null,              -- 'Should the the TextBox be editable or not'; --36
							    TextBoxInputMask                           varchar default null,              -- 'Inputmask for the TextBox'; --37
							    TextBoxAttributeTypeGuid                   uuid    default null,              -- 'The unique guid for the Type of the TextBox. Used for datamanupulation'; --38
							    TextBoxAttributeTypeName                   varchar default null,              -- 'The Name for the Type of the TextBox.'; --39
							    TextBoxAttributeValue                      varchar default null,              -- 'The value for the the TextBox, used for preset values';--40

							-- Properties for ComboBox attribute
							    TestInstructionAttributeComboBoxGuid       uuid    default null,              -- 'The unique guid for the TestInstructionAttributeComboBoxProperties, set by plugin'; --41
							    TestInstructionAttributeComboBoxName       varchar default null,              -- 'The name of the TestInstructionAttributeComboBoxProperties'; --42
							    ComboBoxEditable                           boolean default null,              -- 'Should the the ComboBox be editable or not'; --43
							    ComboBoxInputMask                          varchar default null,              -- 'Inputmask for the ComboBox'; --44
							    ComboBoxAttributeTypeGuid                  uuid    default null,              -- 'The unique guid for the Type of the ComboBox Used for datamanupulation'; --45
							    ComboBoxAttributeTypeName                  varchar default null,              -- 'The Name for the Type of the ComboBox'; --46
							    ComboBoxAttributeValueGuid                 uuid    default null,              -- 'The guid of the value for the the ComboBox, used for preset values'; --47
							    ComboBoxAttributeValue                     varchar default null,              -- 'The value for the the ComboBox, used for preset values';--48
							-- WHen row was create/updated
							    updatedDateTime                            timestamp with time zone not null  -- The Datetime when the row was created/updates 49
							);
						*/

						// Values to insert into database
						sqlResult, err := sqlStatement.Exec(
							availbleTestInstructionAtPluginMessage.OriginalSenderId,
							availbleTestInstructionAtPluginMessage.OriginalSenderName,
							availbleTestInstructionAtPluginMessage.MessageId,
							availbleTestInstructionAtPluginMessage.OrginalCreateDateTime,
							availbleTestInstructionAtPluginMessage.OriginalSystemDomainId,
							availbleTestInstructionAtPluginMessage.OriginalSenderName,

							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].PluginGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].PluginName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].SystemDomainId,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].SystemDomainName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionTypeName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionDescription,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionMouseOver,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionVisible,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionEnable,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionColor,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].DropIDs,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].LandingZoneGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].ChoosenLandingZone,



							common_code.GeneraTimeStampUTC())

						if err != nil {
							// Error while executing
							logger.WithFields(logrus.Fields{
								"ID":                                     "6d2edd97-9548-41fd-9cb1-4addffd563ab",
								"err":                                    err,
								"sqlResult":                              sqlResult,
								"availbleTestInstructionAtPluginMessage": availbleTestInstructionAtPluginMessage,
							}).Error("Error when updating Main Database with data from 'availbleTestInstructionAtPluginMessage'")

							messageSavedInDB = false
						} else {
							//SQL executed OK
							logger.WithFields(logrus.Fields{
								"ID":                                     "033ab15a-4e3a-4d3d-8e99-43a18b8d557c",
								"err":                                    err,
								"sqlResult":                              sqlResult,
								"availbleTestInstructionAtPluginMessage": availbleTestInstructionAtPluginMessage,
							}).Debug("Fenix main Database was updated with data from 'availbleTestInstructionAtPluginMessage'")
						}
					}
				}
			}
		}
	}
}
}

// Return if the messages was saved or not in database
return messageSavedInDB
}
/*
       //TestInstructionResults flows from Plugins towards Fenix Inception
       rpc SendTestInstructionResultTowardsFenix (TestInstructionExecutionResultMessage) returns (AckNackResponse) {
       }

       // Log-posts har sent from Plugins towards Fenix Inception
   x    rpc SendTestExecutionLogTowardsFenix (TestExecutionLogMessage) returns (AckNackResponse) {
       }

       // TestInstructions Timeouts flows from Plugins towards Fenix.
       // Used for telling Fenix when a TestInstruction should be have a timeout and execution for Testcase should be stopped
   x    rpc SendTestInstructionTimeOutTowardsFenix (TestInstructionTimeOutMessage) returns (AckNackResponse) {
       }

   x    rpc SupportedTestDataDomains (SupportedTestDataDomainsWithHeadersMessage) returns (AckNackResponse) {
       }

       // This one will probably be deleted
   //    rpc TestDataFromFilterValues (TestDataFromFilterValuesMessage) returns (AckNackResponse) {

       // Register an avalible TestInstruction
   x    rpc RegisterAvailbleTestInstructions(AvailbleTestInstructionAtPluginMessage) returns (AckNackResponse){
       }

       // Register the different Testdata Domains that is supported
       rpc RegistrateAvailableTestDataDomains(SupportedTestDataDomainsMessage) returns (AckNackResponse) {
       }

       // Register Test Containers that are supported. A Test Container consists of many TestInstructions grouped together into one unit
       rpc RegistrateAailableTestContainers(AvailbleTestContainersAtPluginMessage) returns (AckNackResponse){
       }

       // A Plugin register itself at a distributed gateway or a distributed gateway register itself at Fenix Inception gateway
       rpc RegisterClientAddress (RegisterClientAddressRequest) returns (RegisterClientAddressResponse) {

       // A Plugin can send a INFO- OR WARNING-message to Fenix by this mwethod
       rpc SendMessageToFenix (InformationMessage) returns (AckNackResponse){
*/
'


/*
   	fmt.Println("# Inserting values")

   	var lastInsertId int
   	err = mainDB.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
   	checkErr(err)
   	fmt.Println("last inserted id =", lastInsertId)

   	fmt.Println("# Updating")
   	stmt, err := mainDB.Prepare("update userinfo set username=$1 where uid=$2")
   	checkErr(err)

   	res, err := stmt.Exec("astaxieupdate", lastInsertId)
   	checkErr(err)

   	affect, err := res.RowsAffected()
   	checkErr(err)

   	fmt.Println(affect, "rows changed")

   	fmt.Println("# Querying")
   	rows, err := mainDB.Query("SELECT * FROM userinfo")
   	checkErr(err)

   	for rows.Next() {
   		var uid int
   		var username string
   		var department string
   		var created time.Time
   		err = rows.Scan(&uid, &username, &department, &created)
   		checkErr(err)
   		fmt.Println("uid | username | department | created ")
   		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
   	}

   	fmt.Println("# Deleting")
   	stmt, err = mainDB.Prepare("delete from userinfo where uid=$1")
   	checkErr(err)

   	res, err = stmt.Exec(lastInsertId)
   	checkErr(err)

   	affect, err = res.RowsAffected()
   	checkErr(err)

   	fmt.Println(affect, "rows changed")
   }

   func checkErr(err error) {
   	if err != nil {
   		panic(err)
   	}
   }
*/
