package SendUtrViaMQ

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"jlambert/FenixInception2/go_code/common_code"
)

// **********************************************************************************************************
// Save incoming 'SupportedTestDataDomainsRequest' to Main Database for Fenix Inception
//
func saveSupportedTestDataDomainsMessageInDB(supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (messageSavedInDB bool) {

	messageSavedInDB = true
	fixa
	vidare
	här
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
// Save incoming 'TestInstruction_RT' to Main Database for Fenix Inception
//
func saveTestInstructionRTMessageInDB(testInstructionRT *gRPC.TestInstruction_RT) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL

	var sqlToBeExecuted = "INSERT INTO plugins.TestInstructionRT "
	sqlToBeExecuted = sqlToBeExecuted + "OriginalSenderId, OriginalSenderName, TestInstructionTypeGuid, TestInstructionTypeName, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionGuid, PluginId, PluginName, TestInstructionCreatedDateTime, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionSentDateTime, ChoosenLandingZone, OrginalCreateDateTime, PeerId, "
	sqlToBeExecuted = sqlToBeExecuted + "MessageId, TestInstructionAttributeGuid, TestInstructionAttributeName, TestInstructionAttributeIndex, "
	sqlToBeExecuted = sqlToBeExecuted + "TestInstructionAttributeValue, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                "80c82ad5-0257-4a5c-a429-b4d0d6524622",
			"err":               err,
			"testInstructionRT": testInstructionRT,
		}).Error("Error when Praparing SQL for updating Main Database with data from 'testInstructionRT'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK

		// Get number of Attribute-rows
		numberOfAttributes := len(testInstructionRT.TestInstructionAttributes)
		if numberOfAttributes > 0 {
			// Loop over all Attributes
			for currentAttribute := 0; currentAttribute < numberOfAttributes; currentAttribute++ {

				/*
					create table plugins.TestInstructionRT
					(
					OriginalSenderId               uuid                     not null, -- The Id of the gateway/Fenix that created the message -- 1
					OriginalSenderName             varchar default null,              -- The name of the gateway/Fenix that created the message -- 2
					TestInstructionTypeGuid        uuid                     not null, -- TestInstructionTypeGuid tells Plugin what to do, set by Plugin'; --3
					TestInstructionTypeName        varchar default null,              -- The name of the type of TestInstruction type'; -- 4
					TestInstructionGuid            uuid                     not null, -- TestInstructionGuid is a unique id created when TestInstruction is created in TestCase'; -- 5
					PluginId                       uuid                     not null, -- The unique id of the plugin; -- 6
					PluginName                     varchar default null,              -- The name of the plugin; -- 7
					TestInstructionCreatedDateTime timestamp with time zone not null, -- When the TestInstruction, RunTime-version, was created; -- 8
					TestInstructionSentDateTime    timestamp with time zone not null, -- When the TestInstruction, RunTime-version, was sent from Fenix. Can be resent if errors accours after sent; -- 9
					ChoosenLandingZone             varchar default null,              -- Chosen LandingZone for TestInstruction; -- 10
					OrginalCreateDateTime          timestamp with time zone not null, -- The timestamp when the orignal message was created; -- 11
					PeerId                         uuid                     not null, -- A unique id for all peer TestInstructions that can be processed in parallell with current TestInstruction; -- 12
					MessageId                      uuid                     not null, -- A unique id for the message gererated by Fenix; -- 13
					TestInstructionAttributeGuid   uuid                     not null, -- A Unique id for the attribute, set by Plugin; -- 14
					TestInstructionAttributeName   varchar default null,              -- The name of the Attribute, set by Plugin; -- 15
					TestInstructionAttributeIndex  int                      not null, -- The index for the attrubute; -- 16
					TestInstructionAttributeValue  varchar default null,              -- The value of the attribute; -- 17
					updatedDateTime                timestamp with time zone not null  -- The Datetime when the row was created/updates''; -- 18
					);
				*/

				// Values to insert into database
				sqlResult, err := sqlStatement.Exec(
					testInstructionRT.OriginalSenderId,
					testInstructionRT.OriginalSenderName, testInstructionRT,
					testInstructionRT.TestInstructionTypeGuid,
					testInstructionRT.TestInstructionTypeName,
					testInstructionRT.TestInstructionGuid,
					testInstructionRT.PluginId,
					testInstructionRT.PluginName,
					testInstructionRT.TestInstructionCreatedDateTime,
					testInstructionRT.TestInstructionSentDateTime,
					testInstructionRT.ChoosenLandingZone,
					testInstructionRT.OrginalCreateDateTime,
					testInstructionRT.PeerId,
					testInstructionRT.MessageId,
					testInstructionRT.TestInstructionAttributes[currentAttribute].TestInstructionAttributeGuid,
					testInstructionRT.TestInstructionAttributes[currentAttribute].TestInstructionAttributeName,
					currentAttribute,
					testInstructionRT.TestInstructionAttributes[currentAttribute].TestInstructionAttributeValue,
					common_code.GeneraTimeStampUTC())

				if err != nil {
					// Error while executing
					logger.WithFields(logrus.Fields{
						"ID":                "4affa465-83be-420f-ab0a-dc63a7ad3ad9",
						"err":               err,
						"sqlResult":         sqlResult,
						"testInstructionRT": testInstructionRT,
					}).Error("Error when updating Main Database with data from 'testInstructionRT'")

					messageSavedInDB = false
				} else {
					//SQL executed OK
					logger.WithFields(logrus.Fields{
						"ID":                "74a3ef96-19d5-4881-a224-0599f51a8f93",
						"err":               err,
						"sqlResult":         sqlResult,
						"testInstructionRT": testInstructionRT,
					}).Debug("Fenix main Database was updated with data from 'testInstructionRT'")
				}
			}
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}
