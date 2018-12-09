package worker

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"sync"
)

// logging mutex
var logMutex = &sync.Mutex{}

// setupLogFile sets up the log file for this worker
func setupLogFile() {
	fmt.Println("Setting up filepath")
	filePath := getLogFilePath()
	// Create directory structure to filePath if it doesn't already exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}
	fmt.Println("Created log directory")
	// Delete log file from previous run
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		os.Remove(filePath)
	}
	logPrint("Logger has been setup")
}

// getLogFilePath returns the path to the log file
func getLogFilePath() string {
	userPath := userHomeDir()
	filePath := userPath + "/gopipeline_logs/gopipeline" + StageNumber + "." + StageID + ".log"
	return filePath
}

// userHomeDir returns the current user's home directory
func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

// opens a log file in the user's home directory
func openLogFile() (fp *os.File) {
	filePath := getLogFilePath()
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// logPrint prints a message to the log file
func logPrint(message string) {
	logMutex.Lock()
	f := openLogFile()
	defer f.Close()
	log.SetOutput(f)
	message = "Worker " + StageID + " | Stage " + StageNumber + ": " + message
	log.Println(message)
	logMutex.Unlock()
}

// logMessage prints a message to the console AND the log file
func logMessage(message string) {
	logPrint(message)
	message = "Worker " + StageID + ": " + message
	fmt.Println(message)
}
