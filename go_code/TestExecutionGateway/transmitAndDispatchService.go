package TestExecutionGateway

import (
	"encoding/json"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

// ********************************************************************************************
// Initiate Transmit Engines for messages towards Fenix and Dispatch Engines for messages towards Plugins
//

func initiateAllTransmitAndDispatchEngines() {

	// *** Towards Plugin ***
	// Start a Dispatch Engine, for 'TestInstructionMessageTowardsPlugin' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionMessageTowardsPlugin, common_code.DispatchEngineTowardsPlugin, common_code.GatewayEngine, nil)

	// Start a Dispatch Engine, for 'SupportedTestDataDomainsRequestMessageTowardsPlugin,' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin, common_code.DispatchEngineTowardsPlugin, common_code.GatewayEngine, nil)

	// *** Towards Fenix ***
	// Start a Transmit Engine, for 'informationMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeInformationMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'timeOutMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'spportedTestDataDomainsMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeTestExecutionLogMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'availbleTestInstructionAtPluginMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'testInstructionExecutionResultMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

	// Start a Transmit Engine, for 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' as a go-routine
	go TransmitAndDispatchEngine(common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix, common_code.TransmitEngineTowardsFenix, common_code.GatewayEngine, nil)

}

// ********************************************************************************************
// Forward X-message from channel towards Fenix
//

func TransmitAndDispatchEngine(channelType string, transmitOrDispatchEngineType string, fenixOrGatewayType string, fenixFunctionToCall func()) {

	// Messages towards Plugin
	var (
		testInstructionMessageToBeForwardedTowardsPlugin                 *gRPC.TestInstruction_RT
		supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin *gRPC.SupportedTestDataDomainsRequest
	)

	// Messages towards Fenix
	var (
		informationMessageToBeForwardedTowardsFenix                         *gRPC.InformationMessage
		timeOutMessageToBeForwardedTowardsFenix                             *gRPC.TestInstructionTimeOutMessage
		testExecutionLogMessageToBeForwardedTowardsFenix                    *gRPC.TestExecutionLogMessage
		spportedTestDataDomainsMessageToBeForwardedTowardsFenix             *gRPC.SupportedTestDataDomainsMessage
		availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix     *gRPC.AvailbleTestInstructionAtPluginMessage
		availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix      *gRPC.AvailbleTestContainersAtPluginMessage
		testInstructionExecutionResultMessageToBeForwardedTowardsFenix      *gRPC.TestInstructionExecutionResultMessage
		supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix *gRPC.SupportedTestDataDomainsWithHeadersMessage
	)

	/*
		//  informationMessage towards Fenix
		var (
			InformationMessageChannelTowardsFenix chan *gRPC.InformationMessage

			// testInstructionTimeOutMessage towards Fenix
			TestInstructionTimeOutMessageChannelTowardsFenix chan *gRPC.TestInstructionTimeOutMessage

			// testExecutionLogMessage towards Fenix
			TestExecutionLogMessageChannelTowardsFenix chan *gRPC.TestExecutionLogMessage

			// supportedTestDataDomainsMessage towards Fenix
			SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsMessage

			// availbleTestInstruction<AtPluginMessage towards Fenix
			AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestInstructionAtPluginMessage

			//availbleTestContainersAtPluginMessage towars Fenix
			AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix chan *gRPC.AvailbleTestContainersAtPluginMessage

			//availbleTestContainersAtPluginMessage towars Fenix
			TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix chan *gRPC.TestInstructionExecutionResultMessage

			//supportedTestDataDomainsWithHeadersMessage towars Fenix
			SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix chan *gRPC.SupportedTestDataDomainsWithHeadersMessage
		)

		// Internal queues used by the gateway
		var (
			// TestInstruction Towards Plugin
			TestInstructionMessageChannelTowardsPlugin chan *gRPC.TestInstruction_RT

			// supportedTestDataDomainsRequest Towards Plugin
			SupportedTestDataDomainsRequestChannelTowardsPlugin chan *gRPC.SupportedTestDataDomainsRequest
		)
	*/

	var err error
	var messageId string
	var bucket string
	var messageToBeForwardedByteArray []byte
	var id string
	var infoHeader string
	var info string
	var errorMessage string
	var message string
	var NumberOfMessagesInChannel int
	var pluginId string

	var (
		addressToDial          string
		remoteServerConnection *grpc.ClientConn
		gRpcClientTowardsFenix gRPC.GatewayTowardsFenixClient
		gRpcClientTowardPlugin gRPC.GatewayTowayPluginClient

		gRpcContexType context.Context
		dialSuccess    = false
	)

	var gRPCerr error
	var returnMessageAckNackResponse *gRPC.AckNackResponse
	var returnMessageString string

	// Use background call in gRPC-call
	gRpcContexType = context.Background()

	// Check it is possible to dial parent gateway/Fenix or child gatway/plugin
	addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection = dialChildOrParenGateway(addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection, transmitOrDispatchEngineType, pluginId)

	// Wait for message on channel
	// Do gRPC-call to client gateway or Fenix

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true || dialSuccess == false {
			// Service should be stopped from processing any messages
			time.Sleep(common_code.ServieSleepTime * time.Second)

			// Check it is possible to dial parent gateway/Fenix or child gatway/plugin
			addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection = dialChildOrParenGateway(addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection, transmitOrDispatchEngineType, pluginId)

		} else {

			// Check for ONE message, in DB, to Resend (If so then put IT on channel)
			objectsWherPutOnChannel := checkForSavedTemporaryObjectsInDbThatWillBePutOnChannel(true, channelType, transmitOrDispatchEngineType, &gatewayChannelPackage)
			if objectsWherPutOnChannel == true {
				logger.WithFields(logrus.Fields{
					"ID":                      "84821702-66f1-48bc-a8cd-eee8951fb880",
					"objectsWherPutOnChannel": objectsWherPutOnChannel,
				}).Debug("DB checked for one saved objects and there was one object saved in DB that were put on channel again")
			} else {
				logger.WithFields(logrus.Fields{
					"ID":                      "74168c9d-42b2-4050-a9da-c5ff58910302",
					"objectsWherPutOnChannel": objectsWherPutOnChannel,
				}).Debug("DB checked for saved one object and there were NO objects saved in DB")
			}

			// Run service and process messages

			// Wait for message to come, from defined channel, to transmit engine ****
			// Select on messageType and convert message into byteArray
			// Set variables depending on messageType

			// Decide if it's TransmitEngine or a DispatchEngine
			switch transmitOrDispatchEngineType {

			// Choose DispatchEngine
			case common_code.DispatchEngineTowardsPlugin:
				//  Decide the correct DispatchEngine to use
				switch channelType {
				case common_code.ChannelTypeTestInstructionMessageTowardsPlugin:
					testInstructionMessageToBeForwardedTowardsPlugin = <-gatewayChannelPackage.TestInstructionMessageChannelTowardsPlugin
					messageToBeForwardedByteArray, err = json.Marshal(*testInstructionMessageToBeForwardedTowardsPlugin)
					// If error when marshaling then use the following information
					id = "22318eb6-5d57-4183-8238-0bc40e16ff7f"
					infoHeader = "testInstructionMessageToBeForwardedTowardsPlugin"
					info = ""
					errorMessage = ""
					message = "testInstructionMessageToBeForwardedTowardsPlugin"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					testInstructionMessageToBeForwardedTowardsPlugin.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					testInstructionMessageToBeForwardedTowardsPlugin.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = testInstructionMessageToBeForwardedTowardsPlugin.MessageId
					bucket = common_code.BucketForResendTestInstructionTowardsPlugin
					NumberOfMessagesInChannel = len(gatewayChannelPackage.TestInstructionMessageChannelTowardsPlugin)

				case common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
					supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin = <-gatewayChannelPackage.SupportedTestDataDomainsRequestChannelTowardsPlugin
					messageToBeForwardedByteArray, err = json.Marshal(*testInstructionMessageToBeForwardedTowardsPlugin)
					// If error when marshaling then use the following information
					id = "5f0cc657-97fb-49d0-9c10-e63f6a527676"
					infoHeader = "supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin"
					info = ""
					errorMessage = ""
					message = "supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin.MessageId
					bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin
					NumberOfMessagesInChannel = len(gatewayChannelPackage.SupportedTestDataDomainsRequestChannelTowardsPlugin)

				default:
					LogErrorAndSendInfoToFenix(
						"f5cbd96e-e6a3-4830-abcc-1dccdd1394dc",
						gRPC.InformationMessage_FATAL,
						"No know 'channelType'",
						channelType,
						err.Error(),
						"No know 'channelType' in DispatchEngine",
					)
				}

				// Choose TransmitEngine
			case common_code.TransmitEngineTowardsFenix:
				//  Deside the correct TransmitEngine to use
				switch channelType {
				case common_code.ChannelTypeInformationMessageTowardsFenix:
					informationMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.InformationMessageChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*informationMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "6798902d-36b9-4598-a125-4be3abd4a23b"
					infoHeader = "informationMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "informationMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					informationMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					informationMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = informationMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfInfoMessagesTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.InformationMessageChannelTowardsFenix)

				case common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix:
					timeOutMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.TestInstructionTimeOutMessageChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*timeOutMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "e80978f1-2f2c-43f7-bd4b-218e852f2f72"
					infoHeader = "timeOutMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "timeOutMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					timeOutMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					timeOutMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = timeOutMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOTimeOutMesagesTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.TestInstructionTimeOutMessageChannelTowardsFenix)

				case common_code.ChannelTypeTestExecutionLogMessageTowardsFenix:
					testExecutionLogMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.TestExecutionLogMessageChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*testExecutionLogMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "3f8b01b6-da9b-4cee-b787-6a0e9efe9fed"
					infoHeader = "testExecutionLogMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "testExecutionLogMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					testExecutionLogMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					testExecutionLogMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = testExecutionLogMessageToBeForwardedTowardsFenix.LogMessageId
					bucket = common_code.BucketForResendOfLogMesagesTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.TestExecutionLogMessageChannelTowardsFenix)

				case common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix:
					spportedTestDataDomainsMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*spportedTestDataDomainsMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "d3a33be3-6cf8-459d-bda2-ff64c53d249a"
					infoHeader = "spportedTestDataDomainsMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "spportedTestDataDomainsMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					spportedTestDataDomainsMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					spportedTestDataDomainsMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = spportedTestDataDomainsMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfSupportedTestDataDomainsTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix)

				case common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
					availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "6524a739-e628-4265-a6cb-d64597a20c3e"
					infoHeader = "availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfAvailableTestInstructionsTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix)

				case common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
					availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "14b2aa63-7e69-4e80-9d6b-796c286997c9"
					infoHeader = "availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix"

					messageId = availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfAvailableTestContainersTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix)

				case common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix:
					testInstructionExecutionResultMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*testInstructionExecutionResultMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "72d128f7-e84c-408a-87a9-e0c8b84c6841"
					infoHeader = "testInstructionExecutionResultMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "testInstructionExecutionResultMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					testInstructionExecutionResultMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					testInstructionExecutionResultMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = testInstructionExecutionResultMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfTestInstructionExecutionResultTowardsFenix
					NumberOfMessagesInChannel = len(gatewayChannelPackage.TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix)

				case common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
					supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix = <-gatewayChannelPackage.SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "89df74d0-8ba7-47df-bddc-5864a9b376f1"
					infoHeader = "supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix"

					// ChangeSenderId to this gatway's SenderId before sending the data forward
					supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix.SenderId = gatewayConfig.GatewayIdentification.GatewayId
					supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix.SenderName = gatewayConfig.GatewayIdentification.GatewayName

					// Information used when saving message to local database
					messageId = supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix.MessageId
					bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin
					NumberOfMessagesInChannel = len(gatewayChannelPackage.SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix)

				default:
					LogErrorAndSendInfoToFenix(
						"e22077e8-e2ee-4cb2-aabe-d3603a7bfab4",
						gRPC.InformationMessage_FATAL,
						"No know 'channelType'",
						channelType,
						err.Error(),
						"No know 'channelType' in TransmitEngine",
					)
				}

			default:
				// No known 'transmitOrDispatchEngineType'
				LogErrorAndSendInfoToFenix(
					"79ccdec0-803e-4133-99cb-f0da3a55de7c",
					gRPC.InformationMessage_FATAL,
					"No know 'transmitOrDispatchEngineType'",
					transmitOrDispatchEngineType,
					"No known 'transmitOrDispatchEngineType when dialing to child/parent",
					"No known 'transmitOrDispatchEngineType when dialing to child/parent",
				)
			}

			// Handle error from marshaling []byte
			if err != nil {
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					infoHeader,
					info,
					errorMessage,
					"Error when converting '"+message+"' into a byte array, stopping futher processing of this TestInstruction.",
					//"Error when converting 'testInstructionExecutionResultMessageToBeForwardedTowardsFenix' into a byte array, stopping futher processing of this TestInstruction",
				)

			} else {
				// Converting into byteArray succeeded
				logger.WithFields(logrus.Fields{
					"ID":                   "1eb79d7e-9a44-4315-bbd4-9e008b0c55fc",
					"messageTOBeForwarded": infoHeader,
				}).Debug("Received a new message to be forwarded")

				// Check number of messages in channel
				channelSinaling(NumberOfMessagesInChannel,
					infoHeader,
					"db1bd71e-4ad4-40df-aba5-e80f69899e08")

				//TODO kopier **switch transmitOrDispatchEngineType {** to första Switchen

				// Decide if "we" are Fenix, then save message to SQL-DB otherwise send message forward to parent/child
				switch fenixOrGatewayType {

				case common_code.FenixEngine:
					// Save incoming mesage to SQL-DB and trigger Fenix by calling the function that was sent when initializing
					fenixFunctionToCall()

				case common_code.GatewayEngine:
					// Do the gRPC-call to parent gateway/Fenix or child gatway/plugin
					// Decide if it's TransmitEngine or a DispatchEngine
					switch transmitOrDispatchEngineType {

					// Choose DispatchEngine
					case common_code.DispatchEngineTowardsPlugin:
						//  Decide the correct DispatchEngine to use
						switch channelType {
						case common_code.ChannelTypeTestInstructionMessageTowardsPlugin:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardPlugin.SendTestInstructionTowardsPlugin(gRpcContexType, testInstructionMessageToBeForwardedTowardsPlugin)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardPlugin.GetSupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin)
							returnMessageString = returnMessageAckNackResponse.String()

						default:
							LogErrorAndSendInfoToFenix(
								"03a585ed-3bc9-46a8-b1b9-55426c4ecc84",
								gRPC.InformationMessage_FATAL,
								"No know 'channelType'",
								channelType,
								err.Error(),
								"No know 'channelType' in DispatchEngine",
							)
						}

						// Choose TransmitEngine
					case common_code.TransmitEngineTowardsFenix:
						//  Deside the correct TransmitEngine to use
						switch channelType {
						case common_code.ChannelTypeInformationMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.SendTestInstructionTimeOutTowardsFenix(gRpcContexType, timeOutMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeTestExecutionLogMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.SendTestExecutionLogTowardsFenix(gRpcContexType, testExecutionLogMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.RegistrateAvailableTestDataDomains(gRpcContexType, spportedTestDataDomainsMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.RegisterAvailbleTestInstructions(gRpcContexType, availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.RegistrateAailableTestContainers(gRpcContexType, availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.SendTestInstructionResultTowardsFenix(gRpcContexType, testInstructionExecutionResultMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						case common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
							returnMessageAckNackResponse, gRPCerr = gRpcClientTowardsFenix.SupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix)
							returnMessageString = returnMessageAckNackResponse.String()

						default:
							LogErrorAndSendInfoToFenix(
								"d38c0522-77a6-4f90-b98e-02a4d37f3114",
								gRPC.InformationMessage_FATAL,
								"No know 'channelType'",
								channelType,
								err.Error(),
								"No know 'channelType' in TransmitEngine",
							)
						}

					default:
						// No known 'transmitOrDispatchEngineType'
						LogErrorAndSendInfoToFenix(
							"2c976770-ef52-4bd1-8a52-aaf1f64ebc69",
							gRPC.InformationMessage_FATAL,
							"No know 'transmitOrDispatchEngineType'",
							transmitOrDispatchEngineType,
							"No known 'transmitOrDispatchEngineType when dialing to child/parent",
							"No known 'transmitOrDispatchEngineType when dialing to child/parent",
						)
					}

				default:
					LogErrorAndSendInfoToFenix(
						"721a5122-192d-492f-9882-ab40157b651d",
						gRPC.InformationMessage_FATAL,
						"No know 'fenixOrGatewayType'",
						fenixOrGatewayType,
						err.Error(),
						"No know 'fenixOrGatewayType' in TransmitEngine/Fenix",
					)
				}

				switch fenixOrGatewayType {
				// Are "we" at Fenix or at a gateway
				case common_code.FenixEngine:
					// At Fenix
					// Do nothing

				case common_code.GatewayEngine:
					// At Gateway

					if gRPCerr != nil {
						// Error when sending gRPC to parent gatewayFenix or child gatway/plugin
						dialSuccess = false

						LogErrorAndSendInfoToFenix(
							"e642a035-fd95-45aa-b625-4dcedd65c3a5",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessageString,
							err.Error(),
							"Problem to send message from '"+channelType+"' to parent Gateway/Fenix or or child Gatway/Plugin",
						)

						// Save message to local DB for later processing
						_ = SaveMessageToLocalDB(
							messageId,
							messageToBeForwardedByteArray,
							bucket,
							"05d6219e-0ab8-40e8-b671-0bea0ea24978",
						)

					} else {
						// gRPC Send message OK
						logger.WithFields(logrus.Fields{
							"ID":            "9c3204ea-7f1b-4a34-ae5e-e11d45b1c17f",
							"addressToDial": addressToDial,
						}).Debug("gRPC-send OK for '" + infoHeader + "' to Parent-Gateway or Fenix")

						// Check for All messages, in DB, to Resend (If so then put THEM ALL on channel)
						objectsWherPutOnChannel := checkForSavedTemporaryObjectsInDbThatWillBePutOnChannel(false, channelType, transmitOrDispatchEngineType, &gatewayChannelPackage)
						if objectsWherPutOnChannel == true {
							logger.WithFields(logrus.Fields{
								"ID":                      "b9dbbd63-b6ac-45de-b92e-09d3b049822c",
								"objectsWherPutOnChannel": objectsWherPutOnChannel,
							}).Debug("DB checked for all saved objects and there were objects saved in DB that were put on channel again")
						} else {
							logger.WithFields(logrus.Fields{
								"ID":                      "dd93ce5f-e5f2-4b91-b54d-f720a7b6a437",
								"objectsWherPutOnChannel": objectsWherPutOnChannel,
							}).Debug("DB checked for saved ALL objects and there were NO objects saved in DB")
						}
					}

				default:
					LogErrorAndSendInfoToFenix(
						"af30d945-42c1-4fdf-868f-d8a7d857399b",
						gRPC.InformationMessage_FATAL,
						"No know 'fenixOrGatewayType'",
						fenixOrGatewayType,
						err.Error(),
						"No know 'fenixOrGatewayType' in TransmitEngine/Fenix",
					)
				}
			}
		}
	}
}

// ********************************************************************************************
// Try to dial parent gateway/Fenix or child gateway/plugin

func dialChildOrParenGateway(addressToDial string, dialSuccess bool, gRpcClientTowardsFenix gRPC.GatewayTowardsFenixClient, gRpcClientTowardPlugin gRPC.GatewayTowayPluginClient, remoteServerConnection *grpc.ClientConn, transmitOrDispatchEngineType string, pluginId string) (string, bool, gRPC.GatewayTowardsFenixClient, gRPC.GatewayTowayPluginClient, *grpc.ClientConn) {

	var err error

	// Get address to parent or child depending if we are running a TransmitEngine or a DispatchEngine
	switch transmitOrDispatchEngineType {
	case common_code.DispatchEngineTowardsPlugin:
		// Get address to next child gateway/plugin
		addressToDial = getClientAddressAndPort(pluginId)

	case common_code.TransmitEngineTowardsFenix:
		// Get adress to parent gateway/Fenix
		addressToDial = getParentAddressAndPort()

	default:
		// No known 'transmitOrDispatchEngineType'
		LogErrorAndSendInfoToFenix(
			"2c976770-ef52-4bd1-8a52-aaf1f64ebc69",
			gRPC.InformationMessage_FATAL,
			"No know 'transmitOrDispatchEngineType'",
			transmitOrDispatchEngineType,
			"No known 'transmitOrDispatchEngineType when dialing to child/parent",
			"No known 'transmitOrDispatchEngineType when dialing to child/parent",
		)
	}

	// Set up connection to Parent Gateway/Fenix or child gateway/plugin
	remoteServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())
	if err != nil {
		// Connection Not OK
		dialSuccess = false

		LogErrorAndSendInfoToFenix(
			"f1e2922b-d1b8-4c5d-9861-40ba71733383",
			gRPC.InformationMessage_WARNING,
			"addressToDial",
			addressToDial,
			err.Error(),
			"Coldn't dial Parent Gateway/Fenix or child Gateway/plugin",
		)
	} else {
		// Connection OK
		dialSuccess = true

		// Creates a new gateway gRPC-Client
		switch transmitOrDispatchEngineType {
		case common_code.DispatchEngineTowardsPlugin:
			// Creates a new gateway gRPC-Client towards Plugin
			gRpcClientTowardPlugin = gRPC.NewGatewayTowayPluginClient(remoteServerConnection)

		case common_code.TransmitEngineTowardsFenix:
			// Creates a new gateway gRPC-Client towards Fenix
			gRpcClientTowardsFenix = gRPC.NewGatewayTowardsFenixClient(remoteServerConnection)

		default:
			// No known 'transmitOrDispatchEngineType'
			LogErrorAndSendInfoToFenix(
				"0b819638-a8b7-4aa7-8d6c-a75627491829",
				gRPC.InformationMessage_FATAL,
				"No know 'transmitOrDispatchEngineType'",
				transmitOrDispatchEngineType,
				"No known 'transmitOrDispatchEngineType when creating gRPC-client",
				"No known 'transmitOrDispatchEngineType when creating gRPC-client",
			)
		}

	}

	// Return values to be used by calling function
	return addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection
}

