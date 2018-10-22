// A more involved example, running a simple version of the apriori algorithm
package main

import (
	"math/rand"
	"time"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/types"
)

// Parameters to be passed between functions. Passing a single struct with some nil arguments is easier than
// type casting every single value
type Parameters struct {
	OriginalSet        *SetList // The original, generated set of integers
	CurrentSet         *SetList // The current, in-progress set of integers
	LenCurrentSetItems int      // The length of the items in the current set
	TargetSetLength    int      // The length of the target frequent sets
}

// SetList is a list of sets
type SetList struct {
	List []*Set
}

// Frequency returns the frequency of a given subset within the set list
func (setList *SetList) Frequency(subset *Set) int {
	frequency := 0
	for _, set := range setList.List {
		if set.SupersetOf(subset) {
			frequency++
		}
	}
	return frequency
}

// Add a set to the set list
func (setList *SetList) Add(set *Set) {
	setList.List = append(setList.List, set)
}

// Set of integers, along with the number of times the set appears in OriginalSet
type Set struct {
	_present map[int]bool // Used for a presence check
	Values   []int        // The values in the set
	Count    int          // The number of times this set appears in the superset
}

// Equals checks if this set equals the other set
func (set *Set) Equals(otherSet *Set) bool {
	for index, value := range otherSet.Values {
		if value != set.Values[index] {
			return false
		}
	}
	return true
}

// Contains checks if the current set contains a given value
func (set *Set) Contains(value int) bool {
	_, isPresent := set._present[value]
	return isPresent
}

// SupersetOf checks if the current set is a superset of another set
func (set *Set) SupersetOf(otherSet *Set) bool {
	numElementsPresent := 0
	for _, value := range otherSet.Values {
		if set.Contains(value) {
			numElementsPresent++
		}
	}
	return numElementsPresent == len(otherSet.Values)
}

// Add a value to the set, if it isn't already in the set
func (set *Set) Add(value int) {
	if !set.Contains(value) {
		set.Values = append(set.Values, value)
	}
}

// Generate a sets of numbers from which the apriori algorithm will select common subsets
// This method does not take in any arguments, but they are present in the function signature for type compatibility
func Generate(args ...interface{}) interface{} {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)
	setSize := 50
	numSets := 300
	sets := new(SetList)
	for i := 0; i < numSets; i++ {
		for j := 0; j < setSize; j++ {
			set := new(Set)
			values := randomGenerator.Perm(75)[:setSize]
			for _, value := range values {
				set.Add(value)
			}
			sets.Add(set)
		}
	}
	return Parameters{OriginalSet: sets, LenCurrentSetItems: 0, TargetSetLength: 25}
}

// NextIteration creates the next iteration of sets for the a-priori algorithm
func NextIteration(args ...interface{}) interface{} {
	params := args[0].(Parameters)
	currentSet := new(SetList)
	if params.LenCurrentSetItems == 0 {
		currentSet = FirstIteration(params.OriginalSet)
	} else {
		// TODO: Build multiple-value sets
	}
	params.CurrentSet = currentSet
	params.LenCurrentSetItems++
	return params
}

// FirstIteration creates the first iteration of sets, containing single values
func FirstIteration(originalSetList *SetList) *SetList {
	currentSetList := new(SetList)
	presenceCheckerSet := new(Set)
	for _, set := range originalSetList.List {
		for _, value := range set.Values {
			if !presenceCheckerSet.Contains(value) {
				presenceCheckerSet.Add(value)
				newSet := new(Set)
				newSet.Add(value)
				currentSetList.Add(newSet)
			}
		}
	}
	averageFrequency := CalculateFrequencies(originalSetList, currentSetList)
	currentSetList = FilterSetListByFrequency(currentSetList, averageFrequency)
	return currentSetList
}

// CalculateFrequencies for every set in the current set list, and returns the average frequency
func CalculateFrequencies(originalSetList *SetList, currentSetList *SetList) float64 {
	sum := 0
	for _, set := range currentSetList.List {
		frequency := originalSetList.Frequency(set)
		set.Count = frequency
		sum += frequency
	}
	return float64(sum) / float64(len(currentSetList.List))
}

//FilterSetListByFrequency filters out any set list whose frequency is less than the average
func FilterSetListByFrequency(setList *SetList, averageFrequency float64) *SetList {
	filteredSetList := new(SetList)
	for _, set := range setList.List {
		if float64(set.Count) > averageFrequency {
			filteredSetList.Add(set)
		}
	}
	return filteredSetList
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	// for i := 0; i < 5; i++ {
	// 	functionList = append(functionList, hello)
	// }
	gopipeline.Run(functionList)
}
