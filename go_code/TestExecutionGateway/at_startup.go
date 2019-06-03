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

	//  TODO end message to all known clients that they must reRegister themself to this sgateway server
	// Run goroutine and use queue to count down registrations

	//  Regiater this gateway client to parent gateway
}

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
		CallingSystemId:        gatewayConfig.gatewayIdentification.gatewayId,
		CallingSystemName:      gatewayConfig.gatewayIdentification.gatewayName,
		CallingSystemIpAddress: gatewayConfig.gatewayIdentification.gatewayIpAddress}

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

		// Save Port to memory object
		gatewayConfig.parentgRPCAddress.parentGatewayServer_port = registerClientAddressResponse.ClientPort

		// Update address info in local database
		updateDatabaseFromMemoryForParentAddressInfo(gRPC.ReRegisterToGatewayMessage{
			GatewayAddress:  gatewayConfig.parentgRPCAddress.parentGatewayServer_address,
			GatewayPort:     gatewayConfig.parentgRPCAddress.parentGatewayServer_port,
			GatewayId:       gatewayConfig.parentgRPCAddress.parentGatewayId,
			GatewayName:     gatewayConfig.parentgRPCAddress.parentGatewayName,
			CreatedDateTime: generaTimeStampUTC(),
		})

	}
	return true, nil

}

// *******************************************************************
// Cleanup before closing down
//
var cleanupProcessed bool = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		db.Close()
		logger.WithFields(logrus.Fields{
			"ID": "4f2f77ba-d105-47bc-8120-6b2874faa98d",
		}).Info("Closing local database")

		// Close gRPC connection to Parent gateway/Fenix
		registerGatewayTowardsFenixServer.GracefulStop()

		logger.WithFields(logrus.Fields{
			"ID": "ed85f1c6-98d9-4c08-bf52-2364a42a5bad",
		}).Info("Gracefull stop for: 'registerGatewayTowardsFenixServer'")

		// Close gRPC connection to Children
		registerGatewayTowardsPluginerver.GracefulStop()

		logger.WithFields(logrus.Fields{
			"ID": "a7ee0eea-3bf2-44b7-a5e1-942a15d0d7fd",
		}).Info("Gracefull stop for: 'registerGatewayTowardsPluginerver'")

		// Stop listening towards parent gateway/Fenix
		gatewayTowardsFenixListener.Close()
		logger.WithFields(logrus.Fields{
			"ID": "6fe1108d-31af-4db0-a70a-6bf5fe02b6f8",
		}).Info("Stop listening, from Parent, on port: " + gatewayConfig.gatewayIdentification.gatewaParentCallOnThisPort)

		// Stop listening towards clients
		gatewayTowardsFenixListener.Close()
		logger.WithFields(logrus.Fields{
			"ID": "35b35f32-7e5b-420c-8544-072e868e5bbb",
		}).Info("Stop listening, from Children, on port: " + gatewayConfig.gatewayIdentification.gatewayChildrenCallOnThisPort)

	}
}

// *******************************************************************
// Start all Services
//
func startAllServices() {

	// Cleanup all gRPC connections
	defer cleanup()

	// Ensure that all services don't start before everything has been started
	gatewayMustStopProcessing = true

	// Read 'gatewayConfig.toml' for config parameters
	processConfigFile()

	// Init logger
	initLogger("localLogFile")

	// Initiate Database
	initiateDB("") // Use default database file name

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

	// Update Memory information about parent address and port with that saved in database, database overrule config-file
	updateMemoryAddressForParentAddressInfo()

	// Try to Register this Gateway At Parent
	gatewayTowardsFenixObject.tryToRegisterGatewayAtParent()

	// Ask clients to ReRegister them self to this gateway
	// TODO Make all Clients ReRegister them self
}
