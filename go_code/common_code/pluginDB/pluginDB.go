package pluginDB

import (
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
)

//
const DbBucketForTestInstructionToBeExecutedMEssages = "DbBucketForTestInstructionToBeExecutedMEssages"

// *********************************************************************************
// Get highest gRPC version for KeyValue-DB
//
func GetHighestPluginGRPCVersion() (currentVersion string) {
	maxVersionCount := int32(len(gRPC.CurrentVersionEnum_name) - 1)

	maxVersion := gRPC.CurrentVersionEnum_name[maxVersionCount]
	return maxVersion
}
