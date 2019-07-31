package TestExecutionGateway

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
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
		common_code.Logger.WithFields(logrus.Fields{
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
	clientsAddressAndPortInfo, keyExists := common_code.ClientsAddressAndPort[pluginId]

	if keyExists == true {
		addressAndPort = clientsAddressAndPortInfo.clientAddress + strconv.FormatInt(int64(clientsAddressAndPortInfo.clientPort), 10)

	} else {
		// Key not found, send message to Fenix about Error
		common_code.Logger.WithFields(logrus.Fields{
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
	ParentAddressAndPortInfo := common_code.GatewayConfig.ParentgRPCAddress

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
	returnParentAddressChannel := make(chan common_code.DbResultMessageStruct)

	// Get Clients address
	dbMessage := common_code.dbMessageStruct{
		common_code.DbRead,
		common_code.BucketForParentAddress,
		common_code.BucketKeyForParentAddress,
		nil,
		returnParentAddressChannel}

	// Send Read message to database to receive address
	common_code.dbMessageQueue <- dbMessage
	// Wait for address from channel, then close the channel
	databaseReturnMessage := <-returnParentAddressChannel
	close(returnParentAddressChannel)

	// Check if an error occured
	if databaseReturnMessage.err != nil {
		// Error when reading database
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                         "1a2ce5e8-c872-4fbb-aeb4-df6828f372d3",
			"databaseReturnMessage.err,": databaseReturnMessage.err,
			"databaseReturnMessage.key":  databaseReturnMessage.key,
		}).Warning("No Bucket found or No Key found when reading database")
	} else {

		// Convert saved json object into Go-struct
		err = json.Unmarshal(databaseReturnMessage.value, &parentAddress)
		if err != nil {
			// Problem with unmarshal the json object
			common_code.Logger.WithFields(logrus.Fields{
				"ID":                     "261a6391-abc0-4d9f-a59f-3d0a67e5e52c",
				"databaseReturnMessage,": databaseReturnMessage,
			}).Error("Can't unmarshal gRPCParent-address object from database")

			// Send FATAL information to Fenix
			common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
				OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
				OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
				SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
				SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
				MessageId:                generateUUID(),
				MessageType:              gRPC.InformationMessage_FATAL,
				Message:                  "Can't unmarshal gRPCParent-address object from database",
				OrginalCreateDateTime:    generaTimeStampUTC(),
				OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
				OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
			}
		} else {

			// If Saved data differs from memory data then change in memory object
			// First check ip address
			if parentAddress.GatewayAddress != common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress {
				common_code.Logger.WithFields(logrus.Fields{
					"ID":                           "560c2b17-c71e-45dd-9a38-a3dfd1a2bbd6",
					"parentAddress.GatewayAddress": parentAddress.GatewayAddress,
					"gatewayConfig.ParentgRPCAddress.parentGatewayInitialServer_address": common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress,
				}).Info("Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

				//Send Warning information to Fenix
				common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
					OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
					OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
					SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
					SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
					MessageId:                generateUUID(),
					MessageType:              gRPC.InformationMessage_INFO,
					Message:                  "Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
					OrginalCreateDateTime:    generaTimeStampUTC(),
					OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
					OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
				}

				// Change Address in memory object
				common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress = parentAddress.GatewayAddress
			}

			// Second check port
			if parentAddress.GatewayPort != common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerPort {
				common_code.Logger.WithFields(logrus.Fields{
					"ID":                        "50a3b7ad-6631-42c5-ab5c-777e04ad9728",
					"parentAddress.GatewayPort": parentAddress.GatewayPort,
					"gatewayConfig.ParentgRPCAddress.parentGatewayInitialServer_port": common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerPort,
				}).Info("Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")

				//Send Warning information to Fenix
				common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
					OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
					OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
					SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
					SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
					MessageId:                generateUUID(),
					MessageType:              gRPC.InformationMessage_WARNING,
					Message:                  "Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version",
					OrginalCreateDateTime:    generaTimeStampUTC(),
					OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
					OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
				}

				// Change Port in memory object
				common_code.GatewayConfig.ParentgRPCAddress.ParentGatewayServerPort = parentAddress.GatewayPort
			}

		}
	}
}

