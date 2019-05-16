package GatewayCommonObjects

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net"

	//"jlambert/AllCombinations2/grpc_api/mother_server_grpc_api"
	//"jlambert/AllCombinations2/common_config"
	"database/sql"
	"google.golang.org/grpc"
	"jlambert/AllCombinations2/grpc_api/worker_server_grpc_api"
	"net"
)

type MotherObject_struct struct {
	// Logger object for handling all logging
	// TODO change logging to use Datebased names
	mainlogger *logrus.Logger

	workerIdChannel chan mother_server_grpc_api.WorkerResult

	workerIdToProcessChannel chan int64
	currentTaskName          string
	messageToWorkerChannel   chan worker_server_grpc_api.MessageToWorkerStruct
	fullTaskStructureObject  worker_server_grpc_api.CombinationObjectStruct
}

/*
type clientConnectionInformation_struct struct {
	ip string
	prt string
	uuId	string
}
*/

type clientList_struct struct {
	clientHasRegistered        bool
	clientIp                   string
	clientPort                 string
	clientRegistrationDateTime string
}

var motherObject *MotherObject_struct

var (
	registerGatewayTowardsFenixServer *grpc.Server
	lis                               net.Listener
)

var (
	// Standard Worker gRPC Server
	remoteWorkerServerConnection *grpc.ClientConn
	workerClient                 worker_server_grpc_api.WorkerServerClient

	mother_address_to_dial string = common_config.MotherServer_address + common_config.MotherServer_port
)

// Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type MotherServer struct{}

// The data stored in DB to keep verisons of a project
type projectsObjectInDB_struct struct {
	version_key         string
	time_nanoint        int64
	date_time_string    string
	project_name        string
	project_description string
	deleted             bool
	project_key         string
}

// Object used in transformation from/to GUI-Object to/from DB
type GUIObjectToSave_struct struct {
	Ref         string                                    `json:"$ref"`
	Headers     []worker_server_grpc_api.HeaderStruct     `json:"headers"`
	ValueSets   []worker_server_grpc_api.ValueSetStruct   `json:"valueSets"`
	Rules       []worker_server_grpc_api.RuleStruct       `json:"rules"`
	HeaderTypes []worker_server_grpc_api.HeaderTypeStruct `json:"headerTypes"`
	RuleTypes   []worker_server_grpc_api.RuleTypeStruct   `json:"ruleTypes"`
}

// Variables used for dicide what kind of object to save in DB
// 1) Ordinary Combination Object
// 2) HeadersAndHeaderValues received from jExcel
type objectTypeToSave int

const (
	OrdinaryCombinationObject    objectTypeToSave = 1
	HeadersAndHeaderValuesObject objectTypeToSave = 2
	NewOrdinaryCombinationObject objectTypeToSave = 3
)

// Used to set HeaderTypes
const (
	HeaderTypeStandard_id  = 1
	HeaderTypeStandardText = "Standard"
)
