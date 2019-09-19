package PluginKeyValueDBStore

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	// Common logger for the gateway
	logger *logrus.Logger

	// Database object used for storing any persistent data within Gateway
	db *bolt.DB

	// Database queue used for sending Read questions to database
	dbReadMessageChannel chan DbReadMessageStruct
	// Database queue used for sending Write questions to database
	dbWritewMessageChannel chan gRPC.WriteKeyValueMessage

	// Channel for informationMessage initiated in this gateway
	// removed becasue it resides in object for messages towards Fenix
	//      gatewayTowardsFenixObject.informationMessageChannelTowardsFenix chan *gRPC.InformationMessage
)

// All config parameter will be stored in the following parameter
var gatewayConfig common_code.TomlConfigStruct

// Defines the Read-message sent to Database Engine
type DbReadMessageStruct struct {
	ReadMessage    gRPC.ReadKeyRequestMessage       // Incoming read-message from plugin
	ResultsChannel chan<- DbResultReadMessageStruct // Sending function sends in which channel tp pass back the results on
}

// Defines the Write-message sent to Database Engine
type DbWriteMessageStruct struct {
	WriteMessage   gRPC.WriteKeyValueMessage
	ResultsChannel chan<- DbResultWriteMessageStruct // Sending function sends in which channel tp pass back the results on
}

// Message used for sending back Read-instructions from Database
type DbResultReadMessageStruct struct {
	Err               error  // Error message
	Bucket            string // The Bucket for the message
	ValueSaveTypeId   int    // Defines how the value is saved, As ByteArray, String or as both of them -- 3
	ValueSaveTypeName string // Defines how the value is saved, As ByteArray, String or as both of them
	Value             []byte // The value to be stored as byte array
	ValueString       string // The value to be stored as string
}

// Message used for sending back Write-instructions from Database
type DbResultWriteMessageStruct struct {
	Err    error  // Error message
	Key    string // Key that was Read or Written
	Bucket string // The Bucket for the message
}