// *********************************************************************************
// Log message to local log and then Send message to Fenix
//TODO SPlit logs into one log per messageType
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
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
		}).Debug(messageToFenix)

	case gRPC.InformationMessage_INFO:
		// Log information
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Info(messageToFenix)

		// Send information to Fenix
		common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_INFO,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_WARNING:
		// Log information
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Warning(messageToFenix)

		// Send Warning information to Fenix
		common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_WARNING,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_ERROR:
		// Log information
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Error(messageToFenix)

		// Send Error information to Fenix
		common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         common_code.GatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       common_code.GatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 common_code.GatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               common_code.GatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   common_code.GatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: common_code.GatewayConfig.SystemDomain.GatewayDomainName,
		}

	case gRPC.InformationMessage_FATAL:
		// Only log and then Terminate Gateway due to problems
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                 id,
			infoHeaderForLogging: infoForLogging,
			"error":              errorMessageForLogging,
		}).Fatal(messageToFenix)

	default:
		common_code.Logger.WithFields(logrus.Fields{
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
	returnChannel := make(chan common_code.DbResultMessageStruct)

	dbMessage := common_code.dbMessageStruct{
		common_code.DbWrite,
		bucket,
		key,
		valueToStoreInDB,
		returnChannel}

	// Send message to Database
	common_code.dbMessageQueue <- dbMessage

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
	if numberOfMessagesInChannel >= common_code.UpperBounderyForSignalingMessagesInChannel {

		LogErrorAndSendInfoToFenix(
			id,
			gRPC.InformationMessage_WARNING,
			"Number of message in Channel '"+channelName+"'",
			string(numberOfMessagesInChannel),
			"",
			"Upper WARNING level reached for number of items in channel")
	} else {
		if numberOfMessagesInChannel >= common_code.LowerBounderyForSignalingMessagesInChannel {
			common_code.Logger.WithFields(logrus.Fields{
				"ID": id,
				"Number of message in Channel '" + channelName + "'": numberOfMessagesInChannel,
			}).Info("Lower INFO level reached for number of items in channel")
		} else {
			common_code.Logger.WithFields(logrus.Fields{
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

// Private
func initiateClientAddressMemoryDB() {
	common_code.ClientsAddressAndPort = make(map[string]common_code.ClientsAddressAndPortStruct)
}

// Public
func InitiateClientAddressMemoryDB() {
	initiateClientAddressMemoryDB()
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

		common_code.Logger.WithFields(logrus.Fields{
			"ID":        "291e7759-a002-4363-bbfa-3604bd3255a0",
			"ipAddress": ipAddress,
		}).Debug("Child gateway/Plugin ip-address is the same as this gateway's ip address, will use '127.0.0.1' as key in memory DB.")

		ipAddress = "127.0.0.1"

	}

	// Try to get vaule from key and check if it exits
	clientAddresAndPortInfo, ifExists := common_code.ClientsAddressAndPort[ipAddress]
	if ifExists == true {
		// IP-address exists, Add +1 to current port-counter and save in memory DB
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                      "8800e41f-d46c-4857-b7f3-8ec23f904f77",
			"ipAddress":               ipAddress,
			"clientAddresAndPortInfo": clientAddresAndPortInfo,
		}).Debug("Current ip-address exists in memory DB.")

		clientAddresAndPortInfo.clientPort = clientAddresAndPortInfo.clientPort + 1
		common_code.ClientsAddressAndPort[ipAddress] = clientAddresAndPortInfo
	} else {
		// IP-address does not exists, Use start Port from config and save in memory DB
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                      "35528300-e5ee-4c0a-b41a-15a5d37d8186",
			"ipAddress":               ipAddress,
			"clientAddresAndPortInfo": common_code.GatewayConfig.InitialClientPort.InitialClientPort,
		}).Debug("IP-address not used before, will use standard start port")

		clientAddresAndPortInfo.clientPort = common_code.GatewayConfig.InitialClientPort.InitialClientPort
		common_code.ClientsAddressAndPort[ipAddress] = clientAddresAndPortInfo

	}

	return clientAddresAndPortInfo.clientPort
}

// *********************************************************************************
//  Get preferred outbound ip of this machine
//
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		common_code.Logger.WithFields(logrus.Fields{
			"ID":                      "118517f0-6839-49b5-9f8f-850e82d3de23",
			"err":                     err,
			"clientAddresAndPortInfo": common_code.GatewayConfig.InitialClientPort.InitialClientPort,
		}).Fatal("Error when getting this gateways own IP-address")
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			common_code.Logger.WithFields(logrus.Fields{
				"ID":  "1805f803f-6bfd-49d0-82b6-947677a539bb",
				"err": err,
			}).Fatal("Fatal ERROR when closing udp-connection when checking local IP-address in function 'GetOutboundIP()'")
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	common_code.Logger.WithFields(logrus.Fields{
		"ID":        "2a2a4450-2a09-4084-bc1b-c5ac6ee58bcb",
		"localAddr": localAddr,
	}).Debug("This gateways own address info")

	return localAddr.IP.String()
}

// *********************************************************************************
//  Initiate channels used within gatway
//
func InitiateGatewayChannels() {

	common_code.Logger.WithFields(logrus.Fields{
		"ID": "",
	}).Debug("Initiate local gateway channels")

	// Databasechannel for reading and Writing
	common_code.dbMessageQueue = make(chan common_code.dbMessageStruct, common_code.SuppertedNumberOfMessagesInChannels)

	// *** Towards Fenix ***
	// informationMessage-channel towards Fenix
	common_code.informationMessageChannelTowardsFenix = make(chan *gRPC.InformationMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// testInstructionTimeOutMessage-channel towards Fenix
	common_code.testInstructionTimeOutMessageChannelTowardsFenix = make(chan *gRPC.TestInstructionTimeOutMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// testExecutionLogMessage-channel towards Fenix
	common_code.testExecutionLogMessageChannelTowardsFenix = make(chan *gRPC.TestExecutionLogMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsMessage-channel towards Fenix
	common_code.supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix = make(chan *gRPC.SupportedTestDataDomainsMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// availbleTestInstruction<AtPluginMessage-channel towards Fenix
	common_code.availbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix = make(chan *gRPC.AvailbleTestInstructionAtPluginMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// *** Towards Plugina ***
	// TestInstruction Towards Plugin
	common_code.testInstructionMessageChannelTowardsPlugin = make(chan *gRPC.TestInstruction_RT, common_code.SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsRequest Towards Plugin
	common_code.supportedTestDataDomainsRequestChannelTowardsPlugin = make(chan *gRPC.SupportedTestDataDomainsRequest, common_code.SuppertedNumberOfMessagesInChannels)

}
