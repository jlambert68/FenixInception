package TestExecutionGateway

import (
	"fmt"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TestExecutionMain(configFileAndPath string, logfileForTest string, databaseFile string) {

	// Cleanup all gRPC connections
	defer cleanup()

	// Start all Services as a Gateway Engine and no function references, use nil as function reference
	gatewayOrEndpoint := common_code.FunctionsInsteadOfgRPCStruct{
		FenixOrGatewayTypeOrPlugin: common_code.GatewayEngine,
		CallBackTowardsPlugins: common_code.CallBackTowardsPluginsType{
			CallBackSendTestInstructionTowardsPlugin:        nil,
			CallackGetSupportedTestDataDomainsTowardsPlugin: nil},
		CallBackTowardsFenix: common_code.CallBackTowardsFenixType{
			CallBackRegisterAvailbleTestInstructionsTowardsFenix:   nil,
			CallBackRegistrateAailableTestContainersTowardsFenix:   nil,
			CallBackRegistrateAvailableTestDataDomainsTowardsFenix: nil,
			CallBackSendMessageToFenixTowardsFenix:                 nil,
			CallBackSendTestInstructionTimeOutTowardsFenix:         nil,
			CallBackSendTestExecutionLogTowardsFenix:               nil,
			CallBackSupportedTestDataDomainsTowardsFenix:           nil,
			CallBackSendTestInstructionResultTowardsFenixType:      nil,
		},
	}
	StartAllServices(configFileAndPath, logfileForTest, databaseFile, gatewayOrEndpoint)

	// Just waiting to quit
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	for {
		fmt.Println("sleeping...for another 15 minutes")
		time.Sleep(1200 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}

}
