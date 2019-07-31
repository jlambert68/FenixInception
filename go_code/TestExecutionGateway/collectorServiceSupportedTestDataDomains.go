package TestExecutionGateway

import (
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a supportedTestDataDomainsWithHeadersMessage toward Fenix
//
func (gRPCServerTowardsFenix *common_code.GRPCServerTowardsFenixStruct) SupportedTestDataDomains(ctx context.Context, supportedTestDataDomainsWithHeadersMessage *gRPC.SupportedTestDataDomainsWithHeadersMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	common_code.Logger.WithFields(logrus.Fields{
		"ID": "7572dcf4-687d-469a-b1a1-3e0b884c5766",
		"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
	}).Debug("Incoming gRPC: 'SupportedTestDataDomains'")

	// Put supportedTestDataDomainsWithHeadersMessage on queue for further processing
	common_code.supportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsWithHeadersMessage
	common_code.Logger.WithFields(logrus.Fields{
		"ID": "43fef475-85a7-4b29-926b-e40dffa81503",
	}).Debug("'supportedTestDataDomainsWithHeadersMessage' was put on the channel")

	common_code.Logger.WithFields(logrus.Fields{
		"ID": "5a9d4587-59a4-40b7-a2c6-586ecc7524fb",
	}).Debug("Leaving gRPC: 'SupportedTestDataDomains'")

	// Create message back to child Gateway/Plugin
	returnMessage.Comments = "'supportedTestDataDomainsWithHeadersMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
