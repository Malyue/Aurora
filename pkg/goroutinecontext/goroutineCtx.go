package goroutinecontext

import (
	"context"
	"runtime"
	"sync"

	"github.com/go-eden/routine"
)

const LocaleNameContextKey = "locale_name_context_key"

const bucketsSize = 128

const arm64 = "arm64"

type (
	contextBucket struct {
		lock sync.RWMutex
		data map[int64]context.Context
	}
	contextBuckets struct {
		buckets [bucketsSize]*contextBucket
	}
)

var goroutineContext contextBuckets

func init() {
	for i := range goroutineContext.buckets {
		goroutineContext.buckets[i] = &contextBucket{
			data: make(map[int64]context.Context),
		}
	}
}

// GetContext .
func GetContext() context.Context {
	// mac system use goroutine_context.GetContext() will panic
	if runtime.GOARCH == arm64 {
		return context.Background()
	}

	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.RLock()
	ctx := bucket.data[goid]
	bucket.lock.RUnlock()
	return ctx
}

// SetContext .
func SetContext(ctx context.Context) {
	// mac system use goroutine_context.GetContext() will panic
	if runtime.GOARCH == arm64 {
		return
	}

	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	bucket.data[goid] = ctx
}

// ClearContext .
func ClearContext() {
	// mac system use goroutine_context.GetContext() will panic
	if runtime.GOARCH == arm64 {
		return
	}

	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	delete(bucket.data, goid)
}
