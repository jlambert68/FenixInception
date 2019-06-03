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

func (gatewayObject *gatewayTowardsPluginObject_struct)  cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// CLose database
		gatewayObject.db.Close()
		logger.WithFields(logrus.Fields{
			"ID":            "3c84178b-a365-4d64-babb-4e1956d96684",
		}).Debug("Closing local database")
	}

/*
		// Cleanup before close down application
		GatewayCommonObjects.gatewayTowardsPluginObject	. motherObject.logger.WithFields(logrus.Fields{
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


	var GatewayObject *gatewayTowardsPluginObject_struct
	GatewayObject = &gatewayTowardsPluginObject_struct{}

	defer GatewayObject.cleanup()




	// At startup Register Gateway to Parent Gateway/Fenix
	GatewayObject.registerThisGatewayAtParentGateway()


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

