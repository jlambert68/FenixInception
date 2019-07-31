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
func (gRPCServerTowardsFenix *common_code.GRPCServerTowardsFenixStruct) RegistrateAailableTestContainers(ctx context.Context, availbleTestContainersAtPluginMessage *gRPC.AvailbleTestContainersAtPluginMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	common_code.logger.WithFields(logrus.Fields{
		"ID":                                    "ead1b630-27cd-4c1b-9e55-5bd28198a9f8",
		"availbleTestContainersAtPluginMessage": availbleTestContainersAtPluginMessage,
	}).Debug("Incoming gRPC: 'RegistrateAailableTestContainers'")

	// Put testExecutionLogMessage on queue for further processing
	common_code.availbleTestContainersAtPluginMessageTowardsFenixChannelTowardsFenix <- availbleTestContainersAtPluginMessage
	common_code.logger.WithFields(logrus.Fields{
		"ID": "14c40c76-0dea-4279-84cf-ed78da51aa3e",
	}).Debug("'SupportedTestDataDomainsMessage' was put on the channel")

	common_code.logger.WithFields(logrus.Fields{
		"ID": "b658c4b5-a47e-4ec5-908c-134973d6ec95",
	}).Debug("Leaving gRPC: 'RegistrateAailableTestContainers'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'availbleTestContainersAtPluginMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
