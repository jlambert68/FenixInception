package SendUtrViaMQ

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
		"ID": "c3ad43cc-1ba3-4d93-926f-374c42ddefc0",
	}).Debug("Initiate connection to Main Database for Fenix Inception")

	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	mainDB, err = sql.Open("postgres", dbinfo)

	// If some error occured then close the application
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "bbef4597-af3a-4e92-a5c2-a3ff8db6572e",
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
			"ID":                                    "e2514308-126d-4764-849c-0bdd4f7a770c",
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
				"ID":                                    "ef478a07-7c84-4796-afb2-ef361bc464ec",
				"err":                                   err,
				"sqlResult":                             sqlResult,
				"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
			}).Error("Error when updating Main Database with data from testInstructionExecutionResultMessage")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                                    "39d5f84d-4cae-4478-8927-446f53abd775",
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
			"ID":                      "c149f73b-6d41-48f5-8c65-8d72bc190f8b",
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
				"ID":        "1a35b329-f555-4913-94f0-56d542f29371",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when looking for Test Instructions-peers which has not been finished executing")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "63a07711-a460-4fff-98f4-d81d37a0e09a",
				"sqlResult": sqlResult,
			}).Debug("Found Test Instruction-peers in Fenix main Database")

			// Convert SQL-result into correct return message for function
			var testInstructionId string
			for sqlResult.Next() {
				err = sqlResult.Scan(&testInstructionId)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"ID":        "e9696b6a-740e-4bbf-947e-9fc8211ef522",
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
			"ID":                      "815e6471-2108-4fe0-91e5-03a5115b99aa",
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
				"ID":        "aae1d603-139b-4e1b-af6f-de32dbc0b2cb",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when looking for next TestInstructions-peers which should be executed")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "ec8bd41b-537a-479a-b06e-41bf4ce49aa2",
				"sqlResult": sqlResult,
			}).Debug("sqlResult for Test Instruction-peers to be executed in Fenix main Database")

			// Convert SQL-result into correct return message for function
			var testInstructionId string
			for sqlResult.Next() {
				err = sqlResult.Scan(&testInstructionId)
				if err != nil {
					// Error while looping through result
					logger.WithFields(logrus.Fields{
						"ID":        "d7a3ae0a-3bf6-4692-8d8d-450593d9b141",
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
			"ID":                "fd9268d6-5609-4872-8003-8e842236d25a",
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
				"ID":                "92f12d44-acc0-4960-9fa2-b0de5637ca88",
				"err":               err,
				"sqlResult":         sqlResult,
				"testInstructionId": testInstructionId,
			}).Error("Error when updating Main Database with data that TestInstruction was sent to Plugin")

			messageSavedInDB = false
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":                "8a73dc4a-3b9d-43b2-8282-897b96474c8f",
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
			"ID":                  "56378d2d-19d9-4602-b337-ba5a238ffcd7",
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
				"ID":        "57e09080-7bf0-4aee-bebc-11b47b2d6222",
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
						"ID":        "b30cea6c-5bae-4fe8-bc25-1a98dc3a9c96",
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
					"ID":                   "443e2efd-b22e-425b-adbe-593d2f2517eb",
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

/*
       //TestInstructionResults flows from Plugins towards Fenix Inception
   x    rpc SendTestInstructionResultTowardsFenix (TestInstructionExecutionResultMessage) returns (AckNackResponse) {
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
   x    rpc RegistrateAvailableTestDataDomains(SupportedTestDataDomainsMessage) returns (AckNackResponse) {
       }

       // Register Test Containers that are supported. A Test Container consists of many TestInstructions grouped together into one unit
   x    rpc RegistrateAailableTestContainers(AvailbleTestContainersAtPluginMessage) returns (AckNackResponse){
       }

       // A Plugin register itself at a distributed gateway or a distributed gateway register itself at Fenix Inception gateway
       rpc RegisterClientAddress (RegisterClientAddressRequest) returns (RegisterClientAddressResponse) {

       // A Plugin can send a INFO- OR WARNING-message to Fenix by this mwethod
  x     rpc SendMessageToFenix (InformationMessage) returns (AckNackResponse){
*/

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
