package main

import (
	"flag"
	"log"
	//	"fmt"
	"jlambert/FenixInception2/go_code/TestExecutionGateway"
	//	"time"
)

func main() {

	// Declare flags
	var (
		gatewayUsedInIntegrationTest = flag.Bool("gatewayUsedInIntegrationTest", false, "True/False, deciding if hte gateway should be ran in IntegrationTest mode or not.")
		databasePath                 = flag.String("datbasePath", "", "Relative path to database. This is a mandatory value, when running in Integration test mode, and must include database-name ending with '.db'.")
		logPath                      = flag.String("logdPath", "", "Relative path to log. This is a mandatory value, when running in Integration test mode, and must include config-name ending with '.log'.")
		configPath                   = flag.String("configPath", "", "Relative path to toml-config. This is a mandatory value, when running in Integration test mode, and must include config-name ending with '.toml'.")
	)

	// Once all flags are declared, call `flag.Parse()`
	// Execute the command-line parsing.
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

	/*
		// Transform flags to gateway parameters
		TestExecutionGateway.GatewayInIntegrationTestMode.UsedInIntegrationTest = *testUsedInIntegrationTest
		TestExecutionGateway.GatewayInIntegrationTestMode.IsInSelfIntegrationTestMode = *testIsInSelfIntegrationTestMode
		TestExecutionGateway.GatewayInIntegrationTestMode.StartWithOutAnyParent = *testActAsParent
	*/
	// Call Main function for Gateway
	TestExecutionGateway.TestExecutionMain(*configPath, *logPath, *databasePath)
}
