// Package types contains types needed by packages using the gopipeline library
package types

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
	NodeAvailableMemory uint64 // The amount of memory available on the node (from /proc/meminfo/MemAvailable)
	WorkerMemoryUsage   uint64 // The amount of memory used by the worker process (from /proc/[pid]/statm/Size)
	ExecutionTime       uint64 // The amount of time to process the worker's stage
	Backlog             int    // The number of unprocessed items in the input queue
}
