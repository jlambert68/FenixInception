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
fixa vidare här
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
