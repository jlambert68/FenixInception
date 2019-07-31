package TestExecutionGateway

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
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

	// Start all Dispatch- and Transmit-Engines
	log.Println("Process 'initiateAllTransmitAndDispatchEngines'")
	initiateAllTransmitAndDispatchEngines()

	// Try to Register this Gateway At Parent
	log.Println("Process 'gatewayTowardsFenixObject.tryToRegisterGatewayAtParent'")
	tryToRegisterGatewayAtParent()
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
	common_code.GatewayMustStopProcessing = false

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

// ********************************************************************************************
// Send InformationMessage-messages to this Gateway to simulate message from child towards Fenix
//

func SendMessageToFenix(informationMessageToBeForwarded *gRPC.InformationMessage) error {

	var err error = nil

	// ***** Send InfoMessage to THIS gateway using gRPC-call ****
	ThisAddressAndPortInfo := common_code.GatewayConfig.GatewayIdentification
	addressToDial := ThisAddressAndPortInfo.GatewayIpAddress + ":" + strconv.FormatInt(int64(ThisAddressAndPortInfo.GatewayChildrenCallOnThisPort), 10)

	// Set up connection to Parent Gateway or Fenix
	remoteParentServerConnection, err := grpc.Dial(addressToDial, grpc.WithInsecure())
	if err != nil {
		// Connection Not OK
		LogErrorAndSendInfoToFenix(
			"ad0ce1a4-7755-45f5-b41c-5c80cdaa1b2d",
			gRPC.InformationMessage_WARNING,
			"addressToDial",
			addressToDial,
			err.Error(),
			"Did not connect to Child (Gateway or Plugin) Server!",
		)

		// Convert informationMessageToBeForwarded-struct into a byte array
		informationMessageToBeForwardedByteArray, err := json.Marshal(*informationMessageToBeForwarded)

		if err != nil {
			// Error when Unmarshaling to []byte
			LogErrorAndSendInfoToFenix(
				"6367c1cf-a81f-40dc-a089-4ca6ac90ab12",
				gRPC.InformationMessage_FATAL,
				"testExecutionLogMessageToBeForwarded",
				informationMessageToBeForwarded.String(),
				err.Error(),
				"Error when converting 'informationMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
			)
			return err

		} else {
			// Marshaling to []byte OK

			// Save message to local DB for later processing
			_ = SaveMessageToLocalDB(
				informationMessageToBeForwarded.MessageId,
				informationMessageToBeForwardedByteArray,
				common_code.BucketForResendOfInfoMessagesTowardsFenix,
				"2154b0cc-fdf3-49bd-8e9e-2ebd24710359",
			)
		}
	} else {
		//Connection OK

		// Convert testExecutionLogMessageToBeForwarded-struct into a byte array
		informationMessageToBeForwardedByteArray, err := json.Marshal(*informationMessageToBeForwarded)

		if err != nil {
			// Error when Unmarshaling to []byte
			LogErrorAndSendInfoToFenix(
				"f168c627-00f2-4ba0-a1c8-94782cc27f8e",
				gRPC.InformationMessage_FATAL,
				"informationMessageToBeForwarded",
				informationMessageToBeForwarded.String(),
				err.Error(),
				"Error when converting 'informationMessageToBeForwarded' into a byte array, stopping futher processing of this TestInstruction",
			)
			return err

		} else {
			// Marshaling to []byte OK

			// Save message to local DB for later processing
			_ = SaveMessageToLocalDB(
				informationMessageToBeForwarded.MessageId,
				informationMessageToBeForwardedByteArray,
				common_code.BucketForResendOfInfoMessagesTowardsFenix,
				"01ed6538-efd9-4ce5-af5c-c7869e8a4eb1",
			)

			// Creates a new gateway Client
			gatewayClient := gRPC.NewGatewayTowardsFenixClient(remoteParentServerConnection)

			// ChangeSenderId to this gatway's SenderId before sending the data forward
			informationMessageToBeForwarded.SenderId = common_code.GatewayConfig.GatewayIdentification.GatewayId
			informationMessageToBeForwarded.SenderName = common_code.GatewayConfig.GatewayIdentification.GatewayName

			// Do gRPC-call to client gateway or Fenix
			ctx := context.Background()
			returnMessage, err := gatewayClient.SendMessageToFenix(ctx, informationMessageToBeForwarded)
			if err != nil {
				// Error when rending gRPC to parent
				LogErrorAndSendInfoToFenix(
					"418510f2-2ca7-4f6f-8ba5-6c55463a059c",
					gRPC.InformationMessage_WARNING,
					"returnMessage",
					returnMessage.String(),
					err.Error(),
					"Problem to send 'informationMessageToBeForwarded' to parent-Gateway or Fenix",
				)

				// Save message to local DB for later processing
				_ = SaveMessageToLocalDB(
					informationMessageToBeForwarded.MessageId,
					informationMessageToBeForwardedByteArray,
					common_code.BucketForResendOfInfoMessagesTowardsFenix,
					"0d378676-3dd3-4bb9-a76b-66e410d5c280",
				)
				return err

			} else {
				// gRPC Send message OK
				common_code.Logger.WithFields(logrus.Fields{
					"ID":            "d9074bbc-a110-45de-8559-b063ec6122f1",
					"addressToDial": addressToDial,
				}).Debug("gRPC-send OK of 'informationMessageToBeForwarded' to Parent-Gateway or Fenix")
			}
			return nil
		}
	}
	return nil
}

