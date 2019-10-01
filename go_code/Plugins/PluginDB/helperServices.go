package PluginKeyValueDBStore

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	"github.com/sirupsen/logrus"
)

// *********************************************************************************
// Signal when channel reached certain levels
//

func channelSinaling(
	numberOfMessagesInChannel int,
	channelName string,
	id string) {

	// if channel items is greater than upper warning levels then lg warning
	if numberOfMessagesInChannel >= common_code.UpperBounderyForSignalingMessagesInChannel {

		logger.WithFields(logrus.Fields{
			"ID": id,
			"Number of message in Channel '" + channelName + "'": numberOfMessagesInChannel,
		}).Warn("Upper WARNING level reached for number of items in channel")

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
//  Initiate channels used within gatway
//
func initiateKeyValueStoreDbChannels() {

	logger.WithFields(logrus.Fields{
		"ID": "1d0e73b4-03ef-4f9f-82ec-db3207ec7508",
	}).Debug("Initiate local read/write channels")

	// Databasechannel for reading
	dbReadMessageChannel = make(chan DbReadMessageStruct, suppertedNumberOfMessagesInChannels)

	// Databasechannel for writing
	dbWritewMessageChannel = make(chan DbWriteMessageStruct, suppertedNumberOfMessagesInChannels)

}