// ********************************************************************************************
// Try to get ONE object from db and put on channel or get ALL objects and put on channel
//TODO Add internal structure for checking DB and delete message from DB when put on channel
func checkForSavedTemporaryObjectsInDbThatWillBePutOnChannel(onlyOneObject bool, channelType string, transmitOrDispatchEngineType string, channelToUse *common_code.GatewayChannelPackageStruct) (foundObjectsThatWherePutOnChannel bool) {

	// Decide if one or all objects thould be put back on channel
	if onlyOneObject == true {
		// Only one object
		foundObjectsThatWherePutOnChannel = checkForOneSavedTemporaryObjectsInDbAndPutOnChannel(channelType, transmitOrDispatchEngineType, channelToUse)
	} else {
		// All Objects
		var foundAtleastOneSavedObjectInDb = false
		for {
			foundObjectsThatWherePutOnChannel = checkForOneSavedTemporaryObjectsInDbAndPutOnChannel(channelType, transmitOrDispatchEngineType, channelToUse)
			if foundObjectsThatWherePutOnChannel == false {
				return foundAtleastOneSavedObjectInDb
			} else {
				foundAtleastOneSavedObjectInDb = true
			}
		}
	}

	// Code will never come here
	return true

}

// ********************************************************************************************
// Try to get one object from specific bucket and put on correct channel
//TODO Add internal structure for checking DB and DELETE message from DB when put on channel
func checkForOneSavedTemporaryObjectsInDbAndPutOnChannel(channelType string, transmitOrDispatchEngineType string, channelToUse *common_code.GatewayChannelPackageStruct) (foundObjectsThatWherePutOnChannel bool) {

	var (
		err    error
		bucket string
		id     string
		key    string
	)

	// Messages towards Plugin
	var (
		testInstructionMessageToBeForwardedTowardsPlugin                 *gRPC.TestInstruction_RT
		supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin *gRPC.SupportedTestDataDomainsRequest
	)

	// Messages towards Fenix
	var (
		informationMessageToBeForwardedTowardsFenix                         *gRPC.InformationMessage
		timeOutMessageToBeForwardedTowardsFenix                             *gRPC.TestInstructionTimeOutMessage
		testExecutionLogMessageToBeForwardedTowardsFenix                    *gRPC.TestExecutionLogMessage
		spportedTestDataDomainsMessageToBeForwardedTowardsFenix             *gRPC.SupportedTestDataDomainsMessage
		availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix     *gRPC.AvailbleTestInstructionAtPluginMessage
		availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix      *gRPC.AvailbleTestContainersAtPluginMessage
		testInstructionExecutionResultMessageToBeForwardedTowardsFenix      *gRPC.TestInstructionExecutionResultMessage
		supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix *gRPC.SupportedTestDataDomainsWithHeadersMessage
	)

	// Decide if it's TransmitEngine or a DispatchEngine
	switch transmitOrDispatchEngineType {

	// Choose DispatchEngine
	case common_code.DispatchEngineTowardsPlugin:
		//  Decide the correct DispatchEngine to use
		switch channelType {
		case common_code.ChannelTypeTestInstructionMessageTowardsPlugin:
			id = "d7151e7e-362e-44e8-a27e-fd8dc969e6f9"
			bucket = common_code.BucketForResendTestInstructionTowardsPlugin

		case common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
			id = "1fb35cdf-4689-4372-92d3-937f0420a0ae"
			bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin

		default:
			LogErrorAndSendInfoToFenix(
				"ce3f4832-e7cf-40da-bcb1-4e10110c7eae",
				gRPC.InformationMessage_FATAL,
				"No know 'channelType'",
				channelType,
				"No know 'channelType' in DispatchEngine",
				"No know 'channelType' in DispatchEngine",
			)
		}

		// Choose TransmitEngine
	case common_code.TransmitEngineTowardsFenix:
		//  Deside the correct TransmitEngine to use
		switch channelType {
		case common_code.ChannelTypeInformationMessageTowardsFenix:
			id = "63f924bf-b73f-4065-987d-821acb08bae4"
			bucket = common_code.BucketForResendOfInfoMessagesTowardsFenix

		case common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix:
			id = "f9337ae2-ec8b-494c-ad2d-ded4f1e67c15"
			bucket = common_code.BucketForResendOTimeOutMesagesTowardsFenix

		case common_code.ChannelTypeTestExecutionLogMessageTowardsFenix:
			id = "f4bc8914-e31b-4a52-91fe-bb1ab20a89e5"
			bucket = common_code.BucketForResendOfLogMesagesTowardsFenix

		case common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix:
			id = "dfde3eec-639a-4962-80d9-e660b40bd2d2"
			bucket = common_code.BucketForResendOfSupportedTestDataDomainsTowardsFenix

		case common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
			id = "70fba10a-aefd-4c7f-bb2d-e1fd4452f976"
			bucket = common_code.BucketForResendOfAvailableTestInstructionsTowardsFenix

		case common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
			id = "8e2ce6e9-13bc-4ea4-b51b-58d29cf54735"
			bucket = common_code.BucketForResendOfAvailableTestContainersTowardsFenix

		case common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix:
			id = "89badf8e-8c5e-47a3-9aa0-9dc4cdaf257a"
			bucket = common_code.BucketForResendOfTestInstructionExecutionResultTowardsFenix

		case common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
			id = "4ae89ab0-edad-4d7b-a97e-1dd20a3743ff"
			bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin

		default:
			LogErrorAndSendInfoToFenix(
				"fb66b765-442a-4bf8-b40c-f53b0514fd99",
				gRPC.InformationMessage_FATAL,
				"No know 'channelType'",
				channelType,
				"No know 'channelType' in TransmitEngine",
				"No know 'channelType' in TransmitEngine",
			)
		}

	default:
		// No known 'transmitOrDispatchEngineType'
		LogErrorAndSendInfoToFenix(
			"e63bc548-6238-4bfe-b466-d76af7d574c4",
			gRPC.InformationMessage_FATAL,
			"No know 'transmitOrDispatchEngineType'",
			transmitOrDispatchEngineType,
			"No known 'transmitOrDispatchEngineType when getting first object in bucket",
			"No known 'transmitOrDispatchEngineType when getting first object in bucket",
		)
	}

	// Create the channel that the database object should be sent back on
	returnChannel := make(chan common_code.DbResultMessageStruct)

	// Get first found object in bucket
	dbMessage := common_code.DbMessageStruct{
		common_code.DBGetFirstObjectFromBucket,
		bucket,
		nil,
		nil,
		returnChannel}

	// Send Read message to database to receive object
	dbMessageQueue <- dbMessage
	// Wait for object from channel, then close the channel
	databaseReturnMessage := <-returnChannel
	close(returnChannel)

	// Check if an error occured
	if databaseReturnMessage.Err != nil {
		// No object found
		logger.WithFields(logrus.Fields{
			"ID":     "b6cda913-364c-4296-b3c4-07b89de1a134",
			"Bucket": bucket,
		}).Debug("No object found when reading bucket")

		// Returns info that no object was found
		return false
	}

	// Unmarshal database object into correct data structure
	// Decide if it's TransmitEngine or a DispatchEngine
	switch transmitOrDispatchEngineType {

	// Choose DispatchEngine
	case common_code.DispatchEngineTowardsPlugin:
		//  Decide the correct DispatchEngine to use
		switch channelType {
		case common_code.ChannelTypeTestInstructionMessageTowardsPlugin:
			id = "abb58bad-ee31-4137-b686-1b0654428f5b"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &testInstructionMessageToBeForwardedTowardsPlugin)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.TestInstructionMessageChannelTowardsPlugin <- testInstructionMessageToBeForwardedTowardsPlugin

			// Delete saved message from DB
			id = "3937cf9b-b2ed-4fd8-bfe8-e67d67ea8a6e"
			bucket = common_code.BucketForResendTestInstructionTowardsPlugin
			key = testInstructionMessageToBeForwardedTowardsPlugin.MessageId

		case common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
			id = "63829f26-75aa-40a1-942f-550e97bf1731"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.SupportedTestDataDomainsRequestChannelTowardsPlugin <- supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin

			// Delete saved message from DB
			id = "7577bdbb-e8ca-431f-be51-b2a8563a2ddf"
			bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin
			key = supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin.MessageId

		default:
			LogErrorAndSendInfoToFenix(
				"fb276548-d087-4137-9a9f-9795c76008d7",
				gRPC.InformationMessage_FATAL,
				"No know 'channelType'",
				channelType,
				"No know 'channelType' in DispatchEngine",
				"No know 'channelType' in DispatchEngine",
			)
		}

		// Choose TransmitEngine
	case common_code.TransmitEngineTowardsFenix:
		//  Deside the correct TransmitEngine to use
		switch channelType {
		case common_code.ChannelTypeInformationMessageTowardsFenix:
			id = "bc704210-b491-4f71-8bc4-2f231ad6d27d"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &informationMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.InformationMessageChannelTowardsFenix <- informationMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "f9f427f4-172f-4211-af4f-94f6168da900"
			bucket = common_code.BucketForResendOfInfoMessagesTowardsFenix
			key = informationMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenix:
			id = "09af5d8a-211e-491e-80e2-31ecd5d21040"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &timeOutMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.TestInstructionTimeOutMessageChannelTowardsFenix <- timeOutMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "40275a25-d96c-42ae-9b62-aa5e129fa134"
			bucket = common_code.BucketForResendOTimeOutMesagesTowardsFenix
			key = timeOutMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeTestExecutionLogMessageTowardsFenix:
			id = "ddef7843-5671-49e6-83a5-66d7501dea0c"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &testExecutionLogMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.TestExecutionLogMessageChannelTowardsFenix <- testExecutionLogMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "79f3650c-1c2c-4b94-b972-7076b90f2a16"
			bucket = common_code.BucketForResendOfLogMesagesTowardsFenix
			key = testExecutionLogMessageToBeForwardedTowardsFenix.LogMessageId

		case common_code.ChannelTypeSupportedTestDataDomainsMessageTowardsFenix:
			id = "e7555770-551e-4bf3-9336-2b93bcd9ef46"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &spportedTestDataDomainsMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.SupportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix <- spportedTestDataDomainsMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "3c5fa6a6-8748-4549-826c-ee8cb65ed811"
			bucket = common_code.BucketForResendOfSupportedTestDataDomainsTowardsFenix
			key = spportedTestDataDomainsMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
			id = "c1c0b462-2d3a-4b79-8bd2-0be3a7d22818"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.AvailbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "2f9683b1-ff0a-49a3-9d8d-fc592ac7895c"
			bucket = common_code.BucketForResendOfAvailableTestInstructionsTowardsFenix
			key = availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
			id = "fe55e3b5-3888-4be4-aec8-da4fb5344f30"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.AvailbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "29925a0d-8c35-4f5c-b413-045ecff1a758"
			bucket = common_code.BucketForResendOfAvailableTestContainersTowardsFenix
			key = availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenix:
			id = "8097651f-aeb7-4e4d-b54c-e3c3fdd0357a"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &testInstructionExecutionResultMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.TestInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix <- testInstructionExecutionResultMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "c20a05e6-e985-446e-95f4-65081eb07829"
			bucket = common_code.BucketForResendOfTestInstructionExecutionResultTowardsFenix
			key = testInstructionExecutionResultMessageToBeForwardedTowardsFenix.MessageId

		case common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
			id = "797e8dd6-931f-4cad-a956-7aa4871ad5f4"

			// Convert saved json object into Go-struct
			err = json.Unmarshal(databaseReturnMessage.Value, &supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix)

			if err != nil {
				// Problem with unmarshal the json object
				// Send FATAL information to Fenix
				LogErrorAndSendInfoToFenix(
					id,
					gRPC.InformationMessage_FATAL,
					"Problem when unmarshaling object from DB",
					channelType,
					"Problem when unmarshaling object from DB",
					"Problem when unmarshaling object from DB",
				)
			}

			// Put message on channel
			channelToUse.SupportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix

			// Delete saved message from DB
			id = "77063b14-ade2-4e10-b8c3-a38918db5ef6"
			bucket = common_code.BucketForResendOfGetTestdataDomainsToPlugin
			key = supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix.MessageId

		default:
			LogErrorAndSendInfoToFenix(
				"cb1c4008-cb18-46cb-93a5-c2ae1aadcf7a",
				gRPC.InformationMessage_FATAL,
				"No know 'channelType'",
				channelType,
				"No know 'channelType' in TransmitEngine",
				"No know 'channelType' in TransmitEngine",
			)
		}

	default:
		// No known 'transmitOrDispatchEngineType'
		LogErrorAndSendInfoToFenix(
			"31e65376-91b4-47f7-8dfd-8330bc733e7d",
			gRPC.InformationMessage_FATAL,
			"No know 'transmitOrDispatchEngineType'",
			transmitOrDispatchEngineType,
			"No known 'transmitOrDispatchEngineType when getting first object in bucket",
			"No known 'transmitOrDispatchEngineType when getting first object in bucket",
		)
	}

	//TODO Färdigställ Delete

	// Delete message from Database
	// Create the channel that the database object should be sent back on
	returnChannel = make(chan common_code.DbResultMessageStruct)

	// Get first found object in bucket
	dbMessage = common_code.DbMessageStruct{
		common_code.DBDelete,
		bucket,
		key,
		nil,
		returnChannel}

	// Send Read message to database to receive object
	dbMessageQueue <- dbMessage
	// Wait for object from channel, then close the channel
	databaseReturnMessage = <-returnChannel
	close(returnChannel)

	// Check if an error occured
	if databaseReturnMessage.Err != nil {
		// No object found
		logger.WithFields(logrus.Fields{
			"ID":     "b6cda913-364c-4296-b3c4-07b89de1a134",
			"Bucket": bucket,
		}).Debug("No object found when reading bucket")

		LogErrorAndSendInfoToFenix(
			"b2c195fd-bc18-44b6-a107-1d875e1d536a",
			gRPC.InformationMessage_FATAL,
			"Didn't found key in bucket",
			"Didn't found key in bucket",
			"Didn't found key: '"+key+"' in bucket: '"+bucket+"' when trying to delete key after putting message on channel again",
			"Didn't found key: '"+key+"' in bucket: '"+bucket+"' when trying to delete key after putting message on channel again",
		)

		// Returns info that no object was found
		return false
	}

	//  Return info that one object was found
	return true
}
