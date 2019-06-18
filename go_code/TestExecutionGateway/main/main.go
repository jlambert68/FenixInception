package main

import (
	"flag"
	//	"fmt"
	"jlambert/FenixInception2/go_code/TestExecutionGateway"
	//	"time"
)

func main() {

	// Declare flags
	testIsInSelfIntegrationTestMode := flag.Bool(
		"selfIntegrationTestMode",
		false,
		"Used for Integration Tests and this gateway is the one that is testet.")

	testActAsParent := flag.Bool(
		"useWithOutParent",
		false,
		"Used for Integration Tests and this gateway could be started without any parent of its own.")

	testListenOnThisPort := flag.Int64(
		"asParentIListenOnThisPort",
		0,
		"Used for Integration Tests and this gateway should listen, from clients, on this port.")

	// Once all flags are declared, call `flag.Parse()`
	// Execute the command-line parsing.
	flag.Parse()

	// Transform flags to gateway parameters
	TestExecutionGateway.GatewayInIntegrationTestMode.IsInSelfIntegrationTestMode = *testIsInSelfIntegrationTestMode
	TestExecutionGateway.GatewayInIntegrationTestMode.StartWithOutAnyParent = *testActAsParent
	TestExecutionGateway.GatewayInIntegrationTestMode.ListeningOnThisPortAsParent = *testListenOnThisPort

	// Call Main function for Gateway
	TestExecutionGateway.TestExecution_main()
}
