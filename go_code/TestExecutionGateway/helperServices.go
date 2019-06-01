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
// Log message to local log and then Send message to Fenix
func LogErrorAndSendInfoToFenix(
	id string,
	messageType gRPC.InformationMessage_InformationType,
	infoHeader string,
	info string,
	errorMessage string,
	message string,
) {

	switch messageType {
	case gRPC.InformationMessage_DEBUG:
		// Only logg information and do not send to Fenix
		logger.WithFields(logrus.Fields{
			"ID":       id,
			infoHeader: info,
		}).Debug(message)

	case gRPC.InformationMessage_INFO:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id,
			infoHeader: info,
			"error":    errorMessage,
		}).Info(message)

		// Send information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_INFO,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_WARNING:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id,
			infoHeader: info,
			"error":    errorMessage,
		}).Warning(message)

		// Send Warning information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_WARNING,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_ERROR:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":       id,
			infoHeader: info,
			"error":    errorMessage,
		}).Error(message)

		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.callingSystemId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.callingSystemName,
			SenderId:              gatewayConfig.gatewayIdentification.callingSystemId,
			SenderName:            gatewayConfig.gatewayIdentification.callingSystemName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_ERROR,
			Message:               message,
			OrginalCreateDateTime: generaTimeStampUTC(),
		}

	case gRPC.InformationMessage_FATAL:
		// Only log and then Terminate Gateway due to problems
		logger.WithFields(logrus.Fields{
			"ID":       id,
			infoHeader: info,
			"error":    errorMessage,
		}).Fatal(message)

	default:
		logger.WithFields(logrus.Fields{
			"ID":          "d83326f4-1b06-4d13-8010-70b6c829cc88",
			"messageType": messageType,
		}).Fatal("Unknown messageType, stopping gateway")

	}
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

		LogErrorAndSendInfoToFenix(
			id,
			gRPC.InformationMessage_FATAL,
			"",
			"",
			returnDBMessage.err.Error(),
			"Got an error when Saveing to local DB, Stopping Gateway",
		)
	}
}

// *********************************************************************************
// Signal when channel reached certain levels
//

func channelSinaling(
	numberOfMessagesInChannel int,
	channelName string,
	id string) {
	// if channel items is greater than upper warning levels then lg warning
	if numberOfMessagesInChannel >= CHANNEL_UPPER_MESSAGE_TO_BE_SIGNALED {

		LogErrorAndSendInfoToFenix(
			id,
			gRPC.InformationMessage_WARNING,
			"Number of message in Channel '"+channelName+"'",
			string(numberOfMessagesInChannel),
			"",
			"Upper WARNING level reached for number of items in channel")
	} else {
		if numberOfMessagesInChannel >= CHANNEL_LOWER_MESSAGE_TO_BE_SIGNALED {
			logger.WithFields(logrus.Fields{
				"ID": id,
				"Number of message in Channel '" + channelName + "'": numberOfMessagesInChannel,
			}).Info("Lower INFO level reached for number of items in channel")
		} else {
			logger.WithFields(logrus.Fields{
				"ID": id,
				"Number of message in Channel '" + channelName + "'": numberOfMessagesInChannel,
			}).Debug("Number of items in channel")
		}
	}
}
