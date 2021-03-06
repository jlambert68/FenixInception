package TestExecutionGateway

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a InfoMessage toward Fenix
//
func (gRPCServerTowardsFenix *gRPCServerTowardsFenixStruct) SendTestInstructionTimeOutTowardsFenix(ctx context.Context, testInstructionTimeOutMessage *gRPC.TestInstructionTimeOutMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                            "7f3d4cf0-e4fa-4aae-a5da-990adcdaab22",
		"testInstructionTimeOutMessage": testInstructionTimeOutMessage,
	}).Debug("Incoming gRPC: 'SendTestInstructionTimeOutTowardsFenix'")

	// Put testInstructionTimeOutMessage on queue for further processing
	gatewayChannelPackage.TestInstructionTimeOutMessageChannelTowardsFenix <- testInstructionTimeOutMessage
	logger.WithFields(logrus.Fields{
		"ID": "c1845b2f-08fb-4eed-8c07-543c9fe24d45",
	}).Debug("'testInstructionTimeOutMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "b7a461b1-cd3b-47c7-a6b7-585e774c6642",
	}).Debug("Leaving gRPC: 'SendTestInstructionTimeOutTowardsFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "'TestInstructionTimeOutMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
