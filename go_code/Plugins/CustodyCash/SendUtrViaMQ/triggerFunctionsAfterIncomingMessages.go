package SendUtrViaMQ

import "github.com/sirupsen/logrus"

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case checking if all peers are finished and if so then trigger next peers of TestInstructions to be executed
func newIncomingTestInstructionExecutionResultMessage(peerId string) (err error) {

	var testInstructionsThatAreStillExecuting []string

	// Get all Testinstructions that are a peer to this TestInstruction, and can be run in parallell, and are still executing
	testInstructionsThatAreStillExecuting, err = listPeerTestInstructionPeersWhichIsExecuting(peerId)

	// Trigger next TestInstructions that is waiting to be executed if all current peers are finished
	if err == nil && len(testInstructionsThatAreStillExecuting) == 0 {
		testInstructionPeersThatShouldBeExecutedNext, err := listNextPeersToBeExecuted(peerId)
		if err == nil && len(testInstructionsThatAreStillExecuting) > 0 {
			err = triggerSendNextPeersForExecution(testInstructionPeersThatShouldBeExecutedNext)
		}
	}

	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingSupportedTestDataDomainsWithHeadersMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingTestExecutionLogMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingTestInstructionTimeOutMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingInformationMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingSupportedTestDataDomainsMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingAvailbleTestContainersAtPluginMessage() (err error) {

	err = nil
	return err
}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case ?????????????????????????
func newIncomingAvailbleTestInstructionAtPluginMessage() (err error) {

	err = nil
	return err
}

// ********************************************************************************************
// Trigger the list of TestInstructions to be sent towards Plugins
//
func triggerSendNextPeersForExecution(testInstructionPeersThatShouldBeExecutedNext []string) (err error) {

	var numberOfTestInstructions int
	var testInstructionId string
	var messageSavedInDB bool

	err = nil

	numberOfTestInstructions = len(testInstructionPeersThatShouldBeExecutedNext)
	if numberOfTestInstructions > 0 {

		// Loop through all TestInstructionId's
		for arrayIndex := 0; arrayIndex < numberOfTestInstructions; arrayIndex++ {
			testInstructionId = testInstructionPeersThatShouldBeExecutedNext[arrayIndex]
			// Send the TestInstruction to next Gateway
			err = sendTestInstructionTowardsPlugin(testInstructionId)
			if err != nil {
				// Not sent towards plugin
				logger.WithFields(logrus.Fields{
					"ID":                "6d35b81c-bc28-4314-babb-96300b682f0c",
					"err":               err,
					"testInstructionId": testInstructionId,
				}).Error("TestInstruction was NOT sent towards Plugin")

			} else {
				// Sent towards plugin
				logger.WithFields(logrus.Fields{
					"ID":                "6fd16e88-cfba-41ec-8b46-eb7ab94ea6be",
					"testInstructionId": testInstructionId,
				}).Debug("TestInstruction was sent towards Plugin")

				// Change TestInstruction in DB that it has been sent
				messageSavedInDB = saveInDbThatTestInstructionHasBeenSentToPlugin(testInstructionId)
				if messageSavedInDB == false {
					// Not svaed changes to Database
					logger.WithFields(logrus.Fields{
						"ID":                "4f95e5c3-dde7-48eb-8088-5e7f9ee4f226",
						"err":               err,
						"testInstructionId": testInstructionId,
					}).Error("TestInstruction was sent towards plugin, was NOT saved in Fenix Database")

				} else {
					// Saved changes to database
					logger.WithFields(logrus.Fields{
						"ID":                "58bd7e0f-6df7-41c3-a94d-e137a99c97bd",
						"testInstructionId": testInstructionId,
					}).Debug("TestInstruction was sent towards plugin, was saved in Fenix Database")
				}
			}
		}
	}
	return err
}
