package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"github.com/jlambert68/FenixInception/go_code/common_code"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func (gRPCServerTowardsFenix *common_code.GRPCServerTowardsFenixStruct) RegisterAvailbleTestInstructions(ctx context.Context, availbleTestInstructionAtPluginMessage *gRPC.AvailbleTestInstructionAtPluginMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	common_code.logger.WithFields(logrus.Fields{
		"ID":                                     "a8adaf82-88a7-4e95-8aa2-6e0b6fdc2727",
		"availbleTestInstructionAtPluginMessage": availbleTestInstructionAtPluginMessage,
	}).Debug("Incoming gRPC: 'RegisterAvailbleTestInstructions'")

	// Put availbleTestInstructionAtPluginMessage on queue for further processing
	common_code.availbleTestInstructionAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestInstructionAtPluginMessage
	common_code.logger.WithFields(logrus.Fields{
		"ID": "e6ed9172-39ad-4a7c-bb07-b76f2555b36a",
	}).Debug("'availbleTestInstructionAtPluginMessage' was put on the channel")

	common_code.logger.WithFields(logrus.Fields{
		"ID": "919fd3a0-7d8a-4393-a668-9db565c21ea8",
	}).Debug("Leaving gRPC: 'RegisterAvailbleTestInstructions'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'availbleTestInstructionAtPluginMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
