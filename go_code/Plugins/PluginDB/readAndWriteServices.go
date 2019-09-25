package PluginKeyValueDBStore

import (
	"errors"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/sirupsen/logrus"
)

// **********************************************************************************************************
// Initiate Service for checking Read-channel and Write-channel go-routines
//

func initiateReadANdWriteServices() {

	// Start Read-service
	go readFromDbService()

	// Start Write-service
	go writeToDbService()
}

// **********************************************************************************************************
// Service for checking channel incoming 'WriteKeyValueMessage'. Is run as a go-routine
//
func writeToDbService() {

	var returnMessage DbResultWriteMessageStruct

	// Loop and wait for incoming message from channel
	for {
		keyValueToBeSaved := <-dbWritewMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":                "26782f27-2001-43fb-8e03-f1ce406aaf5a",
			"keyValueToBeSaved": keyValueToBeSaved,
		}).Debug("Message to be saved in KeyValueStore-DB was read from channel")

		// Check number of waiting message in channel
		channelSinaling(
			len(dbWritewMessageChannel),
			"dbWritewMessageChannel",
			"7d942a72-b2ba-4f69-b9ca-1d69f2fec680")

		messageSavedInDB := saveKeyValuetMessageInDB(&keyValueToBeSaved.WriteMessage)

		if messageSavedInDB == true {
			// Message saved in KeyValueStore

			// Create returnMessgae
			returnMessage = DbResultWriteMessageStruct{
				Err:    nil,
				Key:    keyValueToBeSaved.WriteMessage.Key,
				Bucket: keyValueToBeSaved.WriteMessage.Bucket,
			}

		} else {
			// Message NOT saved in KeyValueStore

			// Create returnMessgae
			returnMessage = DbResultWriteMessageStruct{
				Err:    errors.New("Couldn't save Valye in KeyValueStore-DB"),
				Key:    keyValueToBeSaved.WriteMessage.Key,
				Bucket: keyValueToBeSaved.WriteMessage.Bucket,
			}
		}

		// Put returnMessage on channel back to sender
		keyValueToBeSaved.ResultsChannel <- returnMessage
	}

}

// **********************************************************************************************************
// Service for checking channel incoming 'ReadKeyMessage'. Is run as a go-routine
//
func readFromDbService() {

	var returnMessage DbResultReadMessageStruct

	// Loop and wait for incoming message from channel
	for {
		keyToBeRead := <-dbReadMessageChannel

		logger.WithFields(logrus.Fields{
			"ID":          "aa032c6f-fbfe-4168-8a9e-e5e7f0088729",
			"keyToBeRead": keyToBeRead,
		}).Debug("Message used to read KeyValueStore-DB was read from channel")

		// Check number of waiting message in channel
		channelSinaling(
			len(dbReadMessageChannel),
			"dbWritewMessageChannel",
			"2258db93-2b91-4dce-a812-542e33cc7f61")

		messageSavedInDB, err := readKeyValuetMessageFromDB(keyToBeRead.ReadMessage)

		if messageSavedInDB.Acknack == true && err == nil {
			// Message saved in KeyValueStore

			// Create returnMessgae
			returnMessage = DbResultReadMessageStruct{
				err,
				messageSavedInDB.GetKey(),
				messageSavedInDB.GetBucket(),
				int(messageSavedInDB.GetValueSaveType()),
				gRPC.ValueSaveTypeEnum_name[int32(messageSavedInDB.GetValueSaveType())],
				messageSavedInDB.GetValue(),
				messageSavedInDB.GetValueString(),
				messageSavedInDB.GetUpdatedDateTime(),
			}

		} else {
			// Message NOT saved in KeyValueStore

			// Create returnMessgae
			returnMessage = DbResultReadMessageStruct{
				Err:    errors.New("Couldn't read key in KeyValueStore-DB"),
				Key:    keyToBeRead.ReadMessage.Key,
				Bucket: keyToBeRead.ReadMessage.Bucket,
			}
		}

		// Put returnMessage on channel back to sender
		keyToBeRead.ResultsChannel <- returnMessage
	}

}
