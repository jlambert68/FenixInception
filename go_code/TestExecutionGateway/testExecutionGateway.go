package TestExecutionGateway

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)





// Used for only process cleanup once
var cleanupProcessed bool = false

func (gatewayObject *GatewayTowardsPluginObject_struct)  cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		gatewayObject.db.Close()
		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":            "3c84178b-a365-4d64-babb-4e1956d96684",
		}).Debug("Closing local database")
	}

/*
		// Cleanup before close down application
		GatewayCommonObjects.GatewayTowardsPluginObject	. motherObject.logger.WithFields(logrus.Fields{
			"ID": "ab1abcdb-d786-4bcc-b274-b52bd931f43d",
		}).Debug("Clean up and shut down servers")

		motherObject.logger.WithFields(logrus.Fields{
			"ID": "a22e8de2-c38c-4483-95e0-e9dc21dd7ca8",
		}).Debug("Gracefull stop for: registerTaxiHardwareStreamServer")

		registerMotherServer.GracefulStop()

		motherObject.logger.WithFields(logrus.Fields{
			"common_config.MotherServer_port: ": common_config.MotherServer_port,
			"ID": "19593393-3ac9-45ac-96e9-cdf911b167c7",
		}).Debug("Close net.Listing")
		lis.Close()

		//log.Println("Close DB_session: %v", DB_session)
		//DB_session.Close()
*/
	}
}


func TestExecution_main() {

	//var err error


	var GatewayObject *GatewayTowardsPluginObject_struct
	GatewayObject = &GatewayTowardsPluginObject_struct{}

	defer GatewayObject.cleanup()

	// Init logger
	GatewayObject.InitLogger("")


	// At startup Register Gateway to Parent Gateway/Fenix
	GatewayObject.registerThisGatewayAtParentGateway()

	// Find first non allocated port from defined start port
	motherObject.logger.WithFields(logrus.Fields{
		"ID": "56b6419f-d714-4ab0-be62-f3c7f08b9558",
	}).Debug("Mother Server tries to start")

	motherObject.logger.WithFields(logrus.Fields{
		"common_config.MotherServer_port): ": common_config.MotherServer_port,
		"ID": "8f904ace-9d24-452b-891a-5b8e5c247ba2",
	}).Debug("Start listening on:")
	lis, err = net.Listen("tcp", string(common_config.MotherServer_port))

	if err != nil {
		motherObject.logger.WithFields(logrus.Fields{
			"err: ": err,
			"ID": "0fbf0f08-6114-4cd8-9992-6a387955fb5f",
		}).Fatal("failed to listen:")

	} else {
		motherObject.logger.WithFields(logrus.Fields{
			"common_config.MotherServer_port): ": common_config.MotherServer_port,
			"ID": "93496c07-2b6c-4edc-a1f9-3fd555fa1201",
		}).Debug("Success in listening on port:")

	}

	// Creates a new RegisterMotherServer gRPC server
	go func() {
		motherObject.logger.WithFields(logrus.Fields{
			"ID": "ebc26735-9d13-4b13-91b8-20999cd5e254",
		}).Debug("Starting Mother Server")

		registerMotherServer = grpc.NewServer()
		mother_server_grpc_api.RegisterMotherServerServer(registerMotherServer, &MotherServer{})

		motherObject.logger.WithFields(logrus.Fields{
			"common_config.MotherServer_port): ": common_config.MotherServer_port,
			"ID": "cfed87c4-55aa-4cd1-980a-a15eb75ab6fb",
		}).Debug("registerMotherServer for Mother Server started")

		registerMotherServer.Serve(lis)
	}()

	// Start up 'reassebleEngine' to receive incoming worker objects
	go motherObject.reassebleEngine()

	//Start up 'sendWorkLoadToWorker' which sends workloads to workers
	go motherObject.sendWorkLoadToWorker()

	// Start up 'LoadWorkerObjectAndStoreRowsEngine' to extract rows and store them in Database
	go motherObject.LoadWorkerObjectAndStoreRowsEngine()

*/
	// Just waiting to quit
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	for {
		fmt.Println("sleeping...for another 5 minutes")
		time.Sleep(300 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}

	//Wait until user exit
	/*
   for {
	   time.Sleep(10)
   }
   */
}

