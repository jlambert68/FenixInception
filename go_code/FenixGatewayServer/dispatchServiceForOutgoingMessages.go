package FenixGatewayServer

import (
	"encoding/json"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"jlambert/FenixInception2/go_code/TestExecutionGateway"
	"jlambert/FenixInception2/go_code/common_code"
	"time"
)

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func sendTestInstructionTowardsPlugin(testInstructionId string) (err error) {

	err = nil

	// Get Payload for TestInstruction to be sent to Plugin
	testInstructionPayload, err := getTestInstructionPayloadToPlugin(testInstructionId)

	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		TestExecutionGateway.PutTestInstructionMessageOnChannelTowardsPlugin(testInstructionPayload)
		//err = dispatchEngineForTestInstructions(testInstructionPayload)

	}

	return err

}

// ********************************************************************************************
// Send TestInstructions to Fenix own Gateway for then to be forwarded towards Plugin
//

func dispatchEngineForTestInstructions(testInstructionToBeForwarded *gRPC.TestInstruction_RT) (err error) {

	// Service can be started and Stopped by central control of the Gateway
	if gatewayMustStopProcessing == true {
		// Service should be stopped from processing any messages
		time.Sleep(common_code.ServieSleepTime * time.Second)

	} else {
		// Process messages

		logger.WithFields(logrus.Fields{
			"ID":                           "ffd854c0-f570-4d1a-b945-f479df8d3be1",
			"testInstructionToBeForwarded": testInstructionToBeForwarded,
		}).Debug("Received a new TestInstruction that shoud be forwarded")

		// Send TestInstruction to Fenix Gateway using gRPC-call
		addressToDial := getClientAddressAndPort(testInstructionToBeForwarded.PluginId)

		// Set up connection to Fenix Gateway
		remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "73a85fbe-bf9d-4b15-bdf1-d56aca0f7754",
				"addressToDial": addressToDial,
				"error message": err,
			}).Error("Did not connect to Fenix gateway Server!")

		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "0110fa80-c543-4d5c-a02d-ea23fd8969f6",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Fenix Gateway")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			testInstructionToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
			testInstructionToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

			// Do gRPC-call to Fenix Gateway
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendTestInstructionTowardsPlugin(ctx, testInstructionToBeForwarded)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "69c6d6e9-f15e-4e73-a176-8ba947bf9f92",
					"returnMessage": returnMessage,
					"error":         err,
				}).Error("Problem to send TestInstruction to Fenix Gateway")

			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "466ff0c7-eb34-4188-a1bc-d51bd3cda850",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send TestInstruction to Fenix Gateway")

			}
		}
	}

	return err
}

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func sendGetSupportedTestDataDomainsTowardsPlugin(testInstructionId string) (err error) {
	// Get Payload for TestInstruction to be sent to Plugin
	getSupportedTestDataDomainsMessage, err := getTestInstructionPayloadToPlugin(testInstructionId)

	if err != nil {
		// Error when loading Payload
		return err
	} else {
		// Send TestInstruction to Next Gateway
		err = dispatchEngineForTestInstructions(getSupportedTestDataDomainsMessage)

	}

	return err

}

// ********************************************************************************************
// Send getSupportedTestDomains-message to Fenix own Gateway for then to be forwarded towards Plugin
//

func dispatchEngineForSupportedTestDomains(getSupportedTestDomainsMessage *gRPC.TestInstruction_RT) (err error) {

	var messageTypeToSendText string = "getSupportedTestDomainsMessage"

	// Service can be started and Stopped by central control of the Gateway
	if gatewayMustStopProcessing == true {
		// Service should be stopped from processing any messages
		time.Sleep(common_code.ServieSleepTime * time.Second)

	} else {
		// Process messages

		logger.WithFields(logrus.Fields{
			"ID":                  "ffd854c0-f570-4d1a-b945-f479df8d3be1",
			messageTypeToSendText: getSupportedTestDomainsMessage,
		}).Debug("Received a new " + messageTypeToSendText + " that shoud be forwarded")

		// Send Message to Fenix Gateway using gRPC-call
		addressToDial := getClientAddressAndPort(getSupportedTestDomainsMessage.PluginId)

		// Set up connection to Fenix Gateway
		remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "73a85fbe-bf9d-4b15-bdf1-d56aca0f7754",
				"addressToDial": addressToDial,
				"error message": err,
			}).Error("Did not connect to Fenix gateway Server!")

		} else {
			logger.WithFields(logrus.Fields{
				"ID":            "0110fa80-c543-4d5c-a02d-ea23fd8969f6",
				"addressToDial": addressToDial,
			}).Debug("gRPC connection OK to Fenix Gateway")

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			getSupportedTestDomainsMessage.SenderId = gatewayConfig.GatewayIdentification.GatewayId
			getSupportedTestDomainsMessage.SenderName = gatewayConfig.GatewayIdentification.GatewayName

			// Do gRPC-call to Fenix Gateway
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendTestInstructionTowardsPlugin(ctx, getSupportedTestDomainsMessage)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "69c6d6e9-f15e-4e73-a176-8ba947bf9f92",
					"returnMessage": returnMessage,
					"error":         err,
				}).Error("Problem to send " + messageTypeToSendText + " to Fenix Gateway")

			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "466ff0c7-eb34-4188-a1bc-d51bd3cda850",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send " + messageTypeToSendText + " to Fenix Gateway")

			}
		}
	}

	return err
}

