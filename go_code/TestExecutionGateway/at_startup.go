package TestExecutionGateway

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gatewaygRPC "jlambert/FenixInception2/go_code/TestExecutionGateway/Gateway_gRPC_api"
)

func (gatewayObject *gatewayTowardsFenixObject_struct) AtStartUp() {

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

// Register this gateway/client at parent gateway/Fenix
func (gatewayObject *gatewayTowardsFenixObject_struct) registerThisGatewayAtParentGateway() (bool, error) {

	var err error
	var addressToDial string

	// Find parents address and port to call
	addressToDial = gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address + gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port

	// Information sent to parent gateway/Fenix
	registerClientAddressRequest := gatewaygRPC.RegisterClientAddressRequest{
		GRPCVersion:            gatewaygRPC.CurrentVersionEnum_VERSION_0_1_0,
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
		gatewayClient := gatewaygRPC.NewGatewayTowardsFenixClient(remoteGatewayServerConnection)

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

// Replace Parent gateway/Fenix IP-address & port info in memory object if previous connection differs from config-file
func getParentAddressInfo() {

	var parentAddress gatewaygRPC.ReRegisterToGatewayMessage
	var err error

	// Create the channel that the client address should be sent back on
	returnParentAddressChannel := make(chan dbResultMessage_struct)

	// Get Clients address
	dbMessage := dbMessage_struct{
		DB_READ,
		BUCKET_PARENT_ADDRESS,
		BUCKET_KEY_PARENT_ADDRESS,
		nil,
		returnParentAddressChannel}

	// Send Read message to database to receive address
	dbMessageQueue <- dbMessage

	// Wait for address from channel, then close the channel
	parentAddressByteArray := <-returnParentAddressChannel
	close(returnParentAddressChannel)

	// Convert saved json object into Go-struct
	err = json.Unmarshal(parentAddressByteArray.value, &parentAddress)
	if err != nil {
		// Problem with unmarshal the json object
		logger.WithFields(logrus.Fields{
			"ID":                      "261a6391-abc0-4d9f-a59f-3d0a67e5e52c",
			"parentAddressByteArray,": parentAddressByteArray,
		}).Error("Can't unmarshal gRPCClients address object from database")
		//TODO Send Error information to Fenix
	} else {

		// If Saved data differs from memory data then change in memory object
		// First check ip address
		if parentAddress.GatewayAddress != gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address {
			logger.WithFields(logrus.Fields{
				"ID":                           "560c2b17-c71e-45dd-9a38-a3dfd1a2bbd6",
				"parentAddress.GatewayAddress": parentAddress.GatewayAddress,
				"gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address": gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address,
			}).Warning("Ip-address for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")
			//TODO Send Error information to Fenix

			// Change Address in memory object
			gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_address = parentAddress.GatewayAddress
		}

		// Second check port
		if parentAddress.GatewayPort != gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port {
			logger.WithFields(logrus.Fields{
				"ID":                        "50a3b7ad-6631-42c5-ab5c-777e04ad9728",
				"parentAddress.GatewayPort": parentAddress.GatewayPort,
				"gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port": gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port,
			}).Warning("Port for Parent Gateway/Fenix differs for saved in DB and memory object, use DB-version")
			//TODO Send Error information to Fenix

			// Change Port in memory object
			gatewayConfig.parentgRPCAddress.parentGatewayInitialServer_port = parentAddress.GatewayPort
		}

	}
}
