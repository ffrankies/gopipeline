package common

// Identifiers for the messages sent by worker and master nodes
const (
	MsgStageInfo        int = 0
	MsgAddNextStageAddr int = 1
	MsgStageResult      int = 2
	MsgStageStats       int = 3
	MsgStartWorker      int = 4
	MsgBreakConnection  int = 5
	MsgNotifyExit       int = 6
	MsgEndExecution     int = 7
	MsgStartFirstStage  int = 8
)

// Performance log message types
const (
	PerfStartWorker string = "Worker_started         "
	PerfStartExec   string = "Stage_execution_started"
	PerfEndExec     string = "Stage_execution_ended  "
)
