package TestExecutionGateway

import (
	"flag"
	"fmt"
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
	gatewayUsedInIntegrationTest = flag.Bool("gatewayUsedInIntegrationTest", true, "true/false, deciding if the gateway should be ran in IntegrationTest mode or not.")
	databasePath                 = flag.String("datbasePath", "", "Relative path to database. This is a mandatory value, when running in Integration test mode, and must include database-name ending with '.db'.")
	logPath                      = flag.String("logdPath", "", "Relative path to log. This is a mandatory value, when running in Integration test mode, and must include logfile-name ending with '.log'.")
	configPath                   = flag.String("configPath", "", "Relative path to toml-config. This is a mandatory value, when running in Integration test mode, and must include config-name ending with '.toml'.")
	removeFilesAfterTest         = flag.Bool("removeFilesAfterTest", false, "true/false, deciding if temporary files, log and databas should be removed after test is finsihed")
)

func TestMain(m *testing.M) {

	// Parse flags
	flag.Parse()

	// Convert to Rune
	configPathRune := []rune(*configPath)
	logPathRune := []rune(*logPath)
	databasePathRune := []rune(*databasePath)

	// Get last part
	configPathEnding := string(configPathRune[len(*configPath)-5 : len(*configPath)])
	logPathEnding := string(logPathRune[len(*logPath)-4 : len(*logPath)])
	databasePathEnding := string(databasePathRune[len(*databasePath)-3 : len(*databasePath)])

	// Check for mandatory and faulty flags
	var exitBecauseOfMissingFlags = false
	if *databasePath == "" && *gatewayUsedInIntegrationTest == true {
		log.Println("'databasePath' is a mandtory flag when running in IntegrationTest-mode")
		exitBecauseOfMissingFlags = true
	}
	if *databasePath != "" && databasePathEnding != ".db" {
		log.Println("'databasePath' must end on '.db' to ensure a reference to a db file")
		exitBecauseOfMissingFlags = true
	}

	if *logPath == "" && *gatewayUsedInIntegrationTest == true {
		log.Println("'logPath' is a mandtory flag when running in IntegrationTest-mode")
		exitBecauseOfMissingFlags = true
	}
	if *logPath != "" && logPathEnding != ".log" {
		log.Println("'logPath' must end on '.log' to ensure a reference to a log file")
		exitBecauseOfMissingFlags = true
	}

	if *configPath == "" && *gatewayUsedInIntegrationTest == true {
		log.Println("'configPath' is a mandtory flag when running in IntegrationTest-mode")
		exitBecauseOfMissingFlags = true
	}
	if *configPath != "" && configPathEnding != ".toml" {
		log.Println("'configPath' must end on '.toml' to ensure a reference to a config file")
		exitBecauseOfMissingFlags = true
	}

	if exitBecauseOfMissingFlags == true {
		log.Fatalln("CLosing gateway du to missing or faulty mandatory flags")
	}

	// Initiate part 1 for gateway
	InitGatewayPart1(*configPath, *logPath, *databasePath)

	// Cleanup all gRPC connections
	//log.Println("Process 'defer cleanup()'")
	//defer cleanup()

	result := m.Run()

	// Close database
	log.Println("Process 'closeDB'")
	closeDB()

	// Check if should be removed before finish the test
	if *removeFilesAfterTest == true {

		// Remove the database file that was created in this test
		removeFile(*databasePath)

		// Remove the logfile that was created in this test
		removeFile(*logPath)
	}

	fmt.Println("ending test now")
	os.Exit(result)
}

// **************************************************************************
// Initiate the Gateway
//
func InitGatewayPart1(configFileAndPath string, logfileForTest string, databaseFile string) {

	// Read 'gatewayConfig.toml' for config parameters
	log.Println("Process 'processConfigFile'")
	processConfigFile(configFileAndPath) // If empty string then use default toml-config-file name

	// Init logger
	log.Println("Process 'initLogger'")
	initLogger(logfileForTest)

	// Initiate internal gatewau channels
	log.Println("Process 'initiateGatewayChannels'")
	initiateGatewayChannels()

	//  Initiate the memory structure to hold all client gateway/plugin's address information
	log.Println("Process 'initiateClientAddressMemoryDB'")
	initiateClientAddressMemoryDB()

	// Ensure that all services don't start before everything has been started
	log.Println("Process 'gatewayMustStopProcessing'")
	gatewayMustStopProcessing = true

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
	/*
		// Listen to gRPC-calls from parent gateway/Fenix
		log.Println("Process 'startGatewayGRPCServerForMessagesTowardsFenix'")
		startGatewayGRPCServerForMessagesTowardsFenix()
	*/
	// Validate that log-file contains '580d2c7d-b8d3-40f7-b238-eb096d859355' because Parent Gateway is not started and should have been Exited

	// Listen to gRPC-calls from child gateway/Plugin
	log.Println("Process 'startGatewayGRPCServerForMessagesTowardsPlugins'")
	startGatewayGRPCServerForMessagesTowardsPlugins()

	// Listen to gRPC-calls from child gateway/plugin
	log.Println("Process 'startGatewayGRPCServerForMessagesTowardsFenix'")
	startGatewayGRPCServerForMessagesTowardsFenix()

	// Update Memory information about parent address and port with that saved in database, database overrule config-file
	log.Println("Process 'updateMemoryAddressForParentAddressInfo'")
	updateMemoryAddressForParentAddressInfo()

	// Start all services at the same time
	log.Println("Process 'gatewayMustStopProcessing = false'")
	gatewayMustStopProcessing = false

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
	fmt.Println(gatewayConfig)
}

// Validate that gateway did a simulated exit because there are no parent gateway
func TestExitBecasueNoParentGatewayExits(t *testing.T) {
	stringFound := IsExist("580d2c7d-b8d3-40f7-b238-eb096d859355", *logPath)

	if stringFound == false {
		t.Errorf("Gateway didn't exit in a correct way")
	}
}
