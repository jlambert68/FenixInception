package SendUtrViaMQ

import (
	"encoding/json"
	"errors"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/jlambert68/FenixInception/go_code/common_code/pluginDB"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"strconv"

	//"jlambert/FenixInception2/go_code/common_code"
	gRPC_DB "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
)

// ******************gRPC_DB****************************************************************************************
// Save incoming 'SupportedTestDataDomainsRequest' to Main Database for Fenix Inception
//

/*
func saveSupportedTestDataDomainsMessageInDB(supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (messageSavedInDB bool) {
	var messageToBeSavedInDbAsByteArray []byte
	var err error

	messageSavedInDB = true
	supportedTestDataDomainsRequest.
	// Convert message into byteArray
	messageToBeSavedInDbAsByteArray, err = json.Marshal(*supportedTestDataDomainsRequest)

	// Check if there was an error when converting into byteArray
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID": "08af30bd-0641-44d8-93be-aefd8954d671",

		}).Error("'testInstructionMessage' was saved in Plugin database")
	}


	// Return if the messages was saved or not in database
	return messageSavedInDB

}
*/

// **********************************************************************************************************
// Save incoming 'TestInstruction_RT' to Main Database for Fenix Inception
//
func saveTestInstructionRTMessageInDB(testInstructionRT *gRPC.TestInstruction_RT) (messageSavedInDB bool) {

	var messageTBeSavedInDB gRPC_DB.WriteKeyValueMessage
	var messageToBeSavedInDbAsByteArray []byte
	var err error

	messageSavedInDB = true

	// Convert message into byteArray
	messageToBeSavedInDbAsByteArray, err = json.Marshal(*testInstructionRT)

	// Check if there was an error when converting into byteArray
	if err != nil {
		// Error when converting into byteArray
		logger.WithFields(logrus.Fields{
			"ID": "0a4d828f-c898-4fb1-818e-6e7ef512f527",
		}).Error("'testInstructionRT' could not be converted into an byteArray")

		// Set return code
		messageSavedInDB = false
	} else {
		// OK when converting into byteArray

		// Prepare data message to be sent to KeyValue-DB
		messageTBeSavedInDB = gRPC_DB.WriteKeyValueMessage{
			Key:            testInstructionRT.TestInstructionGuid,
			Bucket:         pluginDB.DbBucketForTestInstructionToBeExecutedMEssages,
			ValueSaveType:  gRPC_DB.ValueSaveTypeEnum_BYTEARRAY,
			Value:          messageToBeSavedInDbAsByteArray,
			ValueString:    "",
			CurrentVersion: gRPC_DB.CurrentVersionEnum(gRPC_DB.CurrentVersionEnum_value[pluginDB.GetHighestPluginGRPCVersion()]),
		}

		// Send message, to be sent to KeyValue-DB, via channel
		messageWasSavedInDB := SaveMessageUsingGRPC(messageTBeSavedInDB)

		// Check if message was saved in KeyValue-DB
		if messageWasSavedInDB == false {
			logger.WithFields(logrus.Fields{
				"ID": "bd589b5c-95da-4884-8bfa-9fe487d57c1f",
			}).Error("Message couldn't be saved in KeyValue-DB")

			// Set return code
			messageSavedInDB = false
		}
	}

	// Return that the messages was saved or not in database
	return messageSavedInDB
}

