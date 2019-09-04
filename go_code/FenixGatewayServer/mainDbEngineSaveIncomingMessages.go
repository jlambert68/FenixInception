package FenixGatewayServer

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"jlambert/FenixInception2/go_code/common_code"
)

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
			testExecutionLogMessage.TestInstructionLogCreatedDateTime,
			testExecutionLogMessage.TestInstructionLogSentDateTime,
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
	sqlToBeExecuted = sqlToBeExecuted + "MajorVersion, MinorVersion, "
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
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51);"

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
					for currentAttributeCounter := 0; currentAttributeCounter < numberOfAttributes; currentAttributeCounter++ {

						/*MajorVersion

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
						    MajorVersion                               int                      not null, --  The version numbering system consists of major and minor. A change in major version forces the user to upgrade the TestInstruction --23
						    MinorVersion                               int                      not null, -- The version numbering system consists of major and minor. A change in minor version makes it optional for  the user to upgrade the TestInstruction --24

						--A specific Attribute for the Test Instruction
						    TestInstructionAttributeGuid               uuid                     not null, -- 'The unique guid for the TestInstructionAttribute, set plugin'; --25
						    TestInstructionAttributeTypeGuid           uuid                     not null, -- 'The unique guid for the TestInstructionAttribute-type, set by plugin'; --26
						    TestInstructionAttributeName               varchar default null,              -- 'The name of the TestInstructionAttribute'; --27
						    TestInstructionAttributeDescription        varchar default null,              -- 'The description of the TestInstructionAttribute'; --28
						    TestInstructionAttributeMouseOver          varchar default null,              -- 'The mouse over text for the TestInstructionAttribute'; --29
						    TestInstructionAttributeVisible            boolean                  not null, -- 'Should the TestInstructionAttribute be visible in attributes list in GUI or not'; --30
						    TestInstructionAttributeEnable             boolean                  not null, -- 'Should the TestInstructionAttribute be enabled or not'; --31
						    TestInstructionAttributeMandatory          boolean                  not null, -- 'Should the TestInstructionAttribute be mandatory or not'; --32
						    TestInstructionAttributeVisibleInBlockArea boolean                  not null, -- 'Should the TestInstructionAttribute be visible in TestInstructionBlockAreain GUI or not'; --33
						    TestInstructionAttributeTypeId             uuid                     not null, -- 'The Id for what type the instruction attribute is'; --34
						    TestInstructionAttributeTypeName           varchar default null,              -- 'The Name for what type the instruction attribute is. Can one of the folowing types "TextBox", "ComboBox", "FileSelector", "FunctionSelector"'; --35

						-- Properties for TextBox attribute
						    TestInstructionAttributeInputTextBoxGuid   uuid    default null,              -- 'The unique guid for the TestInstructionAttributeInputTextBoxProperties, set by plugin'; --36
						    TestInstructionAttributeInputTextBoxName   varchar default null,              -- 'The name of the TestInstructionAttributeInputTextBoxProperties'; --37
						    TextBoxEditable                            boolean default null,              -- 'Should the the TextBox be editable or not'; --38
						    TextBoxInputMask                           varchar default null,              -- 'Inputmask for the TextBox'; --39
						    TextBoxAttributeTypeGuid                   uuid    default null,              -- 'The unique guid for the Type of the TextBox. Used for datamanupulation'; --40
						    TextBoxAttributeTypeName                   varchar default null,              -- 'The Name for the Type of the TextBox.'; --41
						    TextBoxAttributeValue                      varchar default null,              -- 'The value for the the TextBox, used for preset values';--42

						-- Properties for ComboBox attribute
						    TestInstructionAttributeComboBoxGuid       uuid    default null,              -- 'The unique guid for the TestInstructionAttributeComboBoxProperties, set by plugin'; --43
						    TestInstructionAttributeComboBoxName       varchar default null,              -- 'The name of the TestInstructionAttributeComboBoxProperties'; --44
						    ComboBoxEditable                           boolean default null,              -- 'Should the the ComboBox be editable or not'; --45
						    ComboBoxInputMask                          varchar default null,              -- 'Inputmask for the ComboBox'; --46
						    ComboBoxAttributeTypeGuid                  uuid    default null,              -- 'The unique guid for the Type of the ComboBox Used for datamanupulation'; --47
						    ComboBoxAttributeTypeName                  varchar default null,              -- 'The Name for the Type of the ComboBox'; --48
						    ComboBoxAttributeValueGuid                 uuid    default null,              -- 'The guid of the value for the the ComboBox, used for preset values'; --49
						    ComboBoxAttributeValue                     varchar default null,              -- 'The value for the the ComboBox, used for preset values';--50
						-- WHen row was create/updated
						    updatedDateTime                            timestamp with time zone not null  -- The Datetime when the row was created/updates 51
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
							availbleTestInstructionAtPluginMessage.
								availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeDescription,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeMouseOver,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeVisible,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeEnable,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeMandatory,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeVisibleInBlockArea,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeType,
							gRPC.TestInstructionAttributeTypeEnum_name[int32(availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].TestInstructionAttributeType)],

							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxAttributeTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TestInstructionAttributeInputTextBoxName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxEditable,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxInputMask,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxAttributeTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxAttributeTypeName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputTextBoxProperty.TextBoxAttributeValue,

							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.TestInstructionAttributeComboBoxGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.TestInstructionAttributeComboBoxName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxEditable,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxInputMask,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxAttributeTypeGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxAttributeTypeName,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxAttributeValueGuid,
							availbleTestInstructionAtPluginMessage.TestInstructions[currentTestInstruction].TestInstructionAttributes[currentAttributeCounter].InputComboBoxProperty.ComboBoxAttributeValue,

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

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// Structure used for passing base Container-message from main message to every container for saveing it togethor with container-information
type containerMessageBaseInformationStruct struct {
	OriginalSenderId         string
	OriginalSenderName       string
	MessageId                string
	OrginalCreateDateTime    string
	OriginalSystemDomainId   string
	OriginalSystemDomainName string
}

// **********************************************************************************************************
// Save incoming 'AvailbleTestContainersAtPluginMessage' to Main Database for Fenix Inception
// If a Container contains a child-container then this function is called in recursion
//
func saveavailbleTestContainersAtPluginMessageInDB(availbleTestContainersAtPluginMessage *gRPC.AvailbleTestContainersAtPluginMessage) (messageSavedInDB bool) {

	var messageBaseInformation containerMessageBaseInformationStruct
	var onlyTopContainers bool = true

	messageSavedInDB = true

	messageBaseInformation = containerMessageBaseInformationStruct{
		OriginalSenderId:         availbleTestContainersAtPluginMessage.OriginalSenderId,
		OriginalSenderName:       availbleTestContainersAtPluginMessage.OriginalSenderName,
		MessageId:                availbleTestContainersAtPluginMessage.MessageId,
		OrginalCreateDateTime:    availbleTestContainersAtPluginMessage.OrginalCreateDateTime,
		OriginalSystemDomainId:   availbleTestContainersAtPluginMessage.OriginalSystemDomainId,
		OriginalSystemDomainName: availbleTestContainersAtPluginMessage.OriginalSystemDomainName,
	}

	// Get number of Child-rows
	numberOfChildren := len(availbleTestContainersAtPluginMessage.TestContainerMessages)
	if numberOfChildren > 0 {
		// Loop over all Children to check if any of them is a TestInstruction
		// Only Containers can be sent in on first level
		for currentChildCounter := 0; currentChildCounter < numberOfChildren; currentChildCounter++ {
			// If this child is not defined as a TopCOntainer then something is wrong with the Container data
			if availbleTestContainersAtPluginMessage.TestContainerMessages[currentChildCounter].TestContainerIsTopContainer == false {
				onlyTopContainers = false
			}
		}

		if onlyTopContainers == false {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                                    "852da9b4-007b-4968-bdcc-6f1598582408",
				"onlyTopContainers":                     onlyTopContainers,
				"availbleTestContainersAtPluginMessage": availbleTestContainersAtPluginMessage,
			}).Error("Error when updating Main Database with data from 'availbleTestContainersAtPluginMessage'")

			messageSavedInDB = false
		} else {
			//Only Top Containers so extracting the data and save it in database
			logger.WithFields(logrus.Fields{
				"ID":                                    "bbbef665-d1f1-4267-b034-960d34c6d101",
				"availbleTestContainersAtPluginMessage": availbleTestContainersAtPluginMessage,
			}).Debug("Only Top Containers so try to extract the data and save it in database")

			// Loop over all Children to extract the container data to be saved in database
			for currentChildCounter := 0; currentChildCounter < numberOfChildren; currentChildCounter++ {
				saveTestContainerMessageInDB(
					messageBaseInformation,
					availbleTestContainersAtPluginMessage.TestContainerMessages[currentChildCounter])
			}
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'TestContainerMessage' to Main Database for Fenix Inception
// If a Container contains a child-container then this function is called in a recursive call
//
func saveTestContainerMessageInDB(messageBaseInformation containerMessageBaseInformationStruct, testContainerMessage *gRPC.TestContainerMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO testInstructions.TestContainerMessage "
	sqlToBeExecuted = sqlToBeExecuted + "PluginGuid, PluginName, SystemDomainId, SystemDomainName, "
	sqlToBeExecuted = sqlToBeExecuted + "TestContainerGuid, TestContainerName, MessageId, TestContainerIsTopContainer, "
	sqlToBeExecuted = sqlToBeExecuted + "ChildProcessingTypeId, ChildProcessingTypeName, ChildTypeId, ChildTypeName, "
	sqlToBeExecuted = sqlToBeExecuted + " ParentContainerReference, ChildContainerReference, TestInstructionsReference, "
	sqlToBeExecuted = sqlToBeExecuted + "OrginalCreateDateTime, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                   "9056ea34-bdfc-47cd-8ef7-974afc5a9bbe",
			"err":                  err,
			"testContainerMessage": testContainerMessage,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'testContainerMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK

		// Get number of Child-rows
		numberOfChildren := len(testContainerMessage.TestContainerChildMessages)
		if numberOfChildren > 0 {
			// Loop over all Children
			for currentChildCounter := 0; currentChildCounter < numberOfChildren; currentChildCounter++ {

				/*
					create table testInstructions.TestContainerMessage
					(
					    OriginalSenderId            uuid                     not null, -- The Id of the gateway/plugin that created the message 1
					    OriginalSenderName          varchar default null,              -- The name of the gateway/plugin that created the message 2
					    SystemDomainId              uuid                     not null, --  The Domain/system's Id where the Sender operates 3
					    SystemDomainName            varchar default null,              -- The Domain/system's Name where the Sender operates 4
					    TestContainerGuid           uuid                     not null, -- A unique id for the TestContainer 5
					    TestContainerName           varchar default null,              -- A name for the TestContainer 6
					    MessageId                   uuid                     not null, -- A unique id for the message 7
					    TestContainerIsTopContainer boolean                  not null, -- Telling if the container is a top container, i.e. has no parent container 8

					    ChildProcessingTypeId       int                      not null, -- The id for child processing type (0 or 1) 9
					    ChildProcessingTypeName     varchar                  not null, -- The name for child processing type ('IsSerialProcessedContainer' or 'IsParallellProcessedContainer') 10
					    ChildTypeId                 int                      not null, -- The id for the type of child  (0 or 1) 11
					    ChildTypeName               varchar                  not null, -- The name for the type of child ('ChildIsTestInstructionContainerMessage' or 'ChildIsTestInstructionMessage') 12
					    ParentContainerReference    uuid    default null,              -- A reference to parent container if it exists 13
					    ChildContainerReference     uuid    default null,              -- A reference to child container if it exists 14
					    TestInstructionsReference   uuid    default null,              -- A reference TestInstruction if it exists 15

					    OrginalCreateDateTime       timestamp with time zone not null, -- 'The timestamp when the orignal container was created'; --16
					    updatedDateTime             timestamp with time zone not null  -- The Datetime when the row was created/updates 17
					);
				*/
				// Values to insert into database
				sqlResult, err := sqlStatement.Exec(
					messageBaseInformation.OriginalSenderId,
					messageBaseInformation.OriginalSenderName,
					messageBaseInformation.OriginalSystemDomainId,
					messageBaseInformation.OriginalSystemDomainName,
					testContainerMessage.TestContainerGuid,
					testContainerMessage.TestContainerName,
					messageBaseInformation.MessageId,
					testContainerMessage.TestContainerIsTopContainer,

					testContainerMessage.TestContainerType,
					gRPC.TestContainerMessage_TestContainerTypeEnum_name[int32(testContainerMessage.TestContainerType)],
					testContainerMessage.TestContainerChildMessages[currentChildCounter].ChildType,
					gRPC.TestContainerMessage_TestContainerTypeEnum_name[int32(testContainerMessage.TestContainerChildMessages[currentChildCounter].ChildType)],
					testContainerMessage.TestContainerChildMessages[currentChildCounter].ParentContainerReference,
					"",
					testContainerMessage.TestContainerChildMessages[currentChildCounter].AvailbleTestInstructionAtPluginMessageChildReferenceGuid,
					common_code.GeneraTimeStampUTC())

				if err != nil {
					// Error while executing
					logger.WithFields(logrus.Fields{
						"ID":                   "af9d5ff3-6552-4e3a-92b2-c11113087a85",
						"err":                  err,
						"sqlResult":            sqlResult,
						"testContainerMessage": testContainerMessage,
					}).Error("Error when updating Main Database with data from 'testContainerMessage'")

					messageSavedInDB = false
				} else {
					//SQL executed OK
					logger.WithFields(logrus.Fields{
						"ID":                   "bbbef665-d1f1-4267-b034-960d34c6d101",
						"err":                  err,
						"sqlResult":            sqlResult,
						"testContainerMessage": testContainerMessage,
					}).Debug("Fenix main Database was updated with data from 'testContainerMessage'")

					// If this child is a container (ChildIsTestInstructionContainerMessage == 0) then call this function in a recursive call
					if testContainerMessage.TestContainerChildMessages[currentChildCounter].ChildType == 0 {
						messageSavedInDB = saveTestContainerMessageInDB(
							messageBaseInformation,
							testContainerMessage.TestContainerChildMessages[currentChildCounter].TestInstructionMessageChild)
					}
				}
			}
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Save incoming 'TestInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func saveTestInstructionExecutionResultMessageInDB(testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (messageSavedInDB bool) {

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
