// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"fmt"

	"github.com/ffrankies/gopipeline"
)

// Run executes the master process code on the given nodelist and list of functions that make up the pipelined module.
// The module parts functions will be executed in the order in which they appear in the moduleParts slice.
func Run(nodeListPath string, moduleParts []gopipeline.AnyFunc) {
	sshConnection := NewSSHConnection("rlogin.cs.vt.edu", "wanyef", 22)
	defer sshConnection.Close()
	out, err := sshConnection.RunCommand("ls")
	if err != nil {
		panic("Failed to run ls: " + err.Error())
	}
	fmt.Println(out)
	out, err = sshConnection.RunCommand("ps")
	if err != nil {
		panic("Failed to run ls: " + err.Error())
	}
	fmt.Println(out)
}
