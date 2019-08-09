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
	startAllServices(configFileAndPath, logfileForTest, databaseFile, common_code.FunctionsInsteadOfgRPCStruct{
		FenixOrGatewayTypeOrPlugin: common_code.GatewayEngine,
		FenixAndPluginFunctionMap: map[common_code.FunctionType]common_code.FunctionReference{
			common_code.ChannelTypeTestInstructionMessageTowardsPluginFunction:                    nil,
			common_code.ChannelTypeSupportedTestDataDomainsRequestMessageTowardsPluginFunction:    nil,
			common_code.ChannelTypeInformationMessageTowardsFenixFunction:                         nil,
			common_code.ChannelTypeTestInstructionTimeOutMessageTowardsFenixFunction:              nil,
			common_code.ChannelTypeTestExecutionLogMessageTowardsFenixFunction:                    nil,
			common_code.ChannelTypeAvailbleTestInstructionsAtPluginMessageTowardsFenixFunction:    nil,
			common_code.ChannelTypeAvailbleTestContainersAtPluginMessageTowardsFenixFunction:      nil,
			common_code.ChannelTypeTestInstructionExecutionResultMessageTowardsFenixFunction:      nil,
			common_code.ChannelTypeSupportedTestDataDomainsWithHeadersMessageTowardsFenixFunction: nil,
		},
	})

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
