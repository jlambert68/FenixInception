package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a testInstructionExecutionResultMessage toward Fenix
//
func (gRPCServerTowardsFenix *GRPCServerTowardsFenixStruct) SendTestInstructionResultTowardsFenix(ctx context.Context, testInstructionExecutionResultMessage *gRPC.TestInstructionExecutionResultMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                                    "65959397-cd47-4c2c-a051-76338d44aeb7",
		"testInstructionExecutionResultMessage": testInstructionExecutionResultMessage,
	}).Debug("Incoming gRPC: 'SendTestInstructionResultTowardsFenix'")

	// Put testInstructionExecutionResultMessage on queue for further processing
	testInstructionExecutionResultMessageTowardsFenixChannelTowardsFenix <- testInstructionExecutionResultMessage
	logger.WithFields(logrus.Fields{
		"ID": "0a7d5580-609f-4a74-8b9c-1c06ee8bd771",
	}).Debug("'SupportedTestDataDomainsMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "16766064-6c77-4f3b-8a31-0dedc4303b7c",
	}).Debug("Leaving gRPC: 'SendTestInstructionResultTowardsFenix'")

	// Create message back to child Gateway/Plugin
	returnMessage.Comments = "'testInstructionExecutionResultMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
