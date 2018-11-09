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