// ************* Tests
// Validate that Gateway should have been stopped because parent gateway is not running
// and this gateway has never been connected to parent gateway
//mvalidateValueExistsInLogFile('580d2c7d-b8d3-40f7-b238-eb096d859355'')

func TestDatabase(t *testing.T) {

	// Create the channel that the client address should be sent back on
	returnParentAddressChannel := make(chan common_code.DbResultMessageStruct)

	// Get Clients address
	dbMessage := common_code.dbMessageStruct{
		common_code.DbRead,
		common_code.BucketForParentAddress,
		common_code.BucketKeyForParentAddress,
		nil,
		returnParentAddressChannel}

	// Send Read message to database to receive address
	common_code.dbMessageQueue <- dbMessage
	// Wait for address from channel, then close the channel
	databaseReturnMessage := <-returnParentAddressChannel
	close(returnParentAddressChannel)

	// Check if an error occured
	if databaseReturnMessage.err != nil {
		// Error when reading database
		t.Errorf("No Bucket found or No Key found when reading database")
	}

}

// Validate that gateway did a simulated exit because there are no parent gateway
func TestExitBecasueNoParentGatewayExits(t *testing.T) {
	stringFound := IsExist("580d2c7d-b8d3-40f7-b238-eb096d859355", *logPath)

	if stringFound == false {
		t.Errorf("Gateway didn't exit in a correct way")
	}
}

//SendMessageToFenix

// Validate that gateway did a simulated exit because there are no parent gateway
func TestSendMessageToFenix(t *testing.T) {
	stringFound := IsExist("580d2c7d-b8d3-40f7-b238-eb096d859355", *logPath)

	var infoMessageSent *gRPC.InformationMessage
	var infoMessageExpected *gRPC.InformationMessage
	//OriginalSenderId         string
	//	OriginalSenderName       string
	//	OriginalSystemDomainId   string
	//	OriginalSystemDomainName string
	//	SenderId                 string
	//	SenderName               string
	//	MessageId                string
	//	MessageType              InformationMessage_InformationType
	//	Message                  string

	// Generate a unique id that can be checked for late
	messageId := generateUUID()

	// The test-message to send
	infoMessageSent = &gRPC.InformationMessage{
		OriginalSenderName:       "OriginalSenderName",
		OriginalSystemDomainId:   "OriginalSystemDomainId",
		OriginalSystemDomainName: "OriginalSystemDomainName",
		SenderId:                 "SenderId",
		SenderName:               "SenderName",
		MessageId:                messageId,
		MessageType:              gRPC.InformationMessage_DEBUG,
		Message:                  "Message",
	}
	// The expected test-message to find in log
	infoMessageExpected = &gRPC.InformationMessage{
		OriginalSenderName:       "OriginalSenderName",
		OriginalSystemDomainId:   "OriginalSystemDomainId",
		OriginalSystemDomainName: "OriginalSystemDomainName",
		SenderId:                 "b34c377c-65d8-4050-afa0-72c30023a65f",
		SenderName:               "Gateway to be tested#1",
		MessageId:                messageId,
		MessageType:              gRPC.InformationMessage_DEBUG,
		Message:                  "Message",
	}

	// Send the message to gRPC
	err := SendMessageToFenix(infoMessageSent)

	if stringFound == false {
		t.Errorf("Gateway didn't exit in a correct way")
	}
}
