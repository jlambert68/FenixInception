package TestExecutionGateway

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a TestExecutionLogMessage toward Fenix
//
func (gRPCServerTowardsFenix *gRPCServerTowardsFenixStruct) SendTestExecutionLogTowardsFenix(ctx context.Context, testExecutionLogMessage *gRPC.TestExecutionLogMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                            "6ef6943e-7a05-42cd-80b5-e35231c14f04",
		"testInstructionTimeOutMessage": testExecutionLogMessage,
	}).Debug("Incoming gRPC: 'TestExecutionLogMessage'")

	// Put testExecutionLogMessage on queue for further processing
	gatewayChannelPackage.TestExecutionLogMessageChannelTowardsFenix <- testExecutionLogMessage
	logger.WithFields(logrus.Fields{
		"ID": "bbd8fdc7-16cb-4dae-b4a1-821fb30ccb70",
	}).Debug("'testExecutionLogMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "f07b61d0-16cc-4b9b-84a3-2e6c642125c3",
	}).Debug("Leaving gRPC: 'SendTexecutionLogTowardsFenix'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "'TestExecutionLogMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
