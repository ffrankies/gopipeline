// A more involved example, running a simple version of the apriori algorithm
package main

import (
	"math/rand"
	"sort"
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
	TargetNumSets      int      // The number of sets to keep from each iterations
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

// ContainsSupersetOf checks whether or not the set list contains a set that is a superset of the given set
func (setList *SetList) ContainsSupersetOf(subset *Set) bool {
	for _, set := range setList.List {
		if set.SupersetOf(subset) {
			return true
		}
	}
	return false
}

// Add a set to the set list
func (setList *SetList) Add(set *Set) {
	if !setList.ContainsSupersetOf(set) {
		setList.List = append(setList.List, set)
	}
}

// Set of integers, along with the number of times the set appears in OriginalSet
type Set struct {
	present map[int]bool // Used for a presence check
	Values  []int        // The values in the set
	Count   int          // The number of times this set appears in the superset
}

// NewSet creates a new set with an initialized map
func NewSet() *Set {
	set := new(Set)
	set.present = make(map[int]bool)
	return set
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
	for _, item := range set.Values {
		if value == item {
			return true
		}
	}
	return false
}

// SupersetOf checks if the current set is a superset of another set
func (set *Set) SupersetOf(otherSet *Set) bool {
	numElementsPresent := 0
	for _, value := range otherSet.Values {
		if set.Contains(value) {
			numElementsPresent++
		}
	}
	isSupersetOf := numElementsPresent == len(otherSet.Values)
	return isSupersetOf
}

// Add a value to the set, if it isn't already in the set. Returns a boolean that can be used to check whether the
// value was added or not
func (set *Set) Add(value int) bool {
	if !set.Contains(value) {
		set.Values = append(set.Values, value)
		set.present[value] = true
		return true
	}
	return false
}

// Split returns a SetList containing each element in the set as a separate set
func (set *Set) Split() *SetList {
	splitSet := new(SetList)
	for _, value := range set.Values {
		newSet := NewSet()
		newSet.Add(value)
		splitSet.Add(newSet)
	}
	return splitSet
}

// Copy returns a copy of the current set
func (set *Set) Copy() *Set {
	copy := NewSet()
	for _, value := range set.Values {
		copy.Add(value)
	}
	return copy
}

// Generate a sets of numbers from which the apriori algorithm will select common subsets
// This method does not take in any arguments, but they are present in the function signature for type compatibility
func Generate(arg interface{}) interface{} {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)
	setSize := 40
	numSets := 150
	sets := new(SetList)
	for i := 0; i < numSets; i++ {
		set := NewSet()
		// Generate setSize * 2 random numbers
		// Cut each set to size setSize
		values := randomGenerator.Perm(setSize * 2)[:setSize]
		for _, value := range values {
			set.Add(value)
		}
		sets.Add(set)
	}
	params := NextIteration(Parameters{OriginalSet: sets, LenCurrentSetItems: 0, TargetNumSets: 100})
	return params
}

// NextIteration creates the next iteration of sets for the a-priori algorithm
func NextIteration(arg interface{}) interface{} {
	params := arg.(Parameters)
	currentSet := new(SetList)
	if params.LenCurrentSetItems == 0 {
		currentSet = BuildInitialSets(params.OriginalSet, params.TargetNumSets)
	} else {
		currentSet = BuildSuccessiveSets(params.OriginalSet, params.CurrentSet, params.TargetNumSets)
	}
	params.CurrentSet = currentSet
	params.LenCurrentSetItems++
	return params
}

// BuildInitialSets creates the first iteration of sets, containing single values
func BuildInitialSets(originalSetList *SetList, targetNumSets int) *SetList {
	uniqueValues := GetUniqueValues(originalSetList)
	currentSetList := uniqueValues.Split()
	filterFrequencies := CalculateFrequencies(originalSetList, currentSetList, targetNumSets)
	currentSetList = FilterSetListByFrequency(currentSetList, filterFrequencies)
	return currentSetList
}

// GetUniqueValues returns a set containing all unique values in the given setList
func GetUniqueValues(setList *SetList) *Set {
	uniqueValueList := NewSet()
	for _, set := range setList.List {
		for _, value := range set.Values {
			if !uniqueValueList.Contains(value) {
				uniqueValueList.Add(value)
			}
		}
	}
	return uniqueValueList
}

// CalculateFrequencies for every set in the current set list, and returns the average frequency
func CalculateFrequencies(originalSetList *SetList, currentSetList *SetList, targetNumSets int) []int {
	sum := 0
	for _, set := range currentSetList.List {
		frequency := originalSetList.Frequency(set)
		set.Count = frequency
		sum += frequency
	}
	targetFrequencies := make([]int, 0)
	for _, set := range currentSetList.List {
		targetFrequencies = append(targetFrequencies, set.Count)
	}
	sort.Ints(targetFrequencies)
	sliceStart := len(targetFrequencies) - targetNumSets
	if sliceStart < 0 {
		sliceStart = 0
	}
	targetFrequencies = targetFrequencies[sliceStart:]
	return targetFrequencies
}

// FilterSetListByFrequency filters out any set list whose frequency is less than the average
func FilterSetListByFrequency(setList *SetList, targetFrequencies []int) *SetList {
	filteredSetList := new(SetList)
	for _, set := range setList.List {
		for _, frequency := range targetFrequencies {
			if set.Count == frequency {
				filteredSetList.Add(set)
			}
		}
	}
	return filteredSetList
}

// BuildSuccessiveSets builds sets that are 1 value larger than the current sets, and filters them based on their
// frequency in the original list of sets
func BuildSuccessiveSets(originalSetList *SetList, currentSetList *SetList, targetNumSets int) *SetList {
	nextSetList := new(SetList)
	uniqueValues := GetUniqueValues(currentSetList)
	for _, set := range currentSetList.List {
		for _, value := range uniqueValues.Values {
			newSet := set.Copy()
			if newSet.Add(value) {
				nextSetList.Add(newSet)
			}
		}
	}
	filterFrequencies := CalculateFrequencies(originalSetList, nextSetList, targetNumSets)
	nextSetList = FilterSetListByFrequency(nextSetList, filterFrequencies)
	return nextSetList
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	functionList = append(functionList, Generate)
	for i := 0; i < 4; i++ {
		functionList = append(functionList, NextIteration)
	}
	gopipeline.Run(functionList, Parameters{})
}
