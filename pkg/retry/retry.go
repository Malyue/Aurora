package retry

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	// default retry count
	DefaultRetryCount = 3
	// default gap time
	DefaultGapTime = 3 * time.Second
)

var instance *RetryOption
var once sync.Once

type DelayRetryFunc func(context.Context, interface{}) (interface{}, bool, error)

// RetryOption
type RetryOption struct {
	// time in the retry
	GapTime time.Duration
	// counts of retry operation
	RetryCount int
	// method of retry
	RetryFunc DelayRetryFunc

	ctx context.Context
}

func NewRetryOption(ctx context.Context, gap time.Duration, retryCount int, func_ DelayRetryFunc) *RetryOption {
	once.Do(func() {
		instance = &RetryOption{
			GapTime:    gap,
			RetryCount: retryCount,
			RetryFunc:  func_,
			ctx:        ctx,
		}
	})
	return instance
}

func (r *RetryOption) Retry(ctx context.Context, req interface{}) (resp interface{}, err error) {
	if r.RetryFunc == nil {
		return nil, errors.New("the retry function is empty")
	}
	for i := 0; i < r.RetryCount; i++ {
		res, needRetry, errx := r.RetryFunc(ctx, req)
		if needRetry || err != nil {
			err = errx
			time.Sleep(r.GapTime)
			continue
		}
		resp = res
		break
	}
	return
}
