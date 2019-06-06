package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix that this gateway should register itself to parent
//
func (gatewayObject *gatewayTowardsPluginObject_struct) PleaseReRegisterClientAddress(ctx context.Context, reRegisterToGatewayMessage *gRPC.ReRegisterToGatewayMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse
	var parentgRPCAddress parentgRPCAddress_struct

	logger.WithFields(logrus.Fields{
		"ID":                         "b9e6bde2-0a59-4459-83c4-d723d50a080c",
		"reRegisterToGatewayMessage": reRegisterToGatewayMessage,
	}).Info("Incoming gRPC: 'PleaseReRegisterClientAddress'")

	//Move data into object that should be save in DB
	parentgRPCAddress.parentGatewayId = reRegisterToGatewayMessage.GatewayId
	parentgRPCAddress.parentGatewayName = reRegisterToGatewayMessage.GatewayName
	parentgRPCAddress.parentGatewayServer_address = reRegisterToGatewayMessage.GatewayAddress
	parentgRPCAddress.parentGatewayServer_port = reRegisterToGatewayMessage.GatewayPort
	parentgRPCAddress.createdDateTime = generaTimeStampUTC()

	// Convert Parent Gateway address info-struct into a byte array
	parentgRPCAddressByteArray, err := json.Marshal(reRegisterToGatewayMessage)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":                "6df384d3-ebe9-4a65-947f-e51814c8544c",
			"parentgRPCAddress": parentgRPCAddress,
			"err":               err,
		}).Error("Error when converting 'parentgRPCAddress' into a byte array, stopping futher processing of Reregistration.")

		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:       gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:                 gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:               gatewayConfig.gatewayIdentification.gatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Error when converting 'parentgRPCAddress' into a byte array, stopping futher processing of Reregistration.",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.systemDomain.gatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.systemDomain.gatewayDomainName,
		}

		returnMessage.Comments = "Error when converting 'parentgRPCAddress' into a byte array, stopping futher processing of Reregistration."
		returnMessage.Acknack = false
		return returnMessage, nil

	}

	// Save parentgRPCAddressByteArray to local database, using local channel
	// Return Channel
	returnChannel := make(chan dbResultMessage_struct)

	dbMessage := dbMessage_struct{
		DB_WRITE,
		BUCKET_PARENT_ADDRESS,
		BUCKET_KEY_PARENT_ADDRESS, // Key allways hardcoded due to one gateway or plugin can only have one parent
		parentgRPCAddressByteArray,
		returnChannel}

	// Send message to Database
	dbMessageQueue <- dbMessage

	// Wait for result on result channel then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "446cb366-9d84-4c1c-a628-ef162f7c1747",
			"err": err,
		}).Error("Got an error when Saveing to local DB")

		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:       gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:                 gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:               gatewayConfig.gatewayIdentification.gatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Got an error when Saveing to local DB",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.systemDomain.gatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.systemDomain.gatewayDomainName,
		}

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Error when saving to local DB, stopping futher processing of Reregistration"
		returnMessage.Acknack = false
		return returnMessage, nil
	}

	logger.WithFields(logrus.Fields{
		"ID": "fdf7081a-e7da-4bf1-a87c-82c51b8f575b",
	}).Debug("Reregistration info was saved in local database")

	// Start Registration of this Gateway/Plugin
	// **** TODO Put message on channel to register this client ****
	registerSuccess, err := registerThisGatewayAtParentGateway()

	if registerSuccess == true {
		logger.WithFields(logrus.Fields{
			"ID": "14ae3650-c0dd-4e23-b197-9706d5dfc8bd",
		}).Debug("Rereregistration to parent gateway/Fenix was successful")

		// Create message back to parent Gateway/Fenix
		returnMessage.Comments = "Reregistration was done"
		returnMessage.Acknack = true
		return returnMessage, nil

	} else {
		logger.WithFields(logrus.Fields{
			"ID":  "e5c273a8-f257-4c9f-bcbe-1697a2de1663",
			"err": err,
		}).Error("Rereregistration to parent gateway/Fenix could not be done")

		returnMessage.Comments = "Reregistration could not be done"
		returnMessage.Acknack = false
		return returnMessage, nil
	}

}
