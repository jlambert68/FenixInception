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
// Initiate Transmit Engine for XTowardsFenix
//
/*
func (gatewayObject *gatewayTowardsFenixObjectStruct) initiateTransmitEngineForXTowardsFenix() {

	// Start Transmit Engine, for X as a go-routine
	go gatewayObject.transmitEngineForXTowardsFenix(gatewayObject.testInstructionExecutionResultMessageTowardsFenixChannel, types.Chan{})
}testInstructionExecutionResultMessageTowardsFenixChanneltestInstructionExecutionResultMessageTowardsFenixChannel
Ä*/
// ********************************************************************************************
// Forward X-messages-messages from channel towards Fenix
//

// TODO initiate all channels at once and set channelType at the same time
func (gatewayObject *gatewayTowardsFenixObjectStruct) transmitEngineForXTowardsFenix(channelType string) {

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
	var NumberOfMessagesInCHannel int

	var (
		addressToDial                string
		remoteParentServerConnection *grpc.ClientConn
		gatewayClient                gRPC.GatewayTowardsFenixClient
		gRpcContexType               context.Context
		dialSuccess                  = false
	)

	var gRRCerr error
	var returnMessageAckNackResponse *gRPC.AckNackResponse
	var returnMessageString string

	// Use background call in gRPC-call
	gRpcContexType = context.Background()

	// Check it is possible to dial parent gateway/Fenix
	addressToDial, dialSuccess, gatewayClient, remoteParentServerConnection = dialParenGateway(addressToDial, dialSuccess, gatewayClient, remoteParentServerConnection)

	// Get address to parent
	addressToDial = getParentAddressAndPort()

	// Set up connection to Parent Gateway or Fenix
	remoteParentServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())
	if dialSuccess == true {
		// Connection OK

		// Wait for messahe on channel
		// Do gRPC-call to client gateway or Fenix

		for {

			// Service can be started and Stopped by central control of the Gateway
			if gatewayMustStopProcessing == true || dialSuccess == false {
				// Service should be stopped from processing any messages
				time.Sleep(ServieSleepTime * time.Second)

				// Check it is possible to dial parent gateway/Fenix
				addressToDial, dialSuccess, gatewayClient, remoteParentServerConnection = dialParenGateway(addressToDial, dialSuccess, gatewayClient, remoteParentServerConnection)

			} else {
				// Run service and process messages

				// Wait for message to come, from defined channel, to transmit engine ****
				// Select on messageType and convert message into byteArray
				// Set variables depending on messageType
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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

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
					NumberOfMessagesInCHannel = len(gatewayTowardsFenixObject.informationMessageChannel)

				default:
					LogErrorAndSendInfoToFenix(
						"12c6cb68-9fd3-4629-91af-ac0fe4b9ede8",
						gRPC.InformationMessage_FATAL,
						"No know 'channelType'",
						channelType,
						err.Error(),
						"No know 'channelType' in TransmitEngine",
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
					channelSinaling(NumberOfMessagesInCHannel,
						infoHeader,
						"db1bd71e-4ad4-40df-aba5-e80f69899e08")

					// Do the gRPC-call to parent gateway or Fenix
					switch channelType {
					case channelTypeInformationMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.SendMessageToFenix(gRpcContexType, informationMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestInstructionTimeOutMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.SendTestInstructionTimeOutTowardsFenix(gRpcContexType, timeOutMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestExecutionLogMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.SendTestExecutionLogTowardsFenix(gRpcContexType, testExecutionLogMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.RegistrateAvailableTestDataDomains(gRpcContexType, spportedTestDataDomainsMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeAvailbleTestInstructionsAtPluginMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.RegisterAvailbleTestInstructions(gRpcContexType, availbleTestInstructionAtPluginMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeAvailbleTestContainersAtPluginMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.RegistrateAailableTestContainers(gRpcContexType, availbleTestContainersAtPluginMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeTestInstructionExecutionResultMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.SendTestInstructionResultTowardsFenix(gRpcContexType, testInstructionExecutionResultMessageToBeForwarded)
						returnMessageString = returnMessageAckNackResponse.String()

					case channelTypeSupportedTestDataDomainsWithHeadersMessage:
						returnMessageAckNackResponse, gRRCerr = gatewayClient.SupportedTestDataDomains(gRpcContexType, supportedTestDataDomainsWithHeadersMessageToBeForwarded)
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

					if gRRCerr != nil {
						// Error when sending gRPC to parent
						dialSuccess = false

						LogErrorAndSendInfoToFenix(
							//TODO SPlit logs into one log per messageType
							"519cc9e7-5a24-4787-a23e-e58808c22ce5",
							gRPC.InformationMessage_WARNING,
							"returnMessage",
							returnMessageString,
							err.Error(),
							"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
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
							"ID":            "8820b7be-2495-44bb-8894-ea9b39486848",
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
// Try to dial parent gateway/Fenix
// input pointer are updated back to their original variables

func dialParenGateway(addressToDial string, dialSuccess bool, gatewayClient gRPC.GatewayTowardsFenixClient, remoteParentServerConnection *grpc.ClientConn) (string, bool, gRPC.GatewayTowardsFenixClient, *grpc.ClientConn) {

	var err error
	// Get address to parent
	addressToDial = getParentAddressAndPort()

	// Set up connection to Parent Gateway or Fenix
	remoteParentServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())
	if err != nil {
		// Connection Not OK
		dialSuccess = false

		LogErrorAndSendInfoToFenix(
			"16738894-468a-4125-a9cc-32fd97afb767",
			gRPC.InformationMessage_WARNING,
			"addressToDial",
			addressToDial,
			err.Error(),
			"COuldn't dial Parent (Gateway or Plugin) Server!",
		)
	} else {
		// Connection OK
		dialSuccess = true

		// Creates a new gateway Client
		gatewayClient = gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

	}

	return addressToDial, dialSuccess, gatewayClient, remoteParentServerConnection
}
