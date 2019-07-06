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

func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateTransmitEnginesTowardsFenix() {

	// Start a Transmit Engine, for 'informationMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeInformationMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'timeOutMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeTestInstructionTimeOutMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'spportedTestDataDomainsMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeTestExecutionLogMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestInstructionAtPluginMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeSupportedTestDataDomainsMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeAvailbleTestInstructionsAtPluginMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'availbleTestContainersAtPluginMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeAvailbleTestContainersAtPluginMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'testInstructionExecutionResultMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeTestInstructionExecutionResultMessage, transmitEngineTowardsFenix)

	// Start a Transmit Engine, for 'supportedTestDataDomainsWithHeadersMessageToBeForwarded' as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(channelTypeSupportedTestDataDomainsWithHeadersMessage, transmitEngineTowardsFenix)

}

// ********************************************************************************************
// Forward X-message from channel towards Fenix
//

func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForXTowardsFenix(channelType string, transmitOrDispatchEngineType string) {

	var (
		informationMessageToBeForwarded                         *gRPC.InformationMessage
		timeOutMessageToBeForwarded                             *gRPC.TestInstructionTimeOutMessage
		testExecutionLogMessageToBeForwarded                    *gRPC.TestExecutionLogMessage
		spportedTestDataDomainsMessageToBeForwarded             *gRPC.SupportedTestDataDomainsMessage
		availbleTestInstructionAtPluginMessageToBeForwarded     *gRPC.AvailbleTestInstructionAtPluginMessage
		availbleTestContainersAtPluginMessageToBeForwarded      *gRPC.AvailbleTestContainersAtPluginMessage
		testInstructionExecutionResultMessageToBeForwarded      *gRPC.TestInstructionExecutionResultMessage
		supportedTestDataDomainsWithHeadersMessageToBeForwarded *gRPC.SupportedTestDataDomainsWithHeadersMessage
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

		gRpcContexType         context.Context
		dialSuccess            = false
	)

	var gRRCerr error
	var returnMessageAckNackResponse *gRPC.AckNackResponse
	var returnMessageString string

	// Use background call in gRPC-call
	gRpcContexType = context.Background()

	// Check it is possible to dial parent gateway/Fenix or child gatway/plugin
	addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection = dialChildOrParenGateway(addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection, transmitOrDispatchEngineType, pluginId)

	if dialSuccess == true {
		// Connection OK

		// Wait for messahe on channel
		// Do gRPC-call to client gateway or Fenix

		for {

			// Service can be started and Stopped by central control of the Gateway
			if gatewayMustStopProcessing == true || dialSuccess == false {
				// Service should be stopped from processing any messages
				time.Sleep(ServieSleepTime * time.Second)

				// Check it is possible to dial parent gateway/Fenix or child gatway/plugin
				addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection = dialChildOrParenGateway(addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTowardPlugin, remoteServerConnection, transmitOrDispatchEngineType, pluginId)

			} else {
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
					case channelTypeTestInstructionMessageChannel:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsRequestMessage:
						returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

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
					case channelTypeInformationMessage:
						informationMessageToBeForwarded = <-gatewayTowardsFenixObject.informationMessageChannel
						messageToBeForwardedByteArray, err = json.Marshal(*informationMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "6798902d-36b9-4598-a125-4be3abd4a23b"
						infoHeader = "InformationMessageToBeForwared"
						info = ""
						errorMessage = ""
						message = "InformationMessageToBeForwared"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						informationMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						informationMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = informationMessageToBeForwarded.MessageId
						bucket = BucketForResendOfInfoMessagesToFenix
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeTestInstructionTimeOutMessage:
						timeOutMessageToBeForwarded = <-gatewayTowardsFenixObject.testInstructionTimeOutMessageChannel
						messageToBeForwardedByteArray, err = json.Marshal(*timeOutMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "e80978f1-2f2c-43f7-bd4b-218e852f2f72"
						infoHeader = "TestInstructionTimeOutMessageToBeForwared"
						info = ""
						errorMessage = ""
						message = "TestInstructionTimeOutMessageToBeForwared"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						timeOutMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						timeOutMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = timeOutMessageToBeForwarded.MessageId
						bucket = BucketForResendOTimeOutMesagesToFenix
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeTestExecutionLogMessage:
						testExecutionLogMessageToBeForwarded = <-gatewayTowardsFenixObject.testExecutionLogMessageChannel
						messageToBeForwardedByteArray, err = json.Marshal(*testExecutionLogMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "3f8b01b6-da9b-4cee-b787-6a0e9efe9fed"
						infoHeader = "TestExecutionLogMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "TestExecutionLogMessageToBeForwarded"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						testExecutionLogMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						testExecutionLogMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = testExecutionLogMessageToBeForwarded.LogMessageId
						bucket = BucketForResendOfLogMesagesToFenix
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeSupportedTestDataDomainsMessage:
						spportedTestDataDomainsMessageToBeForwarded = <-gatewayTowardsFenixObject.supportedTestDataDomainsMessageTowardsFenixChannel
						messageToBeForwardedByteArray, err = json.Marshal(*spportedTestDataDomainsMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "d3a33be3-6cf8-459d-bda2-ff64c53d249a"
						infoHeader = "TestDataDomainsMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "TestDataDomainsMessageToBeForwarded"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						spportedTestDataDomainsMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						spportedTestDataDomainsMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = spportedTestDataDomainsMessageToBeForwarded.MessageId
						bucket = BucketForResendOfSupportedTestDataDomains
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeAvailbleTestInstructionsAtPluginMessage:
						availbleTestInstructionAtPluginMessageToBeForwarded = <-gatewayTowardsFenixObject.availbleTestInstructionAtPluginMessageTowardsFenixChannel
						messageToBeForwardedByteArray, err = json.Marshal(*availbleTestInstructionAtPluginMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "6524a739-e628-4265-a6cb-d64597a20c3e"
						infoHeader = "AvailbleTestInstructionsAtPluginMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "AvailbleTestInstructionsAtPluginMessageToBeForwarded"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						availbleTestInstructionAtPluginMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						availbleTestInstructionAtPluginMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = availbleTestInstructionAtPluginMessageToBeForwarded.MessageId
						bucket = BucketForResendOfAvailableTestInstructionsToFenix
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeAvailbleTestContainersAtPluginMessage:
						availbleTestContainersAtPluginMessageToBeForwarded = <-gatewayTowardsFenixObject.availbleTestContainersAtPluginMessageTowardsFenixChannel
						messageToBeForwardedByteArray, err = json.Marshal(*availbleTestContainersAtPluginMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "14b2aa63-7e69-4e80-9d6b-796c286997c9"
						infoHeader = "AvailbleTestContainersAtPluginMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "AvailbleTestContainersAtPluginMessageToBeForwarded"

						messageId = availbleTestContainersAtPluginMessageToBeForwarded.MessageId
						bucket = BucketForResendOfAvailableTestContainers
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeTestInstructionExecutionResultMessage:
						testInstructionExecutionResultMessageToBeForwarded = <-gatewayTowardsFenixObject.testInstructionExecutionResultMessageTowardsFenixChannel
						messageToBeForwardedByteArray, err = json.Marshal(*testInstructionExecutionResultMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "72d128f7-e84c-408a-87a9-e0c8b84c6841"
						infoHeader = "TestInstructionExecutionResultMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "TestInstructionExecutionResultMessageToBeForwarded"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						testInstructionExecutionResultMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						testInstructionExecutionResultMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = testInstructionExecutionResultMessageToBeForwarded.MessageId
						bucket = BucketForResendOfTestInstructionExecutionResult
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

					case channelTypeSupportedTestDataDomainsWithHeadersMessage:
						supportedTestDataDomainsWithHeadersMessageToBeForwarded = <-gatewayTowardsFenixObject.supportedTestDataDomainsWithHeadersMessageTowardsFenixChannel
						messageToBeForwardedByteArray, err = json.Marshal(*supportedTestDataDomainsWithHeadersMessageToBeForwarded)
						// If error when marshaling then use the following information
						id = "89df74d0-8ba7-47df-bddc-5864a9b376f1"
						infoHeader = "SupportedTestDataDomainsWithHeadersMessageToBeForwarded"
						info = ""
						errorMessage = ""
						message = "SupportedTestDataDomainsWithHeadersMessageToBeForwarded"

						// ChangeSenderId to this gatway's SenderId before sending the data forward
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
						supportedTestDataDomainsWithHeadersMessageToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

						// Information used when saving message to local database
						messageId = supportedTestDataDomainsWithHeadersMessageToBeForwarded.MessageId
						bucket = BucketForResendOfGetTestdataDomainsToPlugin
						NumberOfMessagesInChannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
						//"Error when converting 'testInstructionExecutionResultMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
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
						case channelTypeTestInstructionMessageChannel:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix., informationMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeSupportedTestDataDomainsRequestMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwarded)
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
						case channelTypeInformationMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendMessageToFenix(gRpcContexType, informationMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeTestInstructionTimeOutMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestInstructionTimeOutTowardsFenix(gRpcContexType, timeOutMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeTestExecutionLogMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestExecutionLogTowardsFenix(gRpcContexType, testExecutionLogMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeSupportedTestDataDomainsMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegistrateAvailableTestDataDomains(gRpcContexType, spportedTestDataDomainsMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeAvailbleTestInstructionsAtPluginMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegisterAvailbleTestInstructions(gRpcContexType, availbleTestInstructionAtPluginMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeAvailbleTestContainersAtPluginMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.RegistrateAailableTestContainers(gRpcContexType, availbleTestContainersAtPluginMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeTestInstructionExecutionResultMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SendTestInstructionResultTowardsFenix(gRpcContexType, testInstructionExecutionResultMessageToBeForwarded)
							returnMessageString = returnMessageAckNackResponse.String()

						case channelTypeSupportedTestDataDomainsWithHeadersMessage:
							returnMessageAckNackResponse, gRRCerr = gRpcClientTowardsFenix.SupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsWithHeadersMessageToBeForwarded)
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

						// TODO Check for messages to Resend (If so then put them on channel)

					}
				}
			}
		}
	}
}

// ********************************************************************************************
// Try to dial parent gateway/Fenix or child gateway/plugin

func
dialChildOrParenGateway(addressToDial string, dialSuccess bool, gRpcClientTowardsFenix gRPC.GatewayTowardsFenixClient, gRpcClientTOwardPlugin gRPC.GatewayTowayPluginClient, remoteServerConnection *grpc.ClientConn, transmitOrDispatchEngineType string, pluginId string) (string, bool, gRPC.GatewayTowardsFenixClient, gRPC.GatewayTowayPluginClient, *grpc.ClientConn, ) {

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
			gRpcClientTOwardPlugin = gRPC.NewGatewayTowayPluginClient(remoteServerConnection)

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
	return addressToDial, dialSuccess, gRpcClientTowardsFenix, gRpcClientTOwardPlugin, remoteServerConnection
}
