package worker

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"sync"
)

// logging mutexes
var logMutex = &sync.Mutex{}
var performanceLogMutex = &sync.Mutex{}

// Performance log message types
const (
	PerfStartWorker string = "Worker started"
	PerfStartExec   string = "Stage execution started"
	PerfEndExec     string = "Stage execution ended"
)

// setupLogFile sets up the log file for this worker
func setupLogFile() {
	filePath := getLogFilePath()
	// Create directory structure to filePath if it doesn't already exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}
	// Delete log file from previous run
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		os.Remove(filePath)
	}
	logPrint("Logger has been setup")
}

// getLogFilePath returns the path to the log file
func getLogFilePath() string {
	userPath := userHomeDir()
	filePath := userPath + "/gopipeline_logs/" + StageNumber + "." + StageID + ".log"
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

// setupPerformanceLogFile sets up the performance logger
func setupPerformanceLogFile() {
	filePath := getPerformanceLogFilePath()
	// Create directory structure to filePath if it doesn't already exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}
	fmt.Println("Created log directory")
	// Delete log file from previous run
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		os.Remove(filePath)
	}
}

// getPerformanceLogFilePath returns the path to the log file
func getPerformanceLogFilePath() string {
	userPath := userHomeDir()
	filePath := userPath + "/gopipeline_logs/performance/" + StageNumber + "." + StageID + ".log"
	return filePath
}

// openPerformanceLogFile opens the performance log file
func openPerformanceLogFile() (fp *os.File) {
	filePath := getPerformanceLogFilePath()
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// logPerformance logs a performance message to file
func logPerformance(performanceMessage string) {
	performanceLogMutex.Lock()
	f := openLogFile()
	defer f.Close()
	log.SetOutput(f)
	message := performanceMessage + "," + StageID + "," + StageNumber + "," + performanceMessage
	log.Println(message)
	performanceLogMutex.Unlock()
}
