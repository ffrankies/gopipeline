// A more involved example, running a simple version of the apriori algorithm
package main

import (
	"math/rand"
	"time"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/types"
)

// Generate a sets of numbers from which the apriori algorithm will select common subsets
// This method does not take in any arguments, but they are present in the function signature for type compatibility
func Generate(args ...interface{}) interface{} {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)
	setSize := 20
	numSets := 200
	sets := make([][]int, 0)
	for i := 0; i < numSets; i++ {
		for j := 0; j < setSize; j++ {
			set := randomGenerator.Perm(30)[:setSize]
			sets = append(sets, set)
		}
	}
	return sets
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	// for i := 0; i < 5; i++ {
	// 	functionList = append(functionList, hello)
	// }
	gopipeline.Run(functionList)
}
