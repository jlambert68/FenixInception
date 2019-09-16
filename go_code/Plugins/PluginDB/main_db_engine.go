package PluginKeyValueDBStore

import (
	"database/sql"
	"errors"
	"fmt"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code"
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
// Update 'testInstructionExecutionResultMessage' in, TestInstructions Outgoing message, to Main Database for Fenix Inception
//
func updateTestInstructionWithExecutionResultMessageInDB(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (messageSavedInDB bool) {
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
// Update DB that TestInstruction has been sent to Plugin
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
// Save incoming 'TestInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func saveKeyValuetMessageInDB(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO testInstructions.TestExecutionResult "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, TestInstructionGuid, ResultId, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionResultCreatedDateTime, TestInstructionResultSentDateTime, "
	sqlToBeExecuted = sqlToBeExecuted + "AvailableTestExecutionResultTypeId, AvailableTestExecutionResultTypeName, "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSystemDomainId, OriginalSystemDomainName, MessageId, PeerId, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                                    "ac69559d-9b20-482e-b9eb-5c28865d0b29",
			"err":                                   err,
			"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'testInstructionExecutionResultMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		/*
			create table testInstructions.TestExecutionResult
			(
			    OriginalSenderId                     uuid                     not null, -- The Id of the gateway/plugin that created the message -- 1
			    OriginalSenderName                   varchar default null,              -- The name of the gateway/plugin that created the message -- 2
			    TestInstructionGuid                  uuid                     not null, --  TestInstructionGuid is a unique id for identifying the TestInstruction at Plugin; --3
			    ResultId                             uuid                     not null, -- A unique id created by plugin'; -- 4
			    TestInstructionResultCreatedDateTime timestamp with time zone not null, -- The DateTime when the message was created'; -- 5
			    TestInstructionResultSentDateTime    timestamp with time zone not null, -- The DateTime when the message was sent''; -- 6
			    AvailableTestExecutionResultTypeId   uuid                     not null, -- The Domain/system''''s Id where the Sender operates''; -- 7
			    AvailableTestExecutionResultTypeName varchar                  not null, -- The Domain/system''''s Name where the Sender operates''; -- 8
			    OriginalSystemDomainId               uuid                     not null, -- A unique id for the message created by plugin''; -- 9
			    OriginalSystemDomainName             varchar default null,              -- 'The Domain/system''s name where the Sender operates''; -- 10
			    MessageId                            uuid                     not null, -- A unique id for the message created by plugin''; -- 11
			    PeerId                               uuid                     not null, --A unique id for all peer TestInstructions that can be processed in parallell with current TestInstruction''; -- 12
			    updatedDateTime                      timestamp with time zone not null  -- The Datetime when the row was created/updates''; -- 13
			);

		*/
		// Values to insert into database
		sqlResult, err := sqlStatement.Exec(
			testInstructionExecutionResultMessage.OriginalSenderId,
			testInstructionExecutionResultMessage.OriginalSenderName,
			testInstructionExecutionResultMessage.TestInstructionGuid,
			testInstructionExecutionResultMessage.ResultId,
			testInstructionExecutionResultMessage.TestInstructionResultCreatedDateTime,
			testInstructionExecutionResultMessage.TestInstructionResultSentDateTime,
			testInstructionExecutionResultMessage.TestExecutionResult,
			gRPC.AvailableTestExecutionResultsEnum_name[int32(testInstructionExecutionResultMessage.TestExecutionResult)],
			testInstructionExecutionResultMessage.OriginalSystemDomainId,
			testInstructionExecutionResultMessage.OriginalSystemDomainName,
			testInstructionExecutionResultMessage.MessageId,
			testInstructionExecutionResultMessage.PeerId,
			common_code.GeneraTimeStampUTC())

		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                                    "5007c17f-eddb-4f45-9fda-24368f670850",
				"err":                                   err,
				"sqlResult":                             sqlResult,
				"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
			}).Error("Error when updating Main Database with data from 'testInstructionExecutionResultMessage'")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                                    "e0ea8ad4-0b78-4de5-b1d2-60a63a3b6df2",
				"err":                                   err,
				"sqlResult":                             sqlResult,
				"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
			}).Debug("Fenix main Database was updated with data from 'testInstructionExecutionResultMessage'")
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}
