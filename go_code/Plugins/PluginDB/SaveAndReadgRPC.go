package PluginKeyValueDBStore

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from Plugin that need to store Values in KeyValueStore
//
func (gRPCServerForDbKeyValueStore *gRPCServerForDbKeyValueStoreStruct) WriteToKeyValueStore(ctx context.Context, writeKeyValueMessage *gRPC.WriteKeyValueMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                   "c8ad7a99-a27f-4f77-be96-81733d57d36b",
		"writeKeyValueMessage": writeKeyValueMessage,
	}).Debug("Incoming gRPC: 'WriteToKeyValueStore'")

	// Verify gRPC-proto-message-version in incoming message
	var maxProtoVersionNumner int32 = int32(len(gRPC.CurrentVersionEnum_name) - 1)
	if string(writeKeyValueMessage.GetCurrentVersion()) != gRPC.CurrentVersionEnum_name[maxProtoVersionNumner] {
		// Request from plugin has wrong version, can not execute request
		returnMessage.Comments = "Can't execute request due to wrong gRPC-proto-version. KeyValue store is using version: " + gRPC.CurrentVersionEnum_name[maxProtoVersionNumner]
		returnMessage.Acknack = false
	} else {
		// Correct gRPC-proto-version

		// Create the channel that the database object should be sent back on
		returnChannel := make(chan DbResultWriteMessageStruct)

		// Convert incoming message into message sent to KeyStorEngine
		DbWriteMessage := &DbWriteMessageStruct{
			WriteMessage:   *writeKeyValueMessage,
			ResultsChannel: returnChannel,
		}

		// Put WriteToKeyValueStore on queue for further processing
		dbWritewMessageChannel <- *DbWriteMessage
		logger.WithFields(logrus.Fields{
			"ID": "ec530b7e-a4e1-4edc-b8a5-4cc487dd0f33",
		}).Debug("'WriteToKeyValueStore' was put on the channel to core-DB-engine")

		// Wait for result from core-DB-engine
		databaseReturnMessage := <-returnChannel
		close(returnChannel)

		if databaseReturnMessage.Err == nil {
			// OK respons then convert result into returnMessage
			returnMessage = &gRPC.AckNackResponse{
				Acknack:  true,
				Comments: "Value saved in KeyValue-store",
			}
		} else {
			// Error in Response, put that into returnMessage
			returnMessage = &gRPC.AckNackResponse{
				Acknack:  false,
				Comments: databaseReturnMessage.Err.Error(),
			}
		}
		logger.WithFields(logrus.Fields{
			"ID": "3bebc000-70e6-4937-9fae-fd5ae1f9177f",
		}).Debug("Leaving gRPC: 'WriteToKeyValueStore'")
	}
	return returnMessage, nil
}

// ********************************************************************************************
// Call from Plugin that need to read Values in KeyValueStore
//
func (gRPCServerForDbKeyValueStore *gRPCServerForDbKeyValueStoreStruct) ReadFromKeyValueStore(ctx context.Context, readKeyRequestMessage *gRPC.ReadKeyRequestMessage) (*gRPC.ValueResponseMessage, error) {

	var returnMessage *gRPC.ValueResponseMessage

	logger.WithFields(logrus.Fields{
		"ID":                    "3628b2f2-2e0b-489b-b0af-787a17103274",
		"readKeyRequestMessage": readKeyRequestMessage,
	}).Debug("Incoming gRPC: 'ReadFromKeyValueStore'")

	// Verify gRPC-proto-message-version in incoming message
	var maxProtoVersionNumner int32 = int32(len(gRPC.CurrentVersionEnum_name) - 1)
	if string(readKeyRequestMessage.GetCurrentVersion()) != gRPC.CurrentVersionEnum_name[maxProtoVersionNumner] {
		// Request from plugin has wrong version, can not execute request
		returnMessage.Comments = "Can't execute request due to wrong gRPC-proto-version. KeyValue store is using version: " + gRPC.CurrentVersionEnum_name[maxProtoVersionNumner]
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

		if databaseReturnMessage.Err == nil {
			// OK respons then convert result into returnMessage
			returnMessage = &gRPC.ValueResponseMessage{
				Key:             databaseReturnMessage.Key,
				Bucket:          databaseReturnMessage.Bucket,
				ValueSaveType:   gRPC.ValueSaveTypeEnum(databaseReturnMessage.ValueSaveTypeId),
				Value:           databaseReturnMessage.Value,
				ValueString:     databaseReturnMessage.ValueString,
				UpdatedDateTime: databaseReturnMessage.UpdatedDateTime,
				Acknack:         true,
				Comments:        "Value found in KeyValue-store",
			}
		} else {
			// Error in Response, put that into returnMessage
			returnMessage = &gRPC.ValueResponseMessage{
				Key:             "",
				Bucket:          "",
				ValueSaveType:   0,
				Value:           nil,
				ValueString:     "",
				UpdatedDateTime: "",
				Acknack:         false,
				Comments:        databaseReturnMessage.Err.Error(),
			}
		}
		logger.WithFields(logrus.Fields{
			"ID": "3712717b-ad22-41f5-a05c-b5a3b5794b02",
		}).Debug("Leaving gRPC: 'ReadFromKeyValueStore'")

	}
	return returnMessage, nil
}
