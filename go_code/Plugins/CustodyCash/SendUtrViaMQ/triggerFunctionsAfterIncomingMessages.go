package SendUtrViaMQ

import (
	"github.com/sirupsen/logrus"
)

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case TestInstructionMessage
func newIncomingTestInstructionMessage() {

}

// ************************************************************************************
// Triggers functionality after the incoming message was saved in database
// In this case SupportedTestDataDomainsRequest
func newIncomingSupportedTestDataDomainsRequest() error {
	// Message saved OK
	logger.WithFields(logrus.Fields{
		"ID": "2159c10a-46cc-47ef-a4da-c741580370e2",
	}).Debug("'supportedTestDataDomainsRequest' was received in function 'newIncomingSupportedTestDataDomainsRequest'")

	return nil

}
