package common_code

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gRPC "jlambert/FenixInception2/go_code/common_code/Gateway_gRPC_api"
	"time"
)

// *********************************************************************************
// Generate a new unique uuid
//
func GenerateUUID(logger *logrus.Logger) string {
	var newUuidString = ""

	newUuid, err := uuid.NewUUID()
	newUuidString = newUuid.String()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":            "d67c6284-2629-43c3-9a1c-6cff2dc403fa",
			"error message": err,
		}).Fatal("Couldn't generate a new UUID, stopping execution of Gateway")

	}
	// Return newly created UUID
	return newUuidString
}

// *********************************************************************************
// Genrerate UTC DateTime timestamp
//
func GeneraTimeStampUTC() string {
	now := time.Now()
	return now.String()
}

// *********************************************************************************
// Genrerate DateTime timestamp - "2019-06-27 21:33:17"
//
func generaTimeStampDateDateTime() string {
	now := time.Now()
	return now.String()[0:17]
}

// *********************************************************************************
// Get highest gRPC version
//

func GetHighestGRPCVersion() (currentVersion string) {
	maxVersionCount := int32(len(gRPC.CurrentVersionEnum_name) - 1)

	maxVersion := gRPC.CurrentVersionEnum_name[maxVersionCount]
	return maxVersion
}
