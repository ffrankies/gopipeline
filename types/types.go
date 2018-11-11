// Package types contains types needed by packages using the gopipeline library
package types

import (
	"strconv"
	"sync"
	"time"
)

// AnyFunc is any function with any number of input parameters and a single return value
type AnyFunc func(...interface{}) interface{}

// Message is a generic form of the messages passed between GoPipeline nodes
type Message struct {
	Sender      string      // The ID Of the sender
	Description int         // The message description
	Contents    interface{} // The contents of the message, can be of any type
}

// MessageStageInfo is the message struct for sending a stage's information to master
type MessageStageInfo struct {
	Address string // The address of the stage
	PID     int    // The id of the worker process running the stage
}

// WorkerStats stores the performance statistics about a worker process's execution
type WorkerStats struct {
	NodeAvailableMemory uint64        // The amount of memory available on the node (from /proc/meminfo/MemAvailable)
	WorkerMemoryUsage   uint64        // The amount of memory used by the worker process (from /proc/[pid]/statm/Size)
	ExecutionTime       time.Duration // The amount of time to process the worker's stage
	Backlog             int           // The number of unprocessed items in the input queue
	lock                sync.Mutex    // For concurrency reasons
}

// String representation of WorkerStats
func (workerStats *WorkerStats) String() string {
	workerStats.lock.Lock()
	workerStatsString := "Worker stats: {"
	workerStatsString += " NodeAvailableMemory: " + strconv.FormatUint(workerStats.NodeAvailableMemory, 10)
	workerStatsString += " WorkerMemoryUsage: " + strconv.FormatUint(workerStats.WorkerMemoryUsage, 10)
	workerStatsString += " ExecutionTime: " + strconv.FormatInt(workerStats.ExecutionTime.Nanoseconds(), 10)
	workerStatsString += " Backlog: " + strconv.Itoa(workerStats.Backlog)
	workerStatsString += " }"
	workerStats.lock.Unlock()
	return workerStatsString
}

// UpdateExecutionTime uses a weighted running average to calculate the average execution time of incoming tasks
func (workerStats *WorkerStats) UpdateExecutionTime(executionTime time.Duration) {
	newExecutionTime := float64(executionTime) * (2. / 3.)
	workerStats.lock.Lock()
	var averageExecutionTime time.Duration
	if workerStats.ExecutionTime == 0 {
		averageExecutionTime = executionTime
	} else {
		oldExecutionTime := float64(workerStats.ExecutionTime) * (1. / 3.)
		averageExecutionTime = time.Duration(oldExecutionTime + newExecutionTime)
	}
	workerStats.ExecutionTime = averageExecutionTime
	workerStats.lock.Unlock()
}

// Copy returns a copy of the WorkerStats struct
func (workerStats *WorkerStats) Copy() *WorkerStats {
	workerStatsCopy := new(WorkerStats)
	workerStats.lock.Lock()
	workerStatsCopy.NodeAvailableMemory = workerStats.NodeAvailableMemory
	workerStatsCopy.WorkerMemoryUsage = workerStats.WorkerMemoryUsage
	workerStatsCopy.ExecutionTime = workerStats.ExecutionTime
	workerStatsCopy.Backlog = workerStats.Backlog
	workerStats.lock.Unlock()
	return workerStatsCopy
}
