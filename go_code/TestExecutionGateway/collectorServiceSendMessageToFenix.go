package TestExecutionGateway

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a InfoMessage toward Fenix
//
func (gRPCServerTowardsFenix *gRPCServerTowardsFenixStruct) SendMessageToFenix(ctx context.Context, informationMessage *gRPC.InformationMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                 "25d5b2fd-cfb6-4cb8-9c7c-898e92c50e51",
		"informationMessage": informationMessage,
	}).Debug("Incoming gRPC: 'SendMessageToFenix'")

	// Put supportedTestDataDomainsRequest on queue for further processing
	gatewayChannelPackage.InformationMessageChannelTowardsFenix <- informationMessage
	logger.WithFields(logrus.Fields{
		"ID": "5009cdce-c2a7-4a33-b0fa-e6f55cd805f8",
	}).Debug("InformationMessage was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "becb4e61-05e3-4395-b852-2c26377d8388",
	}).Debug("Leaving gRPC: 'SendMessageToFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "InformationMessage was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}

// ********************************************************************************************
// Initiate xxxx Engine for SupportedTestDataDomainsMessage-messages
//

func initiateSendLocalMessageToFenixEngine() {

	// Start Transmit Engine, for SupportedTestDataDomainsMessage-messages as a go-routine
	logger.WithFields(logrus.Fields{
		"ID": "85e5448f-0d6d-4f72-be4b-70e51fc3d362",
	}).Info("Initiate: 'transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix'")

	go sendLocalMessageToFenixEngine()
}

// ********************************************************************************************
// Local call from this Gateway, when local errors/warnings/problems for sending a InfoMessage towards Fenix
//
func sendLocalMessageToFenixEngine() {

	for {
		// Wait for data comes from 'local' channel for messages initiated in this gateway
		localInformationMessage := <-gatewayChannelPackage.InformationMessageChannelTowardsFenix

		logger.WithFields(logrus.Fields{
			"ID":                      "f74b3a27-77ad-47a4-8e77-ed680f4c419f",
			"localInformationMessage": localInformationMessage,
		}).Debug("Received a new 'localInformationMessage' from local channel that shoud be forwarded")

		// Put localInformationMessage on queue for further processing
		gatewayChannelPackage.InformationMessageChannelTowardsFenix <- localInformationMessage
		logger.WithFields(logrus.Fields{
			"ID": "fbc22674-d6ec-4885-8c38-89a049abafee",
		}).Debug("self created 'MessageToFenix' was put on the channel")

		logger.WithFields(logrus.Fields{
			"ID": "521cc222-88f7-4e52-9315-0a891fc06111",
		}).Debug("Leaving Send self created 'MessageToFenix'")

	}
	return
}
