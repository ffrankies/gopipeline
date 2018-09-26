package master

import (
	"fmt"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/internal/common"
)

func MasterRun() {
	fmt.Println(gopipeline.Thing)
	common.ReadNodeList("/path/to/nodelist")
}
