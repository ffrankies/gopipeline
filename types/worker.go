package types

// Worker represents a worker process running a particular stage on a particular node
type Worker struct {
	ID      string       // The ID of the worker
	Host    string       // The node on which the worker is running
	Stage   int          // The position of the stage it is running
	Address string       // The address of the listener on this Worker
	PID     int          // The PID of the worker
	Stats   *WorkerStats // The performance statistics for this worker
	Exiting bool         // Marks the worker as exiting, so it's not considered for communication
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
	return worker
}
