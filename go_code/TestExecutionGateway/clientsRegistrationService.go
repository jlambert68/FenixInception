package TestExecutionGateway

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from Client Gateway/Plugin with an registration of itself
//
func (gRPCServerTowardsFenixStruct *gRPCServerTowardsFenixStruct) RegisterClientAddress(ctx context.Context, registerClientAddressRequest *gRPC.RegisterClientAddressRequest) (*gRPC.RegisterClientAddressResponse, error) {

	var registerClientAddressResponse *gRPC.RegisterClientAddressResponse
	var clientRPCAddress common_code.ClientsAddressAndPortStruct

	logger.WithFields(logrus.Fields{
		"ID":                           "b88ad310-944f-44ed-bfca-6b06337b81be",
		"registerClientAddressRequest": registerClientAddressRequest,
	}).Info("Incoming gRPC: 'RegisterClientAddress'")

	// Check if calling client is using an old gRPC-version-defitnition file (wrong version)
	// TODO Denna jämförelse är troligen som Äpplen och Päron
	if registerClientAddressRequest.GRPCVersion.String() != common_code.GetHighestGRPCVersion() {

		// Convert time into proto-time format
		protoTimeStamp, _ := ptypes.TimestampProto(common_code.GeneraTimeStampUTC())

		// Send Error information to Fenix
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         registerClientAddressRequest.CallingSystemId,
			OriginalSenderName:       registerClientAddressRequest.CallingSystemName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Child gateway/Plugin is using wrong version of gRPC-defition",
			OrginalCreateDateTime:    protoTimeStamp,
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
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
	clientRPCAddress.ClientId = registerClientAddressRequest.CallingSystemId
	clientRPCAddress.ClientName = registerClientAddressRequest.CallingSystemName
	clientRPCAddress.ClientAddress = registerClientAddressRequest.CallingSystemIpAddress
	clientRPCAddress.ClientPort = nextPort

	// Convert Client Gateway/Plugin address info-struct into a byte array
	childgRPCAddressByteArray, err := json.Marshal(clientRPCAddress)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":                        "f0b16a20-649f-4b69-97f2-649bc686a5a3",
			"childgRPCAddressByteArray": childgRPCAddressByteArray,
			"err":                       err,
		}).Error("Error when converting 'childgRPCAddressByteArray' into a byte array, stopping futher processing of RegisterClientAddress.")

		// Convert time into proto-time format
		protoTimeStamp, _ := ptypes.TimestampProto(common_code.GeneraTimeStampUTC())

		// Send Error information to Fenix
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Error when converting 'clientRPCAddress' into a byte array, stopping futher processing of RegisterClientAddress.",
			OrginalCreateDateTime:    protoTimeStamp,
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
		}

		registerClientAddressResponse = &gRPC.RegisterClientAddressResponse{
			ClientPort: 0,
			Acknack:    false,
			Comments:   "Error when converting 'clientRPCAddress' into a byte array, stopping futher processing of RegisterClientAddress..",
		}

		return registerClientAddressResponse, nil

	}

	// Save childgRPCAddressByteArray to local database, using local channel
	// Return Channel
	returnChannel := make(chan common_code.DbResultMessageStruct)

	dbMessage := common_code.DbMessageStruct{
		common_code.DbWrite,
		common_code.BucketForClients,
		clientRPCAddress.ClientId,
		childgRPCAddressByteArray,
		returnChannel}

	// Send message to Database
	dbMessageQueue <- dbMessage

	// Wait for result on result channel then close returnChannel
	returnDBMessage := <-returnChannel
	close(returnChannel)

	if returnDBMessage.Err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "446cb366-9d84-4c1c-a628-ef162f7c1747",
			"err": err,
		}).Error("Got an error when Saveing to local DB")

		// Convert timestamp into proto-format
		protoTimeStamp, _ := ptypes.TimestampProto(common_code.GeneraTimeStampUTC())

		// Send Error information to Fenix
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- &gRPC.InformationMessage{
			OriginalSenderId:         gatewayConfig.GatewayIdentification.GatewayId,
			OriginalSenderName:       gatewayConfig.GatewayIdentification.GatewayName,
			SenderId:                 gatewayConfig.GatewayIdentification.GatewayId,
			SenderName:               gatewayConfig.GatewayIdentification.GatewayName,
			MessageId:                common_code.GenerateUUID(logger),
			MessageType:              gRPC.InformationMessage_ERROR,
			Message:                  "Got an error when Saveing to local DB",
			OrginalCreateDateTime:    protoTimeStamp,
			OriginalSystemDomainId:   gatewayConfig.SystemDomain.GatewayDomainId,
			OriginalSystemDomainName: gatewayConfig.SystemDomain.GatewayDomainName,
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
