package startup

import (
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
)

var ctx = execution_context.Get()
var logger = log.Get()

func Init() {
	// Add your initialization section here
}