// ********************************************************************************************
// Forward getSupportedTestDomains-messages from incoming channel to child gateway or plugin
//

func (gatewayObject *gatewayTowardsPluginObjectStruct) dispatchEngineForSupportedTestDomains() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// ***** Wait for data comes from channel to dispatch engine *****
			getSupportedTestDomainsToBeForwarded := <-gatewayObject.supportedTestDataDomainsRequestChannelTowardsPlugin

			// Check number of messages in channel
			channelSinaling(len(gatewayObject.testInstructionMessageChannelTowardsPlugin),
				"supportedTestDataDomainsRequestChannelTowardsPlugin",
				"3c85d4c4-260d-4612-a153-2eeeff864621")

			logger.WithFields(logrus.Fields{
				"ID":                                   "05f9ae3c-dc0f-4107-ae85-ae57d6dde6c1",
				"testExecutionLogMessageToBeForwarded": getSupportedTestDomainsToBeForwarded,
			}).Debug("Received a new 'getSupportedTestDomainsToBeForwarded' from channel that shoud be forwarded")

			// ***** Send getSupportedTestDomains-message to client gateway/plugin using gRPC-call *****
			addressToDial := getClientAddressAndPort(getSupportedTestDomainsToBeForwarded.PluginId)

			// Set up connection to Parent Gateway or Fenix
			remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				// Connection Not OK
				LogErrorAndSendInfoToFenix(
					"607ff7f6-4452-49f6-a9ec-608118073a3b",
					gRPC.InformationMessage_WARNING,
					"addressToDial",
					addressToDial,
					err.Error(),
					"Did not connect to child (Gateway or Plugin) Server!",
				)

				// Convert getSupportedTestDomainsToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*getSupportedTestDomainsToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"f1c5f3f6-341e-4fb3-96aa-9886da021c8d",
						gRPC.InformationMessage_FATAL,
						"getSupportedTestDomainsToBeForwarded",
						getSupportedTestDomainsToBeForwarded.String(),
						err.Error(),
						"Error when converting 'getSupportedTestDomainsToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Save message to local DB for later processing
					_ = SaveMessageToLocalDB(
						getSupportedTestDomainsToBeForwarded.MessageId,
						testExecutionLogMessageToBeForwardedByteArray,
						BucketForResendOfGetTestdataDomainsToPlugin,
						"1abe7614-c951-4186-821f-deb2b4203f64",
					)
				}

			} else {
				//Connection OK
				logger.WithFields(logrus.Fields{
					"ID":            "954306f9-d890-48ad-ac5a-68143166b48d",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to child-gateway/Plugin!")

				// Convert getSupportedTestDomainsToBeForwarded-struct into a byte array
				testExecutionLogMessageToBeForwardedByteArray, err := json.Marshal(*getSupportedTestDomainsToBeForwarded)

				if err != nil {
					// Error when Unmarshaling to []byte
					LogErrorAndSendInfoToFenix(
						"12c16a5f-ba94-4e16-9607-214b9281da9a",
						gRPC.InformationMessage_FATAL,
						"getSupportedTestDomainsToBeForwarded",
						getSupportedTestDomainsToBeForwarded.String(),
						err.Error(),
						"Error when converting 'getSupportedTestDomainsToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
					)

				} else {
					// Marshaling to []byte OK

					// Creates a new gateway Client
					gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					getSupportedTestDomainsToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					getSupportedTestDomainsToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Do gRPC-call to child gateway or Plugin
					ctx := context.Background()
					returnMessage, err := gatewayClient.GetSupportedTestDataDomains(ctx, getSupportedTestDomainsToBeForwarded)
					if err != nil {
						// Error when rending gRPC to child
						LogErrorAndSendInfoToFenix(
							"6b77117b-e98c-4451-a91a-d3889438b3d7",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessage.String(),
							err.Error(),
							"Problem to send 'getSupportedTestDomainsToBeForwarded' to child-Gateway or Plugin",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							getSupportedTestDomainsToBeForwarded.MessageId,
							testExecutionLogMessageToBeForwardedByteArray,
							BucketForResendOfGetTestdataDomainsToPlugin,
							"9446b644-66f3-452a-9436-add798b0dad9",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "77eef2f1-116d-4342-a02e-dce88d1e1b52",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK of 'getSupportedTestDomainsToBeForwarded' to child-Gateway or Plugin")

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}
