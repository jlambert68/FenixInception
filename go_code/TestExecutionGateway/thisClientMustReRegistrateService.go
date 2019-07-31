package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"jlambert/FenixInception2/go_code/common_code"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix that this gateway should register itself to parent
//gatewayTowardsPluginObjectStruct
func (gRPCServerTowardsPlugin *common_code.GRPCServerTowardsPluginStruct) PleaseReRegisterClientAddress(ctx context.Context, reRegisterToGatewayMessage *gRPC.ReRegisterToGatewayMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse
	var parentgRPCAddress common_code.ParentgRPCAddressStruct

	common_code.logger.WithFields(logrus.Fields{
		"ID":                         "b9e6bde2-0a59-4459-83c4-d723d50a080c",
		"reRegisterToGatewayMessage": reRegisterToGatewayMessage,
	}).Info("Incoming gRPC: 'PleaseReRegisterClientAddress'")

	//Move data into object that should be save in DB
	parentgRPCAddress.ParentGatewayId = reRegisterToGatewayMessage.GatewayId
	parentgRPCAddress.ParentGatewayName = reRegisterToGatewayMessage.GatewayName
	parentgRPCAddress.ParentGatewayServerAddress = reRegisterToGatewayMessage.GatewayAddress
	parentgRPCAddress.ParentGatewayServerPort = reRegisterToGatewayMessage.GatewayPort
	parentgRPCAddress.CreatedDateTime = generaTimeStampUTC()

	// Convert Parent Gateway address info-struct into a byte array
	parentgRPCAddressByteArray, err := json.Marshal(reRegisterToGatewayMessage)
	if err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"ID":                "6df384d3-ebe9-4a65-947f-e51814c8544c",
			"ParentgRPCAddress": parentgRPCAddress,
			"err":               err,
		}).Error("Error when converting 'ParentgRPCAddress' into a byte array, stopping futher processing of Reregistration.")

		// Send Error information to Fenix
		common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         common_code.gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       common_code.gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 common_code.gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               common_code.gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Error when converting 'ParentgRPCAddress' into a byte array, stopping futher processing of Reregistration.",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   common_code.gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: common_code.gatewayConfig.SystemDomain.GatewayDomainName,
		}

		returnMessage.Comments = "Error when converting 'ParentgRPCAddress' into a byte array, stopping futher processing of Reregistration."
		returnMessage.Acknack = false
		return returnMessage, nil

	}

	// Save parentgRPCAddressByteArray to local database, using local channel
	// Return Channel
	returnChannel := make(chan common_code.dbResultMessageStruct)

	dbMessage := common_code.dbMessageStruct{
		common_code.DbWrite,
		common_code.BucketForParentAddress,
		common_code.BucketKeyForParentAddress, // Key allways hardcoded due to one gateway or plugin can only have one parent
		parentgRPCAddressByteArray,
		returnChannel}

	// Send message to Database
	common_code.dbMessageQueue <- dbMessage

	// Wait for result on result channel then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.err != nil {
		common_code.logger.WithFields(logrus.Fields{
			"ID":  "446cb366-9d84-4c1c-a628-ef162f7c1747",
			"err": err,
		}).Error("Got an error when Saveing to local DB")

		// Send Error information to Fenix
		common_code.informationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         common_code.gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       common_code.gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 common_code.gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               common_code.gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Got an error when Saveing to local DB",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   common_code.gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: common_code.gatewayConfig.SystemDomain.GatewayDomainName,
		}

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Error when saving to local DB, stopping futher processing of Reregistration"
		returnMessage.Acknack = false
		return returnMessage, nil
	}

	common_code.logger.WithFields(logrus.Fields{
		"ID": "fdf7081a-e7da-4bf1-a87c-82c51b8f575b",
	}).Debug("Reregistration info was saved in local database")

	// Start Registration of this Gateway/Plugin
	// **** TODO Put message on channel to register this client ****
	registerSuccess, err := registerThisGatewayAtParentGateway()

	if registerSuccess == true {
		common_code.logger.WithFields(logrus.Fields{
			"ID": "14ae3650-c0dd-4e23-b197-9706d5dfc8bd",
		}).Debug("Rereregistration to parent gateway/Fenix was successful")

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Reregistration was done"
		returnMessage.Acknack = true
		return returnMessage, nil

	} else {
		common_code.logger.WithFields(logrus.Fields{
			"ID":  "e5c273a8-f257-4c9f-bcbe-1697a2de1663",
			"err": err,
		}).Error("Rereregistration to parent gateway/Fenix could not be done")

		returnMessage.Comments = "Reregistration could not be done"
		returnMessage.Acknack = false
		return returnMessage, nil
	}

}
