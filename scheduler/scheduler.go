package scheduler

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// Schedule contains the information needed for scheduling
type Schedule struct {
	freeNodeList *types.PipelineNodeList  // The list of Nodes available for scheduling
	NodeList     *types.PipelineNodeList  // The list of Nodes that have at least one stages running on them
	StageList    *types.PipelineStageList // The list of pipeline Stages, with metadata
	sshUser      string                   // The username to use for logging in with SSH
	sshPort      int                      // The port to use for logging in with SSH
	sshUserPath  string                   // The path to the program command on the remote machines
}

// NewSchedule creates a new scheduler with empty node and stage lists, and populates the empty node list
func NewSchedule(nodeList []string, SSHUser string, SSHPort int, SSHUserPath string) *Schedule {
	schedule := new(Schedule)
	schedule.NodeList = types.NewPipelineNodeList()
	schedule.StageList = types.NewPipelineStageList()
	schedule.freeNodeList = types.NewPipelineNodeList()
	schedule.sshUser = SSHUser
	schedule.sshPort = SSHPort
	schedule.sshUserPath = SSHUserPath
	for _, nodeHostName := range nodeList {
		node := types.NewPipelineNode(nodeHostName, -1)
		schedule.freeNodeList.AddNode(node)
	}
	return schedule
}

// Static does initial static scheduling of the pipeline stages on the available nodes
func (schedule *Schedule) Static(functionList []types.AnyFunc) {
	fmt.Println("Performing static scheduling")
	density := schedule.CalculateFunctionDensity(functionList)
	counter := 0
	schedulingNode := schedule.freeNodeList.Pop()
	for index := range functionList {
		schedule.AssignStageToNode(index, schedulingNode)
		counter++
		if counter == density {
			counter = 0
			density = schedule.CalculateFunctionDensity(functionList)
			if density < 0 { // If scheduling is over, there are no functions to schedule, so density becomes < 0
				break
			}
			schedulingNode = schedule.freeNodeList.Pop()
		}
	}
}

// CalculateFunctionDensity calculates the initial function density in the pipeline
func (schedule *Schedule) CalculateFunctionDensity(functionList []types.AnyFunc) int {
	numFunctionsRemaining := len(functionList) - schedule.StageList.Length()
	numNodesRemaining := schedule.freeNodeList.Length()
	density := math.Ceil(float64(numFunctionsRemaining) / float64(numNodesRemaining))
	return int(density)
}

// AssignStageToFreeNode assigns a single pipeline stage to a single free node
func (schedule *Schedule) AssignStageToFreeNode(functionIndex int) *types.PipelineStage {
	if schedule.freeNodeList.Length() == 0 {
		panic("FATAL ERROR: There are no free nodes to assign this stage to")
	}
	schedulingNode := schedule.freeNodeList.Pop()
	pipelineStage := schedule.AssignStageToNode(functionIndex, schedulingNode)
	return pipelineStage
}

// AssignStageToNode assigns a single pipeline stage (function) to a single node
func (schedule *Schedule) AssignStageToNode(functionIndex int, pipelineNode *types.PipelineNode) *types.PipelineStage {
	_, foundInList := schedule.NodeList.FindNode(pipelineNode.Address)
	pipelineStage := schedule.StageList.AddStage(pipelineNode.Address, functionIndex)
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		schedule.NodeList.AddNode(pipelineNode)
	}
	return pipelineStage
}

// UpdateStageStats updates the worker statistics for a given stage from an incoming message
func (schedule *Schedule) UpdateStageStats(message *types.Message) {
	stage := schedule.StageList.Find(message.Sender)
	stageStats, ok := (message.Contents).(*types.WorkerStats)
	if ok {
		stage.Stats = stageStats
	} else {
		fmt.Println("ERROR: Could not convert message contents to WorkerStats")
	}
}

// UpdateStageInfo updates the stage information for a given stage from an incoming message
func (schedule *Schedule) UpdateStageInfo(message *types.Message) {
	fmt.Println("Received worker info from", message.Sender)
	stage := schedule.StageList.Find(message.Sender)
	stageInfo, ok := (message.Contents).(types.MessageStageInfo)
	if ok {
		stage.NetAddress = stageInfo.Address
		stage.PID = stageInfo.PID
	} else {
		fmt.Println("ERROR: Could not convert message contents to MessageStageInfo")
	}
}

// StartStages starts GoPipeline workers for all the current stages
func (schedule *Schedule) StartStages(program string, masterAddress string) {
	fmt.Println("Starting GoPipeline workers")
	for _, stage := range schedule.StageList.List {
		schedule.startStage(stage, program, masterAddress)
	}
}

// startStage starts a GoPipeline worker for a given stage
func (schedule *Schedule) startStage(stage *types.PipelineStage, program string, masterAddress string) {
	sshConnection := types.NewSSHConnection(stage.Host, schedule.sshUser, schedule.sshPort)
	command := buildWorkerCommand(program, masterAddress, stage, schedule.sshUserPath)
	fmt.Println("Running command:", command, "on node:", stage.Host)
	go sshConnection.RunCommand(command, workerErrorCallback, stage)
}

// workerErrorCallback is the callback for when a worker errors out and dies
func workerErrorCallback(args ...interface{}) {
	stage := args[0].(*types.PipelineStage)
	stage.PID = -2 // Mark stage as errored out
}

// buildWorkerCommand builds the command with which to start a worker.
// The User Path should have a "/" included in the path.
func buildWorkerCommand(program string, masterAddress string, stage *types.PipelineStage, userpath string) string {
	command := userpath + program + " -address=" + masterAddress
	command += " -id=" + stage.StageID
	command += " -position=" + strconv.Itoa(stage.Position)
	command += " worker"
	return command
}

// EstablishWorkerCommunication establishes initial communication between workers by telling them the address
// of the next worker in the pipeline
func (schedule *Schedule) EstablishWorkerCommunication() {
	numPositions := schedule.StageList.Length()
	for position := 1; position < numPositions; position++ {
		nextWorker := schedule.StageList.FindByPosition(position)
		currentWorker := schedule.StageList.FindByPosition(position - 1)
		sendNextWorkerAddress(currentWorker, nextWorker)
	}
}

// sendNextWorkerAddress sends the next worker's address to the given worker
func sendNextWorkerAddress(currentWorker *types.PipelineStage, nextWorker *types.PipelineStage) {
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgAddNextStageAddr
	message.Contents = nextWorker.NetAddress
	fmt.Println("Setting up connection to:", currentWorker.NetAddress)
	connection, err := net.Dial("tcp", currentWorker.NetAddress)
	// defer connection.Close()
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
	fmt.Println("Sent addr", nextWorker.NetAddress, "to:", currentWorker.NetAddress)
}

// Dynamic does dynamic scheduling of the pipeline stages on the available nodes, with the aim of increasing
// throughput and memory utilization
func (schedule *Schedule) Dynamic(program string, masterAddress string) {
	for {
		time.Sleep(1 * time.Second)
		bottleneck := schedule.StageList.FindBottleneck()
		if bottleneck == -1 {
			fmt.Println("There is no bottleneck")
		} else {
			fmt.Println("Found a bottleneck at", bottleneck)
			schedule.scaleStage(bottleneck, program, masterAddress)
		}
	}
}
