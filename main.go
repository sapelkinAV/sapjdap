package main

import (
	"fmt"
	"os"
	"sapelkinav/javadap/launcher"
	"sapelkinav/javadap/utils"
)

const LOG_DIR = "./.sapelkin_debugger"

func main() {

	if err := utils.InitializeLogger("./logs", "debug"); err != nil {
		fmt.Printf("Failed to set up logger: %v\n", err)
		os.Exit(1)
	}

	jarPath := "daphelloworld/build/libs/daphelloworld-0.0.1-SNAPSHOT.jar"
	jdwpPort := 5005

	// Launch the Java application with JDWP enabled in server mode
	fmt.Println("Launching Java application with JDWP enabled...")
	launcher := launcher.NewJdwpLauncher(jarPath, jdwpPort)
	if err := launcher.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch Java process: %v\n", err)
		os.Exit(1)
	}
	defer launcher.Stop()

}
