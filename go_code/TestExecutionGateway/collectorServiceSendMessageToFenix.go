package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix for incoming request for supported TestDomains from a certain Plugin
//
func (gatewayObject *GatewayTowardsPluginObject_struct) SendMessageToFenix(ctx context.Context, informationMessage gRPC.InformationMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	gatewayObject.logger.WithFields(logrus.Fields{
		"ID":                 "25d5b2fd-cfb6-4cb8-9c7c-898e92c50e51",
		"informationMessage": informationMessage,
	}).Debug("Incoming gRPC: 'SendMessageToFenix'")

	// Put supportedTestDataDomainsRequest on queue for further processing
	gatewayObject.informationMessageChannel <- informationMessage
	gatewayObject.logger.WithFields(logrus.Fields{
		"ID": "5009cdce-c2a7-4a33-b0fa-e6f55cd805f8",
	}).Debug("InformationMessage was put on the channel")

	gatewayObject.logger.WithFields(logrus.Fields{
		"ID": "becb4e61-05e3-4395-b852-2c26377d8388",
	}).Debug("Leaving gRPC: 'SendMessageToFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "InformationMessage was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
