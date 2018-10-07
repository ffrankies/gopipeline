// Package gopipeline contains the logic for pipelining a Golang application through a list of compute nodes.
package gopipeline

// AnyFunc is any function with any number of input parameters and a single return value
type AnyFunc func(...interface{}) interface{}
