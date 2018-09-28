// Contains the logic for running the master process of the gopipeline library
package master

import (
	"fmt"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/internal/common"
)

// Executes the master process code on the given nodelist and list of functions that make up the pipelined module.
// The module parts functions will be executed in the order in which they appear in the moduleParts slice.
func MasterRun(nodeListPath string, moduleParts []gopipeline.AnyFunc) {
	common.ReadNodeList(nodeListPath)
	fmt.Println("Number of functions passed =", len(moduleParts))
}
