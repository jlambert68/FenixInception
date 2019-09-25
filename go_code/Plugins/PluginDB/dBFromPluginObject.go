package PluginKeyValueDBStore

import (
	"google.golang.org/grpc"
	"net"
)

var (
	registerKeyValueStoreServer *grpc.Server
	keyValueStoreListener       net.Listener

	// gRPC server used to handle all traffic Towards the Plugins

)

type gRPCServerForDbKeyValueStoreStruct struct{}
