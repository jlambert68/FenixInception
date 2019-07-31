package TestExecutionGateway

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from parent Gateway/Fenix for incoming request for supported TestDomains from a certain Plugin
//
func (gRPCServerTowardsPlugin *common_code.GRPCServerTowardsPluginStruct) GetSupportedTestDataDomains(ctx context.Context, supportedTestDataDomainsRequest *gRPC.SupportedTestDataDomainsRequest) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	common_code.Logger.WithFields(logrus.Fields{
		"ID":                              "97ba5adb-248f-4186-9d7c-2e457af826a7",
		"supportedTestDataDomainsRequest": *supportedTestDataDomainsRequest,
	}).Debug("Incoming gRPC: 'GetSupportedTestDataDomains'")

	// Put supportedTestDataDomainsRequest on queue for further processing
	common_code.supportedTestDataDomainsRequestChannelTowardsPlugin <- supportedTestDataDomainsRequest
	common_code.Logger.WithFields(logrus.Fields{
		"ID": "5009cdce-c2a7-4a33-b0fa-e6f55cd805f8",
	}).Debug("SupportedTestDataDomainsRequest was put on the channel")

	common_code.Logger.WithFields(logrus.Fields{
		"ID": "ece8cd79-bdd2-4ab2-8500-f52ca0a5c226",
	}).Debug("Leaving gRPC: 'GetSupportedTestDataDomains'")

	// Create message back to parent Gateway/Fenix
	returnMessage.Comments = "SupportedTestDataDomainsRequest was forwarded towards plugin"
	returnMessage.Acknack = true

	return returnMessage, nil
}
