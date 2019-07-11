package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"time"
)

// ********************************************************************************************
// Initiate Transmit Engines for messages towards Fenix'
//

func initiateTransmitEnginesTowardsFenix() {

	// *** Towards Plugin ***
	// Start a Dispatch Engine, for 'TestInstructionMessageTowardsPlugin' as a go-routine
	go transmitAndDispatchEngine(channelTypeTestInstructionMessageTowardsPlugin, dispatchEngineTowardsPlugin)

	// Start a Dispatch Engine, for 'SupportedTestDataDomainsRequestMessageTowardsPlugin,' as a go-routine
	go transmitAndDispatchEngine(channelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin, dispatchEngineTowardsPlugin)

	// *** Towards Fenix ***
	// Start a Transmit Engine, for 'informationMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeInformationMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'timeOutMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeTestInstructionTimeOutMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'spportedTestDataDomainsMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeTestExecutionLogMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestInstructionAtPluginMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeSupportedTestDataDomainsMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeAvailbleTestContainersAtPluginMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'testInstructionExecutionResultMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeTestInstructionExecutionResultMessageTowardsFenix, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' as a go-routine
	go transmitAndDispatchEngine(channelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix, transmitEngineTowardsFenix)

}

// ********************************************************************************************
// Forward X-message from channel towards Fenix
//

