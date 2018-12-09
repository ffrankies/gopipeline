package main

import (
	"fmt"
	"strconv"
	"time"
)

func hello(num int) {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello World from Position: " + strconv.Itoa(num))
}

func main() {
	fmt.Println("Started execution at", strconv.FormatInt(time.Now().UnixNano(), 10))
	for i := 0; i < 100; i++ {
		for i := 0; i < 20; i++ {
			hello(i)
		}
	}
	fmt.Println("Finished execution at", strconv.FormatInt(time.Now().UnixNano(), 10))
}
