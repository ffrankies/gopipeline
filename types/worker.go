package types

// Worker represents a worker process running a particular stage on a particular node
type Worker struct {
	ID             string       // The ID of the worker
	Host           string       // The node on which the worker is running
	Stage          int          // The position of the stage it is running
	Address        string       // The address of the listener on this Worker
	PID            int          // The PID of the worker
	Stats          *WorkerStats // The performance statistics for this worker
	Exiting        bool         // Marks the worker as exiting, so it's not considered for communication
	updatedChannel chan int     // The channel for checking if the worker is updated
}

// NewWorker creates a new worker
func NewWorker(id string, host string, stage int) *Worker {
	worker := new(Worker)
	worker.ID = id
	worker.Host = host
	worker.Stage = stage
	worker.Address = ""
	worker.PID = -1
	worker.Stats = new(WorkerStats)
	worker.Exiting = false
	worker.updatedChannel = nil
	return worker
}

// UpdateInfo updates the worker's info
func (w *Worker) UpdateInfo(address string, pid int) {
	w.Address = address
	w.PID = pid
	if w.updatedChannel != nil {
		w.updatedChannel <- 1
	}
}

// WaitUntilUpdated waits until the worker's info has been updated
func (w *Worker) WaitUntilUpdated() {
	if w.PID != -1 && w.Address != "" {
		return
	}
	w.updatedChannel = make(chan int, 1)
	<-w.updatedChannel
	w.updatedChannel = nil
}
