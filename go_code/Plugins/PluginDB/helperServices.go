package PluginKeyValueDBStore

import (
	"encoding/json"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

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
	dbMessageQueue = make(chan common_code.DbMessageStruct, common_code.SuppertedNumberOfMessagesInChannels)

}
