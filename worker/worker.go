package worker

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

type Address struct {
	Data string
}

type ResultData interface {
	Data() string
	Result() string
}

var nextNodeAddress string

func Run(options *common.WorkerOptions, functionList []types.AnyFunc) {
	// Listens for both the master and any other connection
	ln, err := net.Listen("tcp", "0:8081")
	if err != err {
		panic(err)
	}
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Println(conn.RemoteAddr())
	//handle this connection based on master or other workers
	masterdetails := options.MasterAddress + ":" + string(options.MasterPort)
	//Sends my address as a struct data to the master.
	conn1, err := net.Dial("tcp", masterdetails)
	myAddress := conn1.LocalAddr()
	if err != nil {
		panic(err)
	}
	dataAddress := Address{Data: myAddress.String()}
	da := &dataAddress
	fmt.Println(da)
	encoder := gob.NewEncoder(conn1)
	encoder.Encode(da)

	//Listens for port of the next one.
	ln1, err := net.Listen("tcp", myAddress.String())
	if err != nil {
		panic(err)
	}
	for {
		conn1, err := ln1.Accept()
		if err != nil {
			continue
		}
		message, err := bufio.NewReader(conn1).ReadString('\n')
		//Should receive the address of the next node here
		//fmt.Print("Message Received:", string(message))
		nextNodeAddress := string(message)
	}

	//Check for the postion of the node.
	//Since it is the first node in the pipeline
	//No data from the previous is to be received then perform the function
	if options.Position == 0 {
		//Functionlist????
		result := functionList[options.Position]
		connNext, err := net.Dial("tcp", nextNodeAddress)
		results := ResultData{Result: result} //match the return type
		res := &results
		fmt.Println(res)
		encoder := gob.NewEncoder(connNext)
		encoder.Encode(res)

	} else if options.Position == len(functionList) {

	}

}
