package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a SupportedTestDataDomainsMessage toward Fenix
//
func (gRPCServerTowardsFenix *GRPCServerTowardsFenixStruct) RegistrateAvailableTestDataDomains(ctx context.Context, supportedTestDataDomainsMessage *gRPC.SupportedTestDataDomainsMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                              "d3672229-2897-4654-83aa-0e56dd1db4b1",
		"supportedTestDataDomainsMessage": supportedTestDataDomainsMessage,
	}).Debug("Incoming gRPC: 'RegistrateAvailableTestDataDomains'")

	// Put testExecutionLogMessage on queue for further processing
	supportedTestDataDomainsMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsMessage
	logger.WithFields(logrus.Fields{
		"ID": "4d55f441-4e13-41ae-a9eb-e71f288a18e9",
	}).Debug("'SupportedTestDataDomainsMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "fb223a31-88d9-4876-964f-a21d1cdb0a2c",
	}).Debug("Leaving gRPC: 'RegistrateAvailableTestDataDomains'")

	// Create message back to parent Gateway/Plugin
	returnMessage.Comments = "'SupportedTestDataDomainsMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
