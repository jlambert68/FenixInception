// +build integration2

package TestExecutionGateway

import (
	"flag"
	"fmt"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"
)

//func Test(t *testing.T) {
//	grpctest.RunSubTests(t, s{})
//}

var (
	database     = flag.Bool("database", false, "run database integration tests")
	messageQueue = flag.Bool("messageQueue", false, "run message queue integration tests")
)

var (
	logfileForTest = "localLogFileForTest.log"
	databaseFile   = "DbForTest.db"
)

func TestMain(m *testing.M) {

	// Set Gateway in IsInSelfIntegrationTestMode
	GatewayInIntegrationTestMode.IsInSelfIntegrationTestMode = true

	// Parse flags
	flag.Parse()

	// Initiate part 1 for gateway
	InitGatewayPart1()

	/*

		if *database {
		}
		if *messageQueue {
			setupMessageQueue()
		}
	*/

	// Cleanup all gRPC connections
	//log.Println("Process 'defer cleanup()'")
	//defer cleanup()

	result := m.Run()

	// Close database
	log.Println("Process 'closeDB'")
	closeDB()

	// Remove the database file that was created in this test
	removeFile(databaseFile)

	// Remove the logfile that was created in this test
	removeFile(logfileForTest)

	/*

		if *database {
			teardownDatabase()
		}
		if *messageQueue {
			teardownMessageQueue()
		}
	*/

	fmt.Println("ending test now")
	os.Exit(result)
}

// **************************************************************************
// Initiate the Gateway
//
func InitGatewayPart1() {

	// Read 'gatewayConfig.toml' for config parameters
	log.Println("Process 'processConfigFile'")
	processConfigFile("") // Use default toml-config-file name

	// Init logger
	log.Println("Process 'initLogger'")
	initLogger(logfileForTest)

	// Initiate internal gatewau channels
	log.Println("Process 'InitiateGatewayChannels'")
	InitiateGatewayChannels()

	//  Initiate the memory structure to hold all client gateway/plugin's address information
	log.Println("Process 'initiateClientAddressMemoryDB'")
	initiateClientAddressMemoryDB()

	// Ensure that all services don't start before everything has been started
	log.Println("Process 'gatewayMustStopProcessing'")
	common_code.GatewayMustStopProcessing = true

	// Initiate Database
	log.Println("Process 'initiateDB'")
	initiateDB(databaseFile)

	// Start 'dispatchEngineForTestInstructions'
	log.Println("Process 'gatewayTowardsPluginObject.initiateDispatchEngineForTestInstructiona'")
	gatewayTowardsPluginObject.initiateDispatchEngineForTestInstructiona()

	// Start 'dispatchEngineForSupportedTestDomains'
	log.Println("Process 'gatewayTowardsPluginObject.initiateDispatchEngineForSupportedTestDomains'")
	gatewayTowardsPluginObject.initiateDispatchEngineForSupportedTestDomains()

	// Start 'transmitEngineForRegistrateAvailableTestDataDomainsTowardsFenix'
	log.Println("Process 'gatewayTowardsFenixObject.initiateRegistrateAvailableTestDataDomainsTowardsFenix'")
	gatewayTowardsFenixObject.initiateRegistrateAvailableTestDataDomainsTowardsFenix()

	// Start 'transmitEngineForSendMessageToFenix'
	log.Println("Process 'gatewayTowardsFenixObject.initiateSendMessageToFenix'")
	gatewayTowardsFenixObject.initiateSendMessageToFenix()

	// Start 'transmitEngineForSendTestExecutionLogTowardsFenix'
	log.Println("Process 'gatewayTowardsFenixObject.initiateSendTestExecutionLogTowardsFenix('")
	gatewayTowardsFenixObject.initiateSendTestExecutionLogTowardsFenix()

	// Start 'transmitEngineForSendTestInstructionTimeOutTowardsFenix'
	log.Println("Process 'gatewayTowardsFenixObject.initiateSendTestInstructionTimeOutTowardsFenix'")
	gatewayTowardsFenixObject.initiateSendTestInstructionTimeOutTowardsFenix()

	// Try to Register this Gateway At Parent
	log.Println("Process 'gatewayTowardsFenixObject.tryToRegisterGatewayAtParent'")
	gatewayTowardsFenixObject.tryToRegisterGatewayAtParent()
	// Validate that log-file contains '580d2c7d-b8d3-40f7-b238-eb096d859355' because Parent Gateway is not started and should have been Exited

	/*
		// Listen to gRPC-calls from parent gateway/Fenix
		log.Println("Process 'startGatewayGRPCServerForMessagesTowardsFenix'")
		startGatewayGRPCServerForMessagesTowardsFenix()

		// Listen to gRPC-calls from child gateway/Plugin
		log.Println("Process 'startGatewayGRPCServerForMessagesTowardsPlugins'")
		startGatewayGRPCServerForMessagesTowardsPlugins()

		// Update Memory information about parent address and port with that saved in database, database overrule config-file
		log.Println("Process 'updateMemoryAddressForParentAddressInfo'")
		updateMemoryAddressForParentAddressInfo()

		// Start all services at the same time
		log.Println("Process 'gatewayMustStopProcessing = false'")
		gatewayMustStopProcessing = false
	*/
}

// *********************************************************************
// Check if string exists in file
//
func IsExist(str, filepath string) bool {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	isExist, err := regexp.Match(str, b)
	if err != nil {
		panic(err)
	}
	return isExist
}

// *********************************************************************
// Remove a file that was created in this test
//
func removeFile(fileToBeRemoved string) {
	// delete file
	var err = os.Remove(fileToBeRemoved)
	if err != nil {
		log.Printf("Couldn't delete file:%s\n", fileToBeRemoved)
		return
	}

	fmt.Printf("File '%s' Deleted\n", fileToBeRemoved)
}

// ************* Tests
// Validate that Gateway should have been stopped because parent gateway is not running
// and this gateway has never been connected to parent gateway
//mvalidateValueExistsInLogFile('580d2c7d-b8d3-40f7-b238-eb096d859355'')

func TestDatabase(t *testing.T) {
	// Pretend to use the db
	fmt.Println(common_code.GatewayConfig)
}

// Validate that gateway did a simulated exit because there are no parent gateway
func TestExitBecasueNoParentGatewayExits(t *testing.T) {
	stringFound := IsExist("580d2c7d-b8d3-40f7-b238-eb096d859355", logfileForTest)

	if stringFound == false {
		t.Errorf("Gateway didn't exit in a correct way")
	}
}
