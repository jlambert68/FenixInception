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
// Initiate Dispatch Engine for TestInstructions to be executed
//

func (gatewayTowardsPluginObject *gatewayTowardsPluginObjectStruct) initiateDispatchEngineForTestInstructiona() {

	// Start Dispatch Engine, for TestInstruction to be executed, as a go-routine
	go gatewayTowardsPluginObject.dispatchEngineForTestInstructions()
}

// ********************************************************************************************
// Forward TestInstructions from incoming channel to child gateway or plugin
//

func (gatewayTowardsPluginObject *gatewayTowardsPluginObjectStruct) dispatchEngineForTestInstructions() {

	for {

		// Service can be started and Stopped by central control of the Gateway
		if gatewayMustStopProcessing == true {
			// Service should be stopped from processing any messages
			time.Sleep(ServieSleepTime * time.Second)

		} else {
			// Run service and process messages

			// Wait for data comes from channel to dispatch engine
			testInstructionToBeForwarded := <-gatewayTowardsPluginObject.testInstructionMessageChannel

			// Check number of messages in channel
			channelSinaling(len(gatewayTowardsPluginObject.testInstructionMessageChannel),
				"testInstructionTimeOutMessageChannel",
				"9c89b05d-7f4d-4971-b21b-7c7fa9d42e9c")

			logger.WithFields(logrus.Fields{
				"ID":                           "f015de87-d51f-4025-8e27-cf818bbe255d",
				"testInstructionToBeForwarded": testInstructionToBeForwarded,
			}).Debug("Received a new TestInstruction from channel that shoud be forwarded")

			// Send TestInstruction to client using gRPC-call
			addressToDial := getClientAddressAndPort(testInstructionToBeForwarded.PluginId)

			// Set up connection to Client Gateway or Plugin
			remoteChildServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
			if err != nil {
				logger.WithFields(logrus.Fields{
					"ID":            "3ad75b4a-4c7d-44eb-b781-c8cfd18c4fcf",
					"addressToDial": addressToDial,
					"error message": err,
				}).Error("Did not connect to Child (Gateway or Plugin) Server!")

				//Send Error information to Fenix
				gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
					OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
					OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
					SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
					SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
					MessageId:                generateUUID(),
					MessageType:              gRPC.InformationMessage_ERROR,
					Message:                  "Did not connect to Child (Gateway or Plugin) Server!",
					OrginalCreateDateTime:    generaTimeStampUTC(),
					OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
					OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
				}
			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "4510ae97-1753-4d4c-a7f5-3987054bd969",
					"addressToDial": addressToDial,
				}).Debug("gRPC connection OK to child-gateway- or Plugin-Server!")

				// Creates a new gateway Client
				gatewayClient := gRPC.NewGatewayTowayPluginClient(remoteChildServerConnection)

				// ChangeSenderId to this gatway's SenderId before sending the data forward
				testInstructionToBeForwarded.SenderId = gatewayConfig.GatewayIdentification.GatewayId
				testInstructionToBeForwarded.SenderName = gatewayConfig.GatewayIdentification.GatewayName

				// Do gRPC-call to client gateway or Plugin
				ctx := context.Background()
				returnMessage, err := gatewayClient.SendTestInstructionTowardsPlugin(ctx, testInstructionToBeForwarded)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"ID":            "4a0f0753-6394-4115-aa56-cf26c9354f83",
						"returnMessage": returnMessage,
						"error":         err,
					}).Error("Problem to send TestInstruction to child-Gateway or Plugin")

					//Send Error information to Fenix
					gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
						OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
						OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
						SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
						SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
						MessageId:                generateUUID(),
						MessageType:              gRPC.InformationMessage_ERROR,
						Message:                  "Problem to send TestInstruction to child-Gateway or Plugin",
						OrginalCreateDateTime:    generaTimeStampUTC(),
						OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
						OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
					}
				} else {
					logger.WithFields(logrus.Fields{
						"ID":            "ef7f99b4-dcc0-4583-b16a-de3bf093880f",
						"addressToDial": addressToDial,
					}).Debug("gRPC-send TestInstruction to child-Gateway or Plugin")

					// Set New Database Status and Timestamp on TestInstruction
					testInstructionToBeForwarded.Metadata.DbStatusEnum = gRPC.DbStatusEnum_TESTINSTRUCTION_SENT_TO_CHILD
					testInstructionToBeForwarded.Metadata.LastUpdatedDateTimeInDB = time.Now().String()

					// Convert TestInstruction struct into a byte array
					testInstructionByteArray, err := json.Marshal(testInstructionToBeForwarded)
					if err != nil {
						logger.WithFields(logrus.Fields{
							"ID":              "b6cbe063-65a1-4d79-90f1-9849fb830e16",
							"testInstruction": testInstructionToBeForwarded,
							"err":             err,
						}).Error("Error when converting TestInstruction into a byte array, stopping futher processing of this TestInstruction")

						//TSend Error information to Fenix
						gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
							OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
							OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
							SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
							SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
							MessageId:                generateUUID(),
							MessageType:              gRPC.InformationMessage_ERROR,
							Message:                  "Error when converting TestInstruction into a byte array, stopping futher processing of this TestInstruction",
							OrginalCreateDateTime:    generaTimeStampUTC(),
							OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
							OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
						}

					} else {

						// Save TestInstruction to local database, using local channel
						returnChannel := make(chan dbResultMessageStruct)

						dbMessage := dbMessageStruct{
							DbWrite,
							"TestInstructions",
							testInstructionToBeForwarded.TestInstructionGuid,
							testInstructionByteArray,
							returnChannel}

						// Send message to Database
						dbMessageQueue <- dbMessage

						// Wait for result on result channel and then close returnChannel
						returnDBMessage := <-returnChannel
						close(returnChannel)

						if returnDBMessage.err != nil {
							logger.WithFields(logrus.Fields{
								"ID":  "69e2f6bd-5ca3-4292-893d-075ff5ac79e0",
								"err": err,
							}).Error("Got an error when Saveing to local DB")

							//Send Error information to Fenix
							gatewayTowardsFenixObject.informationMessageChannel <- &gRPC.InformationMessage{
								OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
								OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
								SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
								SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
								MessageId:                generateUUID(),
								MessageType:              gRPC.InformationMessage_ERROR,
								Message:                  "Got an error when Saveing to local DB",
								OrginalCreateDateTime:    generaTimeStampUTC(),
								OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
								OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
							}

						}

					}
				}
			}
		}
	}
}
