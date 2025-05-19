package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sapelkinav/javadap/awesomejdwp/client"
	"sapelkinav/javadap/launcher"
	"sapelkinav/javadap/utils"
	"syscall"
	"time"
)

const LOG_DIR = "./.logs"

func main() {

	if err := utils.InitializeLogger(LOG_DIR, "debug"); err != nil {
		fmt.Printf("Failed to set up logger: %v\n", err)
		os.Exit(1)
	}

	jarPath := "daphelloworld/build/libs/daphelloworld-0.0.1-SNAPSHOT.jar"
	jdwpPort := 5005

	// Launch the Java application with JDWP enabled in server mode
	fmt.Println("Launching Java application with JDWP enabled...")
	launcher := launcher.NewJavaLauncher(jarPath, jdwpPort)
	if err := launcher.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch Java process: %v\n", err)
		os.Exit(1)
	}
	defer launcher.Stop()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	jdwpClient := client.NewJdwpClient("localhost:5005", ctx)

	err := jdwpClient.Connect()
	defer jdwpClient.Close()

	fmt.Println(err)
	jdwpClient.HelloWorld()

	gracefulShutdown(cancel)

}

func gracefulShutdown(
	cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	running := true
	for running {
		select {
		case sig := <-signalChan:
			fmt.Printf("Received signal: %v. Shutting down gracefully...\n", sig)
			cancel()        // Cancel context to notify all goroutines
			running = false // Exit the loop
		default:
			// Do periodic work here if needed
			// For example, check status, perform periodic tasks, etc.

			// Sleep a bit to prevent CPU spinning
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("Application shutdown complete")
}