// **********************************************************************************************************
// Call KeyValue-DB using gRPC to save message
//
func SaveMessageUsingGRPC(messageTBeSavedInDB gRPC_DB.WriteKeyValueMessage) (messageSavedInDB bool) {

	var (
		addressToDial                  string
		remotePluginDbServerConnection *grpc.ClientConn
		gRpcClientTowardPluginDB       gRPC_DB.PluginDBClient
		//gRPCresponse                   gRPC_DB.AckNackResponse

		gRpcContexType context.Context
	)
	//var gRPCerr error
	var err error

	messageSavedInDB = true

	// Use background call in gRPC-call
	gRpcContexType = context.Background()

	// Get address to next child gateway/plugin
	addressToDial, err = getPluginDBAddressAndPort()

	if err != nil {
		// Couldn't get Address to PluginDB
		messageSavedInDB = false

		logger.WithFields(logrus.Fields{
			"ID":            "6b3d5eb4-8bcf-43ce-a3d7-0b6fb8578ea3",
			"addressToDial": addressToDial,
			"err":           err,
		}).Error("Couldn't get gRPC-address to KeyValue-DB")
	} else {
		// Got gRPC-address to PluginDB

		// Set up connection to Plugin
		// logger.WithFields(logrus.Fields{
		//"ID": "6b3d5eb4-8bcf-43ce-a3d7-0b6fb8578ea3",
		//	"addressToDial": addressToDial,

		remotePluginDbServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())

		// Check if the dial was a success
		if err != nil {
			// Could not dial PluginDB

			messageSavedInDB = false

			logger.WithFields(logrus.Fields{
				"ID":  "1cb7cdf1-c17e-4b62-9f2d-611d4e865f79",
				"err": err,
			}).Error("Error when dialing KeyValue-DB")

		} else {
			// Dial was successful towards PluginDB

			// Creates a new gateway gRPC-Client towards Plugin
			gRpcClientTowardPluginDB = gRPC_DB.NewPluginDBClient(remotePluginDbServerConnection)

			// Send message to KeyValue-DB using gRPC-call
			gRPCresponse, gRPCerr := gRpcClientTowardPluginDB.WriteToKeyValueStore(gRpcContexType, &messageTBeSavedInDB)

			// Chek of there was an error in the gRPC-call
			if gRPCerr != nil {
				messageSavedInDB = false

				logger.WithFields(logrus.Fields{
					"ID":           "9fc03ecb-a25c-4fe7-9d8f-029521883125",
					"gRPCresponse": gRPCresponse,
					"gRPCerr":      gRPCerr,
				}).Error("Error when doing gRPC-call to KeyValue-DB")
			} else {
				// SUccess in doing the gRPC-call

				// Check if the the message was saved
				if gRPCresponse.Acknack == false {
					// Message was not saved
					messageSavedInDB = false
					logger.WithFields(logrus.Fields{
						"ID":           "549a16be-e85a-4456-bc67-c3e14852275e",
						"gRPCresponse": gRPCresponse,
					}).Error("THe message was not saved in the KeyValue-DB")
				} else {
					// Message was saved in KeyValue-DB
					logger.WithFields(logrus.Fields{
						"ID":           "1012a125-7dbb-43cb-86a1-4ab2b36d1396",
						"gRPCresponse": gRPCresponse,
					}).Debug("Message was saved in KeyValue-DB")
				}
			}
		}
	}
	return messageSavedInDB
}

// **********************************************************************************************************
// Get address and port for KeyValue-DB
//
func getPluginDBAddressAndPort() (addressAndPort string, err error) {

	err = nil
	addressAndPort = ""

	// Check if KeyValueStoreIpAddress is a correct formated ip addess
	if net.ParseIP(gatewayConfig.KeyValueStoreIdentification.KeyValueStoreIpAddress) != nil {
		err = errors.New("KeyValueStore IpAddress is not correc formated: '" + gatewayConfig.KeyValueStoreIdentification.KeyValueStoreIpAddress + "'")
	} else {
		// Correct formated
		// Check if port is bigger than 999, just th be sure that the port was correct typed
		if gatewayConfig.InitialClientPort.InitialClientPort < 1000 {
			err = errors.New("KeyValueStore port is less then 1000: '" + strconv.FormatInt(int64(gatewayConfig.InitialClientPort.InitialClientPort), 10) + "'")
		}

		addressAndPort = gatewayConfig.KeyValueStoreIdentification.KeyValueStoreIpAddress + ":" + strconv.FormatInt(int64(gatewayConfig.InitialClientPort.InitialClientPort), 10)
	}

	return addressAndPort, err
}
