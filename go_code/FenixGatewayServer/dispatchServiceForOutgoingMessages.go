package FenixGatewayServer

// **********************************************************************************************************
// Save 'testInstructionExecutionResultMessage' to Main Database for Fenix Inception
//
func sentTestInstructionTowardsPlugin(testInstructionId string) (err error) {
	returnMessageAckNackResponse, gRPCerr = gRpcClientTowardPlugin.SendTestInstructionTowardsPlugin(gRpcContexType, testInstructionMessageToBeForwardedTowardsPlugin)
	returnMessageString = returnMessageAckNackResponse.String()

}
