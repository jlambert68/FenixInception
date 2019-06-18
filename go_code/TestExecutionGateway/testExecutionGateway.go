package TestExecutionGateway

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TestExecution_main() {

	// Cleanup all gRPC connections
	defer cleanup()

	// Start all Services
	startAllServices()

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
