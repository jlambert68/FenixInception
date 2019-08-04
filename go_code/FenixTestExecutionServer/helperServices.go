package FenixTestExecutionServer

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

// *********************************************************************************
// Get Clients IP-address and Port from final PluginId from Memory object
//
func getClientAddressAndPort(pluginId string) (addressAndPort string) {
	// Get Client-info object from Clients-memory-object
	clientsAddressAndPortInfo, keyExists := clientsAddressAndPort[pluginId]

	if keyExists == true {
		addressAndPort = clientsAddressAndPortInfo.ClientAddress + strconv.FormatInt(int64(clientsAddressAndPortInfo.ClientPort), 10)

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
// Log message to local log and then Send message to Fenix
//TODO SPlit logs into one log per messageType
// TODO change into 'func LogErrorAndSaveTOFenixDB('
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
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_INFO,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    common_code.GeneraTimeStampUTC(),
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
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_WARNING,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    common_code.GeneraTimeStampUTC(),
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
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  messageToFenix,
			OrginalCreateDateTime:    common_code.GeneraTimeStampUTC(),
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
	returnChannel := make(chan common_code.DbResultMessageStruct)

	dbMessage := common_code.DbMessageStruct{
		common_code.DbWrite,
		bucket,
		key,
		valueToStoreInDB,
		returnChannel}

	// Send message to Database
	dbMessageQueue <- dbMessage

	// Wait for result on result channel and then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.Err != nil {

		LogErrorAndSendInfoToFenix(
			id,
			gRPC.InformationMessage_FATAL,
			"",
			"",
			returnDBMessage.Err.Error(),
			"Got an error when Saveing to local DB, Stopping Gateway",
		)
		return false
	}
	return true
}

// *********************************************************************************
// Signal when channel reached certain levels
// TODO Move this to common_code

func channelSinaling(
	numberOfMessagesInChannel int,
	channelName string,
	id string,
	logger *logrus.Logger) {
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
// Initiate map for clients addresses

// Private
func initiateClientAddressMemoryDB() {
	clientsAddressAndPort = make(map[string]common_code.ClientsAddressAndPortStruct)
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

		clientAddresAndPortInfo.ClientPort = clientAddresAndPortInfo.ClientPort + 1
		clientsAddressAndPort[ipAddress] = clientAddresAndPortInfo
	} else {
		// IP-address does not exists, Use start Port from config and save in memory DB
		logger.WithFields(logrus.Fields{
			"ID":                      "35528300-e5ee-4c0a-b41a-15a5d37d8186",
			"ipAddress":               ipAddress,
			"clientAddresAndPortInfo": gatewayConfig.InitialClientPort.InitialClientPort,
		}).Debug("IP-address not used before, will use standard start port")

		clientAddresAndPortInfo.ClientPort = gatewayConfig.InitialClientPort.InitialClientPort
		clientsAddressAndPort[ipAddress] = clientAddresAndPortInfo

	}

	return clientAddresAndPortInfo.ClientPort
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
func InitiateGatewayChannels() {

	logger.WithFields(logrus.Fields{
		"ID": "",
	}).Debug("Initiate local gateway channels")

	// Databasechannel for reading and Writing
	dbMessageQueue = make(chan common_code.DbMessageStruct, common_code.SuppertedNumberOfMessagesInChannels)

	// *** Towards Fenix ***
	// informationMessage-channel towards Fenix
	gatewayChannelPackage.InformationMessageChannelTowardsFenix = make(chan *gRPC.InformationMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// testInstructionTimeOutMessage-channel towards Fenix
	gatewayChannelPackage.TestInstructionTimeOutMessageChannelTowardsFenix = make(chan *gRPC.TestInstructionTimeOutMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// testExecutionLogMessage-channel towards Fenix
	gatewayChannelPackage.TestExecutionLogMessageChannelTowardsFenix = make(chan *gRPC.TestExecutionLogMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsMessage-channel towards Fenix
	gatewayChannelPackage.SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix = make(chan *gRPC.SupportedTestDataDomainsMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// availbleTestInstruction<AtPluginMessage-channel towards Fenix
	gatewayChannelPackage.AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix = make(chan *gRPC.AvailbleTestInstructionAtPluginMessage, common_code.SuppertedNumberOfMessagesInChannels)

	// *** Towards Plugina ***
	// TestInstruction Towards Plugin
	gatewayChannelPackage.TestInstructionMessageChannelTowardsPlugin = make(chan *gRPC.TestInstruction_RT, common_code.SuppertedNumberOfMessagesInChannels)

	// supportedTestDataDomainsRequest Towards Plugin
	gatewayChannelPackage.SupportedTestDataDomainsRequestChannelTowardsPlugin = make(chan *gRPC.SupportedTestDataDomainsRequest, common_code.SuppertedNumberOfMessagesInChannels)

}
