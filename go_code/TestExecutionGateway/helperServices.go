package TestExecutionGateway

import (
"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)


// Generate a new unique uuid
func generateUUID() (string) {
	var new_uuid_string string = ""

	new_uuid, err := uuid.NewUUID()
	new_uuid_string = new_uuid.String()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":            "d67c6284-2629-43c3-9a1c-6cff2dc403fa",
			"error message": err,
		}).Error("Couldn't generate a new UUID")
		//TODO Send message to Fenix
	} else
}