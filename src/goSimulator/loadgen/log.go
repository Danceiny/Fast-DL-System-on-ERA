package loadgen

import (
	"goSimulator/logging"
)

var logger logging.Logger

func init() {
	logger = logging.NewSimpleLogger()
}
