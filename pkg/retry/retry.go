package retry

import (
	"time"
)

const (
	// default retry count
	DefaultRetryCount = 3
	// default gap time
	DefaultGapTime = 3 * time.Second
)

//type DelayRetryFunc func(ctx context.Context, interface{}) (interface{}, bool, error)

type RetryOption struct {
	GapTime    time.Duration
	RetryCount int
	//RetryFunc
}
