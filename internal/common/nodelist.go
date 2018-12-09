package common

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ReadNodeList reads in a file containing the list of nodes on which to run the pipelined module.
// Each node is expected to be on a new line. No checking is done for whether or not the nodes are valid, however empty
// lines are filtered out.
func ReadNodeList(path string) []string {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	contentsString := string(contents)
	contentsArray := strings.Split(contentsString, "\n")
	filteredContentsArray := make([]string, 0)
	for _, value := range contentsArray {
		if len(value) > 1 {
			filteredContentsArray = append(filteredContentsArray, value)
		}
	}
	for _, value := range filteredContentsArray {
		fmt.Println(value)
	}
	return filteredContentsArray
}
