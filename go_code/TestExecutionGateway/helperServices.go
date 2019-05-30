package TestExecutionGateway

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"time"
)

// *********************************************************************************
// Generate a new unique uuid
//
func generateUUID() string {
	var new_uuid_string string = ""

	new_uuid, err := uuid.NewUUID()
	new_uuid_string = new_uuid.String()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":            "d67c6284-2629-43c3-9a1c-6cff2dc403fa",
			"error message": err,
		}).Fatal("Couldn't generate a new UUID, stopping execution of Gateway")

	}
	// Return newly created UUID
	return new_uuid_string
}

// *********************************************************************************
// Genrerate UTC DateTime timestamp
//
func generaTimeStampUTC() string {
	now := time.Now()
	return now.String()
}

// *********************************************************************************
// Get Clients IP-address and Port from final PluginId from Memory object
//
func getClientAddressAndPort(pluginId string) (addressAndPort string) {
	// Get Client-info object from Clients-memory-object
	clientsAddressAndPortInfo, keyExists := clientsAddressAndPort[pluginId]

	if keyExists == true {
		addressAndPort = clientsAddressAndPortInfo.clientAddress + clientsAddressAndPortInfo.clientPort

	} else {
		// Key not found, send message to Fenix about Error
		logger.WithFields(logrus.Fields{
			"ID":       "80a08767-b706-4e6f-bcee-b1b19202d165",
			"pluginId": pluginId,
		}).Fatal("Couldn't not find address-data about Plugin, stopping execution of Gateway")

	}
	return addressAndPort

}

// *********************************************************************************
// Get Parents IP-address and Port from Memory object
//
func getParentAddressAndPort() (addressAndPort string) {
	// Get Parent-info object from Clients-memory-object
	ParentAddressAndPortInfo := gatewayConfig.parentgRPCAddress

	addressAndPort = ParentAddressAndPortInfo.parentGatewayServer_address + ParentAddressAndPortInfo.parentGatewayServer_port

	return addressAndPort

}

// *********************************************************************************
// Save Message in Local DB
//
func SaveMessageToLocalDB(
	key string,
	valueToStoreInDB []byte,
	bucket string,
	id string,
) {

	// Create return channel for save-status from DB
	returnChannel := make(chan dbResultMessage_struct)

	dbMessage := dbMessage_struct{
		DB_WRITE,
		bucket,
		key,
		valueToStoreInDB,
		returnChannel}

	// Send message to Database
	dbMessageQueue <- dbMessage

	// Wait for result on result channel and then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.err != nil {

		LogErrorAndSendInfoToFenixFor_SendTestExecutionLogTowardsFenix(
			id,
			gRPC.InformationMessage_FATAL,
			"",
			"",
			returnDBMessage.err.Error(),
			"Got an error when Saveing to local DB",
		)
	}
}
