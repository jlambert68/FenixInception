package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from Client Gateway/Plugin with an registration of itself
//
func (gatewayObject *gatewayTowardsFenixObjectStruct) RegisterClientAddress(ctx context.Context, registerClientAddressRequest *gRPC.RegisterClientAddressRequest) (*gRPC.RegisterClientAddressResponse, error) {

	var registerClientAddressResponse *gRPC.RegisterClientAddressResponse
	var clientRPCAddress clientsAddressAndPortStruct

	logger.WithFields(logrus.Fields{
		"ID":                           "b88ad310-944f-44ed-bfca-6b06337b81be",
		"registerClientAddressRequest": registerClientAddressRequest,
	}).Info("Incoming gRPC: 'RegisterClientAddress'")

	// Check if calling client is using an old gRPC-version-defitnition file (wrong version)
	// TODO Denna jämförelse är troligen som Äpplen och Päron
	if registerClientAddressRequest.GRPCVersion.String() != getHighestGRPCVersion() {
		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         registerClientAddressRequest.CallingSystemId,
			OriginalSenderName:       registerClientAddressRequest.CallingSystemName,
			SenderId:                 gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:               gatewayConfig.gatewayIdentification.gatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Child gateway/Plugin is using wrong version of gRPC-defition",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.systemDomain.gatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.systemDomain.gatewayDomainName,
		}

		registerClientAddressResponse = &gRPC.RegisterClientAddressResponse{
			ClientPort: 0,
			Acknack:    false,
			Comments:   "This gateway/Plugin is using wrong version of gRPC-defition",
		}

		return registerClientAddressResponse, nil
	}

	// *** Save Clents Address Info in local database ***
	nextPort := getNextFreeClientPort(registerClientAddressRequest.CallingSystemIpAddress)

	//Move data into object that should be save in DB
	clientRPCAddress.clientId = registerClientAddressRequest.CallingSystemId
	clientRPCAddress.clientName = registerClientAddressRequest.CallingSystemName
	clientRPCAddress.clientAddress = registerClientAddressRequest.CallingSystemIpAddress
	clientRPCAddress.clientPort = nextPort

	// Convert Client Gateway/Plugin address info-struct into a byte array
	childgRPCAddressByteArray, err := json.Marshal(clientRPCAddress)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":                        "f0b16a20-649f-4b69-97f2-649bc686a5a3",
			"childgRPCAddressByteArray": childgRPCAddressByteArray,
			"err":                       err,
		}).Error("Error when converting 'childgRPCAddressByteArray' into a byte array, stopping futher processing of RegisterClientAddress.")

		// Send Error information to Fenix
		localInformationMessageChannel <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.gatewayIdentification.gatewayId,
			OriginalSenderName:       gatewayConfig.gatewayIdentification.gatewayName,
			SenderId:                 gatewayConfig.gatewayIdentification.gatewayId,
			SenderName:               gatewayConfig.gatewayIdentification.gatewayName,
			MessageId:                generateUUID(),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Error when converting 'parentgRPCAddress' into a byte array, stopping futher processing of RegisterClientAddress.",
			OrginalCreateDateTime:    generaTimeStampUTC(),
			OriginalSystemDomainId:   gatewayConfig.systemDomain.gatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.systemDomain.gatewayDomainName,
		}

		registerClientAddressResponse = &gRPC.RegisterClientAddressResponse{
			ClientPort: 0,
			Acknack:    false,
			Comments:   "Error when converting 'parentgRPCAddress' into a byte array, stopping futher processing of RegisterClientAddress..",
		}

		return registerClientAddressResponse, nil

	}

	// Save childgRPCAddressByteArray to local database, using local channel
	// Return Channel
	returnChannel := make(chan dbResultMessageStruct)

	dbMessage := dbMessageStruct{
		DbWrite,
		BucketForParentAddress,
		BucketKeyForParentAddress, // Key allways hardcoded due to one gateway or plugin can only have one parent
		childgRPCAddressByteArray,
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

		// Create message back to child Gateway/Plugin
		registerClientAddressResponse = &gRPC.RegisterClientAddressResponse{
			ClientPort: 0,
			Acknack:    false,
			Comments:   "Got an error when Saveing to local DB",
		}

		return registerClientAddressResponse, nil

	} else {

		// Return port for client to listen to
		logger.WithFields(logrus.Fields{
			"ID": "b6be47f6-c6ee-44c8-b21c-d64891e70038",
		}).Debug("Registration of client was done")

		// Create message back to child Gateway/Plugin
		registerClientAddressResponse = &gRPC.RegisterClientAddressResponse{
			ClientPort: nextPort,
			Acknack:    true,
			Comments:   "Success in doing registration, please listen on attached port",
		}

		return registerClientAddressResponse, nil
	}

}
