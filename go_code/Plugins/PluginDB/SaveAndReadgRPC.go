package PluginKeyValueDBStore

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)
WriteToKeyValueStore (WriteKeyValueMessage) returns (AckNackResponse) {
}


rpc ReadFromKeyValueStore (ReadKeyRequestMessage) returns (ValueResponseMessage) {
}
// ********************************************************************************************
// Call from Plugin that need to store Values in KeyValueStore
//
func (gRPCServerForDbKeyValueStore *gRPCServerForDbKeyValueStoreStruct) WriteToKeyValueStore(ctx context.Context, writeKeyValueMessage *gRPC.WriteKeyValueMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                   "c8ad7a99-a27f-4f77-be96-81733d57d36b",
		"writeKeyValueMessage": writeKeyValueMessage,
	}).Debug("Incoming gRPC: 'WriteToKeyValueStore'")

	// Put writeKeyValueMessage on queue for further processing
	dbWritewMessageChannel <- writeKeyValueMessage
	logger.WithFields(logrus.Fields{
		"ID": "0c11316c-ea0a-45f7-9e53-ecc842f40a16",
	}).Debug("'writeKeyValueMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "5821e625-2c0a-4b84-bcdf-ef99c4af2c1e",
	}).Debug("Leaving gRPC: 'WriteToKeyValueStore'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'writeKeyValueMessage' was saved in KeyValueStore"
	returnMessage.Acknack = true

	return returnMessage, nil
}

// ********************************************************************************************
// Call from Plugin that need to read Values in KeyValueStore
//
func (gRPCServerForDbKeyValueStore *gRPCServerForDbKeyValueStoreStruct) ReadFromKeyValueStore(ctx context.Context, readKeyRequestMessage *gRPC.ReadKeyRequestMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                    "3628b2f2-2e0b-489b-b0af-787a17103274",
		"readKeyRequestMessage": readKeyRequestMessage,
	}).Debug("Incoming gRPC: 'ReadFromKeyValueStore'")

	// Verify gRPC-proto-message-version in incoming message
	var maxProtoVersionNumner int32 = int32(len(gRPC.CurrentVersionEnum_name) - 1)
	if string(readKeyRequestMessage.GetCurrentVersion()) != gRPC.CurrentVersionEnum_name[maxProtoVersionNumner] {
		// Request from plugin has wrong version, can not execute request
		returnMessage.Comments = "Can execute request due to wrong gRPC-proto-version. KeyValue store is using version: " + gRPC.CurrentVersionEnum_name[maxProtoVersionNumner]
		returnMessage.Acknack = false
	} else {
		// Correct gRPC-proto-version

		// Create the channel that the database object should be sent back on
		returnChannel := make(chan DbResultReadMessageStruct)


		// Convert incoming message into message sent to KeyStorEngine
		DbReadMessage := &DbReadMessageStruct{
			ReadMessage:    *readKeyRequestMessage,
			ResultsChannel: returnChannel,
		}

		// Put readKeyRequestMessage on queue for further processing
		dbReadMessageChannel <- *DbReadMessage
		logger.WithFields(logrus.Fields{
			"ID": "6050780bb-287b-4381-84a4-5b6937758d93",
		}).Debug("'readKeyRequestMessage' was put on the channel to core-DB-engine")

		// Wait for result from core-DB-engine
		databaseReturnMessage := <-returnChannel
		close(returnChannel)

		databaseReturnMessage.

		logger.WithFields(logrus.Fields{
			"ID": "3712717b-ad22-41f5-a05c-b5a3b5794b02",
		}).Debug("Leaving gRPC: 'ReadFromKeyValueStore'")

		// Create message back to parent Gateway/Plugin
		returnMessage.Comments = "'readKeyRequestMessage' was processed by Plugin-KeyValueStore-Engine"
		returnMessage.Acknack = true

	}
	return returnMessage, nil
}
