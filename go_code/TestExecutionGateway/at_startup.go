package TestExecutionGateway

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

// *******************************************************************
// Start up this gateway by register it to parent gateway/Fenix
//
func (gatewayObject *gatewayTowardsFenixObject_struct) tryToRegisterGatewayAtParent() {

	// Initiate map used for handle CLients address and port info
	clientsAddressAndPort = make(map[string]clientsAddressAndPort_struct)

	// Register gateway/client at parent Gateway/Fenix
	resultBool, err := gatewayObject.registerThisGatewayAtParentGateway()
	if err != nil || resultBool == false {
		// If this gateway never has been connected to parent gateway/Fenix then Exit
		// due to that parent doesn't know this gateways address yet
		if gatewayClientHasBeenConnectedToParentGateway == false {
			logger.WithFields(logrus.Fields{
				"ID": "c7ea051f-37b2-41d2-820e-5050a560cfbb",
			}).Fatal("This gateway has never been connected to parent gateway/Fenix so Exit, becasue Parent Gateway/Fenix doesn't know the address to this gaeway")
		} else {
			logger.WithFields(logrus.Fields{
				"ID": "35b7981d-ed97-48bf-8f5e-9807da4cced4",
			}).Warning("Parent Gateway/Fenix is not alive so Waiting for Gateway/Fenix to reconnect")
		}
	}

	//  end message to all known clients that they must reRegister themself to this sgateway server
	// Run goroutine and use queue to count down registrations

	//  Regiater this gateway client to parent gateway
}

// TODO `json:"page"` fixa detta för de objekt som ska sparas i localDB

// *******************************************************************
// Register this gateway/client at parent gateway/Fenix
//
func (gatewayObject *gatewayTowardsFenixObject_struct) registerThisGatewayAtParentGateway() (bool, error) {

	var err error
	var addressToDial string

	// Find parents address and port to call
	addressToDial = gatewayConfig.parentgRPCAddress.parentGatewayServer_address + gatewayConfig.parentgRPCAddress.parentGatewayServer_port

	// Information sent to parent gateway/Fenix
	registerClientAddressRequest := gRPC.RegisterClientAddressRequest{
		GRPCVersion:            gRPC.CurrentVersionEnum_VERSION_0_1_0,
		CallingSystemId:        gatewayConfig.gatewayIdentification.callingSystemId,
		CallingSystemName:      gatewayConfig.gatewayIdentification.callingSystemName,
		CallingSystemIpAddress: gatewayConfig.gatewayIdentification.callingSystemIpAddress}

	// If no port received from parent Gateway/Fenix then port should be cleared
	incomingPortForCallsFromParentGateway = ""

	// Set up connection to Parent Gateway/Fenix Server
	remoteGatewayServerConnection, err = grpc.Dial(addressToDial, grpc.WithInsecure())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":            "fbd24ed4-638a-43ac-a07b-c622f0ab325c",
			"addressToDial": addressToDial,
			"error message": err,
		}).Warning("Did not connect to Parent Gateway/Fenix Server!")
		return false, err

	} else {
		logger.WithFields(logrus.Fields{
			"ID":            "14d029db-0031-4837-b139-7b04b707fabf",
			"addressToDial": addressToDial,
		}).Debug("gRPC connection OK to Worker Server!")

		// Creates a new Gateway Client
		gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteGatewayServerConnection)

		ctx := context.Background()
		registerClientAddressResponse, err := gatewayClient.RegisterClientAddress(ctx, &registerClientAddressRequest)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"ID":            "364e2a90-1c8b-47df-be64-b73457317911",
				"returnMessage": registerClientAddressResponse,
				"error":         err,
			}).Error("Problem to register gateway to parent Gateway/Fenix")

			return false, err

		} else {
			// Check answer from parent Gateway/Fenix
			if registerClientAddressResponse.Acknack == false {
				logger.WithFields(logrus.Fields{
					"ID":            "236cf5a9-c038-4a2a-b9bc-24015becec18",
					"returnMessage": registerClientAddressResponse,
					"error":         err,
				}).Warning("Parent Gateway/Fenix did not allow registration")

				return false, err

			} else {
				logger.WithFields(logrus.Fields{
					"ID":            "116024c5-268b-4688-97ca-272ab3db385f",
					"returnMessage": registerClientAddressResponse,
					"error":         err,
				}).Debug("Success in registering gateway to parent Gateway/Fenix")
			}
		}

		// Take care of port-information sent back from Parent Gateway/Fenix
		incomingPortForCallsFromParentGateway = registerClientAddressResponse.ClientPort
		// TODO Save port in DB for use if gateway restarts

	}
	return true, nil

}

// *******************************************************************
// Start all Services
//
func startAllServices() {

	// Ensure that all services don't start before everything has been started
	gatewayMustStopProcessing = true

	// Init logger
	initLogger("")

	// Read 'gatewayConfig.toml' for config parameters
	processConfigFile()

	// Initiate Database

	// Start 'dispatchEngineForTestInstructions'
	gatewayTowardsPluginObject.initiateDispatchEngineForTestInstructiona()

	// Start 'dispatchEngineForSupportedTestDomains'
	gatewayTowardsPluginObject.initiateDispatchEngineForSupportedTestDomains()

	// Start 'transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix'
	gatewayTowardsFenixObject.initiateRegistrateAvailableTestDataDomainsTowardsFenix()

	// Start 'transmitEngineForSendMessageToFenix'
	gatewayTowardsFenixObject.initiateSendMessageToFenix()

	// Start 'transmitEngineForSendTestExecutionLogTowardsFenix'
	gatewayTowardsFenixObject.initiateSendTestExecutionLogTowardsFenix()

	// Start 'transmitEngineForSendTestInstructionTimeOutTowardsFenix'
	gatewayTowardsFenixObject.initiateSendTestInstructionTimeOutTowardsFenix()

	// Start all services at the same time
	gatewayMustStopProcessing = false

	// Try to Register this Gateway At Parent
	gatewayTowardsFenixObject.tryToRegisterGatewayAtParent()

}
