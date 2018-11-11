package worker

import (
	"fmt"
	"os"
	"strconv"
	"time"

	procreader "github.com/c9s/goprocinfo/linux"
)

func getStatsGoRoutine() {
	for {
		time.Sleep(1 * time.Second)
		WorkerStatistics.NodeAvailableMemory = readAvailableMemory()
		WorkerStatistics.WorkerMemoryUsage = readMemoryUsage()
		fmt.Println("====Worker Statistics for Stage " + StageID + " ====")
		fmt.Println(WorkerStatistics)
	}
}

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

func readMemoryUsage() uint64 {
	procPath := "/proc/" + strconv.Itoa(os.Getpid()) + "/statm"
	procStatm, err := procreader.ReadProcessStatm(procPath)
	if err != nil {
		logMessage(err.Error())
		panic(err)
	}
	return procStatm.Size
}
