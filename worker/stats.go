package worker

import (
	"fmt"
	"os"
	"strconv"
	"time"

	procreader "github.com/c9s/goprocinfo/linux"
)

// trackStatsGoroutine is meant to track the performance statistics of the given worker, and send them to master
func trackStatsGoroutine() {
	for {
		time.Sleep(1 * time.Second)
		nodeAvailableMemory := readAvailableMemory()
		workerMemoryUsage := readMemoryUsage()
		WorkerStatistics.UpdateMemoryUsage(workerMemoryUsage, nodeAvailableMemory)
		fmt.Println("====Worker Statistics for Stage " + StageID + " ====")
		fmt.Println(WorkerStatistics)
	}
}

// readAvailableMemory reads the /proc file system to find the amount of memory available on the node
func readAvailableMemory() uint64 {
	procPath := "/proc/meminfo"
	systemMemoryInfo, err := procreader.ReadMemInfo(procPath)
	if err != nil {
		logMessage(err.Error())
		panic(err)
	}
	availableMemory := systemMemoryInfo.MemAvailable
	if availableMemory == 0 { // MemAvailable doesn't always work. If it doesn't, use MemFree instead
		availableMemory = systemMemoryInfo.MemFree
	}
	return availableMemory
}

// readMemoryUsage reads the /proc file system to find the amount of memory used by the worker process
func readMemoryUsage() uint64 {
	procPath := "/proc/" + strconv.Itoa(os.Getpid()) + "/statm"
	procStatm, err := procreader.ReadProcessStatm(procPath)
	if err != nil {
		logMessage(err.Error())
		panic(err)
	}
	return procStatm.Size
}