func transmitAndDispatchEngine(channelType string, transmitOrDispatchEngineType string) {

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

	var gRRCerr error
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
			time.Sleep(ServieSleepTime * time.Second)

			// Check it is possible to dial parent gateway/Fenix or child gatway/plugin
			addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection = dialChildOrParenGateway(addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection, transmitOrDispatchEngineType, pluginId)

		} else {

			// Check for messages, in DB, to Resend (If so then put them on channel)
			objectsWherPutOnChannel := checkForSavedTemporaryObjectsInDbThatWillBePutOnChannel(channelType, transmitOrDispatchEngineType)
			if objectsWherPutOnChannel == true {
				logger.WithFields(logrus.Fields{
					"ID":                      "84821702-66f1-48bc-a8cd-eee8951fb880",
					"objectsWherPutOnChannel": objectsWherPutOnChannel,
				}).Debug("DB checked for saved objects and there were at least one object saved in DB that were put on channel again")
			} else {
				logger.WithFields(logrus.Fields{
					"ID":                      "84821702-66f1-48bc-a8cd-eee8951fb880",
					"objectsWherPutOnChannel": objectsWherPutOnChannel,
				}).Debug("DB checked for saved objects and there were NO objects saved in DB")
			}

			// Run service and process messages

			// Wait for message to come, from defined channel, to transmit engine ****
			// Select on messageType and convert message into byteArray
			// Set variables depending on messageType

			// Decide if it's TransmitEngine or a DispatchEngine
			switch transmitOrDispatchEngineType {

			// Choose DispatchEngine
			case dispatchEngineTowardsPlugin:
				//  Decide the correct DispatchEngine to use
				switch channelType {
				case channelTypeTestInstructionMessageTowardsPlugin:
					testInstructionMessageToBeForwardedTowardsPlugin = <-testInstructionMessageChannelTowardsPlugin
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
					bucket = BucketForResendTestInstructionTowardsPlugin
					NumberOfMessagesInChannel = len(testInstructionMessageChannelTowardsPlugin)

				case channelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
					supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin = <-supportedTestDataDomainsRequestChannelTowardsPlugin
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
					bucket = BucketForResendOfGetTestdataDomainsToPlugin
					NumberOfMessagesInChannel = len(supportedTestDataDomainsRequestChannelTowardsPlugin)

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
			case transmitEngineTowardsFenix:
				//  Deside the correct TransmitEngine to use
				switch channelType {
				case channelTypeInformationMessageTowardsFenix:
					informationMessageToBeForwardedTowardsFenix = <-informationMessageChannelTowardsFenix
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
					bucket = BucketForResendOfInfoMessagesTowardsFenix
					NumberOfMessagesInChannel = len(informationMessageChannelTowardsFenix)

				case channelTypeTestInstructionTimeOutMessageTowardsFenix:
					timeOutMessageToBeForwardedTowardsFenix = <-testInstructionTimeOutMessageChannelTowardsFenix
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
					bucket = BucketForResendOTimeOutMesagesTowardsFenix
					NumberOfMessagesInChannel = len(testInstructionTimeOutMessageChannelTowardsFenix)

				case channelTypeTestExecutionLogMessageTowardsFenix:
					testExecutionLogMessageToBeForwardedTowardsFenix = <-testExecutionLogMessageChannelTowardsFenix
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
					bucket = BucketForResendOfLogMesagesTowardsFenix
					NumberOfMessagesInChannel = len(testExecutionLogMessageChannelTowardsFenix)

				case channelTypeSupportedTestDataDomainsMessageTowardsFenix:
					spportedTestDataDomainsMessageToBeForwardedTowardsFenix = <-supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix
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
					bucket = BucketForResendOfSupportedTestDataDomainsTowardsFenix
					NumberOfMessagesInChannel = len(supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix)

				case channelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
					availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix = <-availbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix
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
					bucket = BucketForResendOfAvailableTestInstructionsTowardsFenix
					NumberOfMessagesInChannel = len(availbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix)

				case channelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
					availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix = <-availbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix
					messageToBeForwardedByteArray, err = json.Marshal(*availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix)
					// If error when marshaling then use the following information
					id = "14b2aa63-7e69-4e80-9d6b-796c286997c9"
					infoHeader = "availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix"
					info = ""
					errorMessage = ""
					message = "availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix"

					messageId = availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix.MessageId
					bucket = BucketForResendOfAvailableTestContainersTowardsFenix
					NumberOfMessagesInChannel = len(availbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix)

				case channelTypeTestInstructionExecutionResultMessageTowardsFenix:
					testInstructionExecutionResultMessageToBeForwardedTowardsFenix = <-testInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix
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
					bucket = BucketForResendOfTestInstructionExecutionResultTowardsFenix
					NumberOfMessagesInChannel = len(testInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix)

				case channelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
					supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix = <-supportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix
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
					bucket = BucketForResendOfGetTestdataDomainsToPlugin
					NumberOfMessagesInChannel = len(supportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix)

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

				// Do the gRPC-call to parent gateway/Fenix or child gatway/plugin
				// Decide if it's TransmitEngine or a DispatchEngine
				switch transmitOrDispatchEngineType {

				// Choose DispatchEngine
				case dispatchEngineTowardsPlugin:
					//  Decide the correct DispatchEngine to use
					switch channelType {
					case channelTypeTestInstructionMessageTowardsPlugin:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardPlugin.SendTestInstructionTowardsPlugin(gRpcContexType, testInstructionMessageToBeForwardedTowardsPlugin)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsRequestMessageTowardsPlugin:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardPlugin.GetSupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsRequestMessageToBeForwardedTowardsPlugin)
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
				case transmitEngineTowardsFenix:
					//  Deside the correct TransmitEngine to use
					switch channelType {
					case channelTypeInformationMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestInstructionTimeOutMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestInstructionTimeOutTowardsFenix(gRpcContexType, timeOutMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestExecutionLogMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestExecutionLogTowardsFenix(gRpcContexType, testExecutionLogMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegistrateAvailableTestDataDomains(gRpcContexType, spportedTestDataDomainsMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegisterAvailbleTestInstructions(gRpcContexType, availbleTestInstructionAtPluginMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeAvailbleTestContainersAtPluginMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegistrateAailableTestContainers(gRpcContexType, availbleTestContainersAtPluginMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestInstructionExecutionResultMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestInstructionResultTowardsFenix(gRpcContexType, testInstructionExecutionResultMessageToBeForwardedTowardsFenix)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenix:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsWithHeadersMessageToBeForwardedTowardsFenix)
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

				if gRRCerr != nil {
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
	case dispatchEngineTowardsPlugin:
		// Get address to next child gateway/plugin
		addressToDial = getClientAddressAndPort(pluginId)

	case transmitEngineTowardsFenix:
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
		case dispatchEngineTowardsPlugin:
			// Creates a new gateway gRPC-Client towards Plugin
			gRpcClientTowardPlugin = gRPC.NewGatewayTowayPluginClient(remoteServerConnection)

		case transmitEngineTowardsFenix:
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
// Try to dial parent gateway/Fenix or child gateway/plugin
//TODO Add internal structure for checking DB and delete message from DB when put on channel
func checkForSavedTemporaryObjectsInDbThatWillBePutOnChannel(channelType string, transmitOrDispatchEngineType string) (foundObjectsThatWherePutOnChannel bool) {
	return true
}
