package TestExecutionGateway

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"strconv"
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

	addressAndPort = ParentAddressAndPortInfo.parentGatewayServer_address + ":" + strconv.Itoa(ParentAddressAndPortInfo.parentGatewayServer_port)

	return addressAndPort

}

// *******************************************************************
// Replace Parent gateway/Fenix IP-address & port info in memory object if previous connection, saved in database, differs from config-file
//
func updateMemoryAddressForParentAddressInfo() {

	var parentAddress gRPC.ReRegisterToGatewayMessage
	var err error

	// Create the channel that the client address should be sent back on
	returnParentAddressChannel := make(chan dbResultMessage_struct)

	// Get Clients address
	dbMessage := dbMessage_struct{
		DB_READ,
		BUCKET_PARENT_ADDRESS,
		BUCKET_KEY_PARENT_ADDRESS,
		nil,
		returnParentAddressChannel}

	// Send Read message to database to receive address
	dbMessageQueue <- dbMessage

	// Wait for address from channel, then close the channel
	parentAddressByteArray := <-returnParentAddressChannel
	close(returnParentAddressChannel)

	// Convert saved json object into Go-struct
	err = json.Unmarshal(parentAddressByteArray.value, &parentAddress)
	if err != nil {
		// Problem with unmarshal the json object
		logger.WithFields(logrus.Fields{
			"ID":                      "261a6391-abc0-4d9f-a59f-3d0a67e5e52c",
			"parentAddressByteArray,": parentAddressByteArray,
		}).Error("Can't unmarshal gRPCParent-address object from database")

		// Send FATAL information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
			MessageId:             generateUUID(),
			MessageType:           gRPC.InformationMessage_FATAL,
			Message:               "Can't unmarshal gRPCParent-address object from database",
			OrginalCreateDateTime: generaTimeStampUTC(),
		}
	} else {

		// If Saved data differs from memory data then change in memory object
		// First check ip address
		if parentAddress.GatewayAddress != gatewayConfig.parentgRPCAddress.parentGatewayServer_address {
			logger.WithFields(logrus.Fields{
				"ID":                           "560c2b17-c71e-45dd-9a38-a3dfd1a2bbd6",
				"parentAddress.GatewayAddress": parentAddress.GatewayAddress,
				"gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address": gatewayConfig.parentgRPCAddress.parentGatewayServer_address,
			}).Info("Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

			//Send Warning  information to Fenix
			localInformationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
				OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
				SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
				SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
				MessageId:             generateUUID(),
				MessageType:           gRPC.InformationMessage_INFO,
				Message:               "Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

			// Change Address in memory object
			gatewayConfig.parentgRPCAddress.parentGatewayServer_address = parentAddress.GatewayAddress
		}

		// Second check port
		if parentAddress.GatewayPort != gatewayConfig.parentgRPCAddress.parentGatewayServer_port {
			logger.WithFields(logrus.Fields{
				"ID":                        "50a3b7ad-6631-42c5-ab5c-777e04ad9728",
				"parentAddress.GatewayPort": parentAddress.GatewayPort,
				"gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port": gatewayConfig.parentgRPCAddress.parentGatewayServer_port,
			}).Info("Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

			//Send Warning information to Fenix
			localInformationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
				OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
				SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
				SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
				MessageId:             generateUUID(),
				MessageType:           gRPC.InformationMessage_WARNING,
				Message:               "Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
				OrginalCreateDateTime: generaTimeStampUTC(),
			}

			// Change Port in memory object
			gatewayConfig.parentgRPCAddress.parentGatewayServer_port = parentAddress.GatewayPort
		}

	}
}

// *******************************************************************
// Replace Parent gateway/Fenix IP-address & port info in databse, from Memory if previous connection differs from memory object
//

func updateDatabaseFromMemoryForParentAddressInfo(reRegisterToGatewayMessage gRPC.ReRegisterToGatewayMessage) {

	// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
	reRegisterToGatewayMessageToSavedAsByteArray, err := json.Marshal(reRegisterToGatewayMessage)

	if err != nil {
		// Error when Unmarshaling to []byte
		LogErrorAndSendInfoToFenix(
			"6366259e-a206-42b1-b99a-d4a18a56da96",
			gRPC.InformationMessage_FATAL,
			"reRegisterToGatewayMessage",
			reRegisterToGatewayMessage.String(),
			err.Error(),
			"Error when converting 'reRegisterToGatewayMessage' into a byte array, stopping processing and stopping Gateway",
		)

	} else {
		// Marshaling to []byte OK

		// Save message to local DB for later processing
		_ = SaveMessageToLocalDB(
			BUCKET_KEY_PARENT_ADDRESS,
			reRegisterToGatewayMessageToSavedAsByteArray,
			BUCKET_PARENT_ADDRESS,
			"f1ae7544-a190-4e36-b527-5abdd86c0c61",
		)
	}
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
			OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
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
			OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
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
			OriginalSenderId:      gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:    gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:              gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:            gatewayConfig.gatewayIdentification.gatewayName,
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
) (saveOK bool) {

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
		return false
	}
	return true
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
