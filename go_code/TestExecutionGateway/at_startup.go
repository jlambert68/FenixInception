package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
	"strconv"
)

// *******************************************************************
// Start up this gateway by register it to parent gateway/Fenix
//
func (gatewayObject *gatewayTowardsFenixObjectStruct) tryToRegisterGatewayAtParent() {

	// Check if this gateway i used in integration test mode and only should start without any connection to parent gateway
	if GatewayInIntegrationTestMode.StartWithOutAnyParent == false {
		// Register gateway/client at parent Gateway/Fenix
		resultBool, err := registerThisGatewayAtParentGateway()
		if err != nil || resultBool == false {
			// If this gateway never has been connected to parent gateway/Fenix then Exit
			// due to that parent doesn't know this gateways address yet
			if gatewayConfig.ParentgRPCAddress.ConnectionToParentDoneAtLeastOnce == false {
				if GatewayInIntegrationTestMode.IsInSelfIntegrationTestMode == false {
					// Gateway is NOT in IsInSelfIntegrationTestMode
					logger.WithFields(logrus.Fields{
						"ID": "c7ea051f-37b2-41d2-820e-5050a560cfbb",
					}).Fatal("This gateway has never been connected to parent gateway/Fenix so Exit, because Parent Gateway/Fenix doesn't know the address to this gateway")
				} else {
					// Gateway IS in IsInSelfIntegrationTestMode
					logger.WithFields(logrus.Fields{
						"ID": "580d2c7d-b8d3-40f7-b238-eb096d859355",
					}).Error("This gateway has never been connected to parent gateway/Fenix so Exit, because Parent Gateway/Fenix doesn't know the address to this gateway")
				}
			} else {
				logger.WithFields(logrus.Fields{
					"ID": "35b7981d-ed97-48bf-8f5e-9807da4cced4",
				}).Warning("Parent Gateway/Fenix is not alive so Waiting for Gateway/Fenix to reconnect")
			}
		} else {
			logger.WithFields(logrus.Fields{
				"ID": "9401d538-0d13-4213-bab8-d5e546784738",
			}).Debug("Success in connectiing to parent Gateway/Fenix ")
		}
	} else {
		logger.WithFields(logrus.Fields{
			"ID": "db3709f2-a848-44be-82a0-509cc4bc08db",
		}).Debug("No Connection done to parent Gateway/Fenix. This gateway is used for Integrations Tests ")
	}

}

// *******************************************************************
// Register this gateway/client at parent gateway/Fenix
//
func registerThisGatewayAtParentGateway() (bool, error) {

	var err error
	var addressToDial string

	// Find parents address and port to call
	addressToDial = gatewayConfig.ParentgRPCAddress.ParentGatewayServerAddress + ":" + strconv.FormatInt(int64(gatewayConfig.ParentgRPCAddress.ParentGatewayServerPort), 10)

	// Information sent to parent gateway/Fenix
	registerClientAddressRequest := gRPC.RegisterClientAddressRequest{
		GRPCVersion:            gRPC.CurrentVersionEnum_VERSION_0_1_0,
		CallingSystemId:        gatewayConfig.GatewayIdentification.GatewayId,
		CallingSystemName:      gatewayConfig.GatewayIdentification.GatewayName,
		CallingSystemIpAddress: gatewayConfig.GatewayIdentification.GatewayIpAddress}

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
		}).Debug("gRPC connection OK to Parent gateway")

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
		gatewayConfig.GatewayIdentification.GatewaParentCallOnThisPort = registerClientAddressResponse.ClientPort

		// Save information that about that registration was successful. Used for knowning that registration was made at least once
		gatewayConfig.ParentgRPCAddress.ConnectionToParentDoneAtLeastOnce = true
		gatewayConfig.ParentgRPCAddress.ConnectionToParentLastConnectionDateTime = generaTimeStampUTC()

		// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
		gatewayIdentificationByteArray, err := json.Marshal(gatewayConfig.GatewayIdentification)

		if err != nil {
			// Error when Marshaling to []byte
			LogErrorAndSendInfoToFenix(
				"a1a37a7e-a9f1-4374-8ccd-239bc4a63e08",
				gRPC.InformationMessage_FATAL,
				"gatewayIdentificationByteArray",
				"No data available...",
				err.Error(),
				"Error when converting 'gatewayConfig.GatewayIdentification' into a byte array, stopping futher processing.",
			)
		} else {
			// Marshaling to []byte OK

			// Save gateWayIdentifaction information to local DB
			_ = SaveMessageToLocalDB(
				BucketKeyForGatewayIdentificationInfo,
				gatewayIdentificationByteArray,
				BucketForGatewayIdentificationInfo,
				"ddc66d46-b9b9-45f4-8da8-75acdb17b8be",
			)
		}

		// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
		parentgRPCAddressByteArray, err := json.Marshal(gatewayConfig.ParentgRPCAddress)

		if err != nil {
			// Error when Marshaling to []byte
			LogErrorAndSendInfoToFenix(
				"322ecb3d-7c10-4638-a9ec-242c695efc0e",
				gRPC.InformationMessage_FATAL,
				"parentgRPCAddressByteArray",
				"No data available...",
				err.Error(),
				"Error when converting 'ParentgRPCAddress' into a byte array, stopping futher processing.",
			)
		} else {
			// Marshaling to []byte OK

			// Save message to local DB for later processing
			_ = SaveMessageToLocalDB(
				BucketKeyForParentAddress,
				parentgRPCAddressByteArray,
				BucketForParentAddress,
				"36dc5c00-64b6-4122-af7f-69962442889e",
			)
		}

	}
	return true, nil

}

// *******************************************************************
// Loop all stored clients in DB and ask them to reRegister themself to this Gateway
//

func askClientsToReRegisterTHemSelf() {
	// Initiate map used for handle Clients address and port info
	// tabort detta här clientsAddressAndPort = make(map[string]clientsAddressAndPortStruct)
}

// *******************************************************************
// Cleanup before closing down
//
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		closeDB()

		// Stop gRPC-relateds listing and close towards Fenix
		stopGatewayGRPCServerForMessagesTowardsFenix()

		// Stop gRPC-relateds listing and close towards Plugins
		stopGatewayGRPCServerForMessagesTowardsPlugins()

	}
}

// *******************************************************************
// Start all Services
//
func startAllServices() {
	// Read 'gatewayConfig.toml' for config parameters
	processConfigFile("") // Use default toml-config-file name

	// Init logger
	initLogger("localLogFile.log")

	// Initiate internal gatewau channels
	initiateGatewayChannels()

	//  Initiate the memory structure to hold all client gateway/plugin's address information
	initiateClientAddressMemoryDB()

	// Ensure that all services don't start before everything has been started
	gatewayMustStopProcessing = true

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

	// Update Memory object with parameters that was recceived at start up using flags
	updateMemoryObjectWithFlagOverrideParameters()

	// Try to Register this Gateway At Parent
	gatewayTowardsFenixObject.tryToRegisterGatewayAtParent()

	// Listen to gRPC-calls from parent gateway/Fenix
	//startGatewayGRPCServerForMessagesTowardsFenix()

	// Listen to gRPC-calls from child gateway/Plugin
	startGatewayGRPCServerForMessagesTowardsPlugins()

	// Update Memory information about parent address and port with that saved in database, database overrule config-file
	updateMemoryAddressForParentAddressInfo()

	// Start all services at the same time
	gatewayMustStopProcessing = false

	// Ask clients to ReRegister them self to this gateway
	// TODO Make all Clients ReRegister them self

	// Release all saved messages to channls
	// TODO Release all massages to channels
}
