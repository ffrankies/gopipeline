package types

import (
	"strconv"
	"sync"
	"time"
)

// WorkerStats stores the performance statistics about a worker process's execution
type WorkerStats struct {
	NodeAvailableMemory  uint64        // The amount of memory available on the node (from /proc/meminfo/MemAvailable)
	WorkerMemoryUsage    uint64        // The amount of memory used by the worker process (from /proc/[pid]/statm/Size)
	MaxWorkerMemoryUsage uint64        // The maximum amount of memory used by the worker process
	ExecutionTime        time.Duration // The amount of time to process the worker's stage
	Backlog              int           // The number of unprocessed items in the input queue
	lock                 sync.Mutex    // For concurrency reasons
}

// String representation of WorkerStats
func (workerStats *WorkerStats) String() string {
	workerStats.lock.Lock()
	workerStatsString := "Worker stats: {"
	workerStatsString += " NodeAvailableMemory: " + strconv.FormatUint(workerStats.NodeAvailableMemory, 10)
	workerStatsString += " WorkerMemoryUsage: " + strconv.FormatUint(workerStats.WorkerMemoryUsage, 10)
	workerStatsString += " MaxWorkerMemoryUsage: " + strconv.FormatUint(workerStats.MaxWorkerMemoryUsage, 10)
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

// UpdateMemoryUsage uses a weighted running average to calculate the average memory usage of the given stage
func (workerStats *WorkerStats) UpdateMemoryUsage(memoryUsage uint64, availableMemory uint64) {
	newMemoryUsage := float64(memoryUsage) * (2. / 3.)
	workerStats.lock.Lock()
	var averageMemoryUsage uint64
	if workerStats.WorkerMemoryUsage == 0 {
		averageMemoryUsage = uint64(newMemoryUsage)
	} else {
		oldMemoryUsage := float64(workerStats.WorkerMemoryUsage) * (1. / 3.)
		averageMemoryUsage = uint64(oldMemoryUsage + newMemoryUsage)
	}
	if workerStats.MaxWorkerMemoryUsage < memoryUsage {
		workerStats.MaxWorkerMemoryUsage = memoryUsage
	}
	workerStats.NodeAvailableMemory = availableMemory
	workerStats.WorkerMemoryUsage = averageMemoryUsage
	workerStats.lock.Unlock()
}

// UpdateBacklog updates the backlog with the number of elements in the input queue
func (workerStats *WorkerStats) UpdateBacklog(backlog int) {
	workerStats.lock.Lock()
	workerStats.Backlog = backlog
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
