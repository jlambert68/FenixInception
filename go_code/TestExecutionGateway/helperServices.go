package TestExecutionGateway

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"net"
	"strconv"
	"time"
)

// *********************************************************************************
// Generate a new unique uuid
//
func generateUUID() string {
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
func generaTimeStampUTC() string {
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
// Get Clients IP-address and Port from final PluginId from Memory object
//
func getClientAddressAndPort(pluginId string) (addressAndPort string) {
	// Get Client-info object from Clients-memory-object
	clientsAddressAndPortInfo, keyExists := clientsAddressAndPort[pluginId]

	if keyExists == true {
		addressAndPort = clientsAddressAndPortInfo.clientAddress + strconv.FormatInt(int64(clientsAddressAndPortInfo.clientPort), 10)

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
	ParentAddressAndPortInfo := gatewayConfig.ParentgRPCAddress

	addressAndPort = ParentAddressAndPortInfo.ParentGatewayServerAddress + ":" + strconv.FormatInt(int64(ParentAddressAndPortInfo.ParentGatewayServerPort), 10)

	return addressAndPort

}

// *******************************************************************
// Replace Parent gateway/Fenix IP-address & port info in memory object if previous connection, saved in database, differs from config-file
//
func updateMemoryAddressForParentAddressInfo() {

	var parentAddress gRPC.ReRegisterToGatewayMessage
	var err error

	// Create the channel that the client address should be sent back on
	returnParentAddressChannel := make(chan dbResultMessageStruct)

	// Get Clients address
	dbMessage := dbMessageStruct{
		DbRead,
		BucketForParentAddress,
		BucketKeyForParentAddress,
		nil,
		returnParentAddressChannel}

	// Send Read message to database to receive address
	dbMessageQueue <- dbMessage
	// Wait for address from channel, then close the channel
	databaseReturnMessage := <-returnParentAddressChannel
	close(returnParentAddressChannel)

	// Check if an error occured
	if databaseReturnMessage.err != nil {
		// Error when reading database
		logger.WithFields(logrus.Fields{
			"ID":                         "1a2ce5e8-c872-4fbb-aeb4-df6828f372d3",
			"databaseReturnMessage.err,": databaseReturnMessage.err,
			"databaseReturnMessage.key":  databaseReturnMessage.key,
		}).Warning("No Bucket found or No Key found when reading database")
	} else {

		// Convert saved json object into Go-struct
		err = json.Unmarshal(databaseReturnMessage.value, &parentAddress)
		if err != nil {
			// Problem with unmarshal the json object
			logger.WithFields(logrus.Fields{
				"ID":                     "261a6391-abc0-4d9f-a59f-3d0a67e5e52c",
				"databaseReturnMessage,": databaseReturnMessage,
			}).Error("Can't unmarshal gRPCParent-address object from database")

			// Send FATAL information to Fenix
			gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
				OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
				OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
				SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
				SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
				MessageId:                generateUUID(),
				MessageType:              gRPC.InformationMessage_FATAL,
				Message:                  "Can't unmarshal gRPCParent-address object from database",
				OrginalCreateDateTime:    generaTimeStampUTC(),
				OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
				OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
			}
		} else {

			// If Saved data differs from memory data then change in memory object
			// First check ip address
			if parentAddress.GatewayAddress != gatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress {
				logger.WithFields(logrus.Fields{
					"ID":                           "560c2b17-c71e-45dd-9a38-a3dfd1a2bbd6",
					"parentAddress.GatewayAddress": parentAddress.GatewayAddress,
					"gatewayConfig.ParentgRPCAddress.parentGatewayInitialServer_address": gatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress,
				}).Info("Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

				//Send Warning information to Fenix
				gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
					OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
					SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
					SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
					MessageId:                generateUUID(),
					MessageType:              gRPC.InformationMessage_INFO,
					Message:                  "Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
					OrginalCreateDateTime:    generaTimeStampUTC(),
					OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
					OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
				}

				// Change Address in memory object
				gatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress = parentAddress.GatewayAddress
			}

			// Second check port
			if parentAddress.GatewayPort != gatewayConfig.ParentgRPCAddress.ParentGatewayServerPort {
				logger.WithFields(logrus.Fields{
					"ID":                        "50a3b7ad-6631-42c5-ab5c-777e04ad9728",
					"parentAddress.GatewayPort": parentAddress.GatewayPort,
					"gatewayConfig.ParentgRPCAddress.parentGatewayInitialServer_port": gatewayConfig.ParentgRPCAddress.ParentGatewayServerPort,
				}).Info("Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

				//Send Warning information to Fenix
				gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
					OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
					SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
					SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
					MessageId:                generateUUID(),
					MessageType:              gRPC.InformationMessage_WARNING,
					Message:                  "Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
					OrginalCreateDateTime:    generaTimeStampUTC(),
					OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
					OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
				}

				// Change Port in memory object
				gatewayConfig.ParentgRPCAddress.ParentGatewayServerPort = parentAddress.GatewayPort
			}

		}
	}
}

// *******************************************************************
// Replace Parent gateway/Fenix IP-address & port info in databse, from Memory if previous connection differs from memory object
//
//NOT USED (190612)
/*
func updateDatabaseFromMemoryForParentAddressInfo_ShouldNotBeUsed(GatewayIdentification GatewayIdentificationStruct) {

	// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
	reRegisterToGatewayMessageToSavedAsByteArray, err := json.Marshal(GatewayIdentification)

	if err != nil {
		// Error when Unmarshaling to []byte
		LogErrorAndSendInfoToFenix(
			"6366259e-a206-42b1-b99a-d4a18a56da96",
			gRPC.InformationMessage_FATAL,
			"GatewayIdentification",
			"no data available...",
			err.Error(),
			"Error when converting 'GatewayIdentification' into a byte array, stopping processing and stopping Gateway",
		)

	} else {
		// Marshaling to []byte OK

		// Save message to local DB for later processing
		_ = SaveMessageToLocalDB(
			BucketKeyForParentAddress,
			reRegisterToGatewayMessageToSavedAsByteArray,
			BucketForParentAddress,
			"f1ae7544-a190-4e36-b527-5abdd86c0c61",
		)
	}
}*/

// *********************************************************************************
// Log message to local log and then Send message to Fenix
func LogErrorAndSendInfoToFenix(
	id string,
	messageType gRPC.InformationMessage_InformationType,
	infoHeaderForLogging string,
	infoForLogging string,
	errorMessageForLogging string,
	messageToFenix string,
) {

	switch messageType {
	case gRPC.InformationMessage_DEBUG:
		// Only logg information and do not send to Fenix
		logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
		}).Debug(messageToFenix)

	case gRPC.InformationMessage_INFO:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Info(messageToFenix)

		// Send information to Fenix
		gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_INFO,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_WARNING:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Warning(messageToFenix)

		// Send Warning information to Fenix
		gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_WARNING,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_ERROR:
		// Log information
		logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Error(messageToFenix)

		// Send Error information to Fenix
		gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_FATAL:
		// Only log and then Terminate Gateway due to problems
		logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Fatal(messageToFenix)

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
	returnChannel := make(chan dbResultMessageStruct)

	dbMessage := dbMessageStruct{
		DbWrite,
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
	if numberOfMessagesInChannel >= UpperBounderyForSignalingMessagesInChannel {

		LogErrorAndSendInfoToFenix(
			id,
			gRPC.InformationMessage_WARNING,
			"Number of message in Channel '"+channelName+"'",
			string(numberOfMessagesInChannel),
			"",
			"Upper WARNING level reached for number of items in channel")
	} else {
		if numberOfMessagesInChannel >= LowerBounderyForSignalingMessagesInChannel {
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

// *********************************************************************************
// Get highest gRPC version
//

func getHighestGRPCVersion() (currentVersion string) {
	maxVersionCount := int32(len(gRPC.CurrentVersionEnum_name) - 1)

	maxVersion := gRPC.CurrentVersionEnum_name[maxVersionCount]
	return maxVersion
}

// *********************************************************************************
// Initiate map for clients addresses
//
func initiateClientAddressMemoryDB() {
	clientsAddressAndPort = make(map[string]clientsAddressAndPortStruct)
}

// *********************************************************************************
// Get next free Client port
//

func getNextFreeClientPort(ipAddress string) (port int32) {
	// If decide if incoming ip-address is localhost, 127.0.0.1, 0.0.0.0, localhost, 192.168.2.x="this gateways ip)

	var thisGataewaysIpAddress string
	var childUseSamteIp bool

	// Get this gateways IP-address
	thisGataewaysIpAddress = GetOutboundIP()

	// Decide if child-gateway use the IP-address as this gateway
	switch ipAddress {
	case "127.0.0.1":
		childUseSamteIp = true
	case "0.0.0.0":
		childUseSamteIp = true
	case "localhost":
		childUseSamteIp = true
	case thisGataewaysIpAddress:
		childUseSamteIp = true
	default:
		childUseSamteIp = false
	}

	// If child's address is the same as gateway address then use 127.0.0.1 as address as key
	if childUseSamteIp == true {

		logger.WithFields(logrus.Fields{
			"ID":        "291e7759-a002-4363-bbfa-3604bd3255a0",
			"ipAddress": ipAddress,
		}).Debug("Child gateway/Plugin ip-address is the same as this gateway's ip address, will use '127.0.0.1' as key in memory DB.")

		ipAddress = "127.0.0.1"

	}

	// Try to get vaule from key and check if it exits
	clientAddresAndPortInfo, ifExists := clientsAddressAndPort[ipAddress]
	if ifExists == true {
		// IP-address exists, Add +1 to current port-counter and save in memory DB
		logger.WithFields(logrus.Fields{
			"ID":                      "8800e41f-d46c-4857-b7f3-8ec23f904f77",
			"ipAddress":               ipAddress,
			"clientAddresAndPortInfo": clientAddresAndPortInfo,
		}).Debug("Current ip-address exists in memory DB.")

		clientAddresAndPortInfo.clientPort = clientAddresAndPortInfo.clientPort + 1
		clientsAddressAndPort[ipAddress] = clientAddresAndPortInfo
	} else {
		// IP-address does not exists, Use start Port from config and save in memory DB
		logger.WithFields(logrus.Fields{
			"ID":                      "35528300-e5ee-4c0a-b41a-15a5d37d8186",
			"ipAddress":               ipAddress,
			"clientAddresAndPortInfo": gatewayConfig.InitialClientPort.InitialClientPort,
		}).Debug("IP-address not used before, will use standard start port")

		clientAddresAndPortInfo.clientPort = gatewayConfig.InitialClientPort.InitialClientPort
		clientsAddressAndPort[ipAddress] = clientAddresAndPortInfo

	}

	return clientAddresAndPortInfo.clientPort
}

// *********************************************************************************
//  Get preferred outbound ip of this machine
//
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":                      "118517f0-6839-49b5-9f8f-850e82d3de23",
			"err":                     err,
			"clientAddresAndPortInfo": gatewayConfig.InitialClientPort.InitialClientPort,
		}).Fatal("Error when getting this gateways own IP-address")
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":  "1805f803f-6bfd-49d0-82b6-947677a539bb",
				"err": err,
			}).Fatal("Fatal ERROR when closing udp-connection when checking local IP-address in function 'GetOutboundIP()'")
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	logger.WithFields(logrus.Fields{
		"ID":        "2a2a4450-2a09-4084-bc1b-c5ac6ee58bcb",
		"localAddr": localAddr,
	}).Debug("This gateways own address info")

	return localAddr.IP.String()
}

