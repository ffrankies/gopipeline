package common

// AnyFunc is any function with any number of input parameters and a single return value
type AnyFunc func(...interface{}) interface{}
