package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix for incoming request for supported TestDomains from a certain Plugin
//
func (gRPCServerTowardsPlugin *GRPCServerTowardsPluginStruct) GetSupportedTestDataDomains(ctx context.Context, supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID":                              "97ba5adb-248f-4186-9d7c-2e457af826a7",
		"supportedTestDataDomainsRequest": *supportedTestDataDomainsRequest,
	}).Debug("Incoming gRPC: 'GetSupportedTestDataDomains'")

	// Put supportedTestDataDomainsRequest on queue for further processing
	gatewayTowardsPluginObject.supportedTestDataDomainsRequestChannelTowardsPlugin <- supportedTestDataDomainsRequest
	logger.WithFields(logrus.Fields{
		"ID": "5009cdce-c2a7-4a33-b0fa-e6f55cd805f8",
	}).Debug("SupportedTestDataDomainsRequest was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "ece8cd79-bdd2-4ab2-8500-f52ca0a5c226",
	}).Debug("Leaving gRPC: 'GetSupportedTestDataDomains'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "SupportedTestDataDomainsRequest was forwarded towards plugin"
	returnMessage.Acknack = true

	return returnMessage, nil
}