// *********************************************************************************
//  Initiate channels used within gatway
//
func initiateGatewayChannels() {

	logger.WithFields(logrus.Fields{
		"ID": "",
	}).Debug("Initiate local gateway channels")

	// Databasechannel for reading and Writing
	dbMessageQueue = make(chan dbMessageStruct, SuppertedNumberOfMessagesInChannels)

	// *** Towards Fenix ***
	// informationMessage-channel towards Fenix
	gatewayTowardsFenixObject.informationMessageChannel = make(chan *gRPC.InformationMessage, SuppertedNumberOfMessagesInChannels)

	// testInstructionTimeOutMessage-channel towards Fenix
	gatewayTowardsFenixObject.testInstructionTimeOutMessageChannel = make(chan *gRPC.TestInstructionTimeOutMessage, SuppertedNumberOfMessagesInChannels)

	// testExecutionLogMessage-channel towards Fenix
	gatewayTowardsFenixObject.testExecutionLogMessageChannel = make(chan *gRPC.TestExecutionLogMessage, SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsMessage-channel towards Fenix
	gatewayTowardsFenixObject.supportedTestDataDomainsMessageTowardsFenixChannel = make(chan *gRPC.SupportedTestDataDomainsMessage, SuppertedNumberOfMessagesInChannels)

	// availbleTestInstruction<AtPluginMessage-channel towards Fenix
	gatewayTowardsFenixObject.availbleTestInstructionAtPluginMessageTowardsFenixChannel = make(chan *gRPC.AvailbleTestInstructionAtPluginMessage, SuppertedNumberOfMessagesInChannels)

	// *** Towards Plugina ***
	// TestInstruction Towards Plugin
	gatewayTowardsPluginObject.testInstructionMessageChannel = make(chan *gRPC.TestInstruction_RT, SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsRequest Towards Plugin
	gatewayTowardsPluginObject.supportedTestDataDomainsRequestChannel = make(chan *gRPC.SupportedTestDataDomainsRequest, SuppertedNumberOfMessagesInChannels)

}

/*
NOT USED

// *********************************************************************************
//  Update Memory object with parameters that was recceived at start up using flags
//
func updateMemoryObjectWithFlagOverrideParameters() {

	// **** As slave gateway for tests ****
	// THis gateway should listen on this port from calls from clients
	if GatewayInIntegrationTestMode.ListeningOnThisPortAsParent != 0 {
		gatewayConfig.GatewayIdentification.GatewayChildrenCallOnThisPort = int32(GatewayInIntegrationTestMode.ListeningOnThisPortAsParent)
	}

	// THis gateway should call parent gateway/Fenix on this adress
	if GatewayInIntegrationTestMode.parentIsListeningOnThisPort != 0 {
		gatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress = GatewayInIntegrationTestMode.parentIsListeningOnThisAddress
	}

	// **** As the tested gateway ****
	// THis gateway should call parent gateway/Fenix on this port
	if GatewayInIntegrationTestMode.parentIsListeningOnThisPort != 0 {
		gatewayConfig.ParentgRPCAddress.ParentGatewayServerPort = int32(GatewayInIntegrationTestMode.parentIsListeningOnThisPort)
	}
}
*/
