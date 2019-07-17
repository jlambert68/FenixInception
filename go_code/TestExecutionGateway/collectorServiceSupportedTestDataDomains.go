package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// ********************************************************************************************
// Call from this(bufferd in DB)/child Gateway/Plugin for incoming request for forwarding a supportedTestDataDomainsWithHeadersMessage toward Fenix
//
func (gRPCServerTowardsFenix *GRPCServerTowardsFenixStruct) SupportedTestDataDomains(ctx context.Context, supportedTestDataDomainsWithHeadersMessage *gRPC.SupportedTestDataDomainsWithHeadersMessage) (*gRPC.AckNackResponse, error) {

	var returnMessage *gRPC.AckNackResponse

	logger.WithFields(logrus.Fields{
		"ID": "7572dcf4-687d-469a-b1a1-3e0b884c5766",
		"supportedTestDataDomainsWithHeadersMessage": supportedTestDataDomainsWithHeadersMessage,
	}).Debug("Incoming gRPC: 'SupportedTestDataDomains'")

	// Put supportedTestDataDomainsWithHeadersMessage on queue for further processing
	supportedTestDataDomainsWithHeadersMessageTowardsFenixChannelTowardsFenix <- supportedTestDataDomainsWithHeadersMessage
	logger.WithFields(logrus.Fields{
		"ID": "43fef475-85a7-4b29-926b-e40dffa81503",
	}).Debug("'supportedTestDataDomainsWithHeadersMessage' was put on the channel")

	logger.WithFields(logrus.Fields{
		"ID": "5a9d4587-59a4-40b7-a2c6-586ecc7524fb",
	}).Debug("Leaving gRPC: 'SupportedTestDataDomains'")

	// Create message back to child Gateway/Plugin
	returnMessage.Comments = "'supportedTestDataDomainsWithHeadersMessage' was forwarded towards Fenix"
	returnMessage.Acknack = true

	return returnMessage, nil
}
