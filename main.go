package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"os"
	"os/signal"
	"sapelkinav/javadap/jdwp/client"
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
	javaLauncher := launcher.NewJavaLauncher(jarPath, jdwpPort)
	if err := javaLauncher.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch Java process: %v\n", err)
		os.Exit(1)
	}
	defer javaLauncher.Stop()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var socket io.ReadWriteCloser
	var err error
	for i := 0; i < 5; i++ {
		if socket, err = net.Dial("tcp", fmt.Sprintf("localhost:%v", jdwpPort)); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if socket == nil {
		log.Warn().Err(err).Msg("Failed to connect to the socket. Error")
		return
	}

	con, err := client.Open(ctx, socket)
	defer socket.Close()

	version, err := con.GetVersion()
	if err != nil {
		fmt.Println(version)
	}

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
