// Package types contains types needed by packages using the gopipeline library
package types

// AnyFunc is any function with any number of input parameters and a single return value
type AnyFunc func(...interface{}) interface{}
