package timingWheel

import (
	"errors"
	"sync/atomic"
	"time"
	"unsafe"

	"Aurora/internal/apps/access-server/pkg/timingWheel/delayqueue"
)

type TimingWheel struct {
	tick        int64
	wheelSize   int64
	interval    int64
	currentTime int64
	buckets     []*bucket
	queue       *delayqueue.DelayQueue

	overflowWheel unsafe.Pointer

	exitChan  chan struct{}
	waitGroup waitGroupWrapper
}

func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be more that or equal to 1ms"))
	}

	startMs := timeToMs(time.Now().UTC())

	return newTimingWheel(tickMs, wheelSize, startMs, delayqueue.New(int(wheelSize)))
}

func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, queue *delayqueue.DelayQueue) *TimingWheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}

	return &TimingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * wheelSize,
		buckets:     buckets,
		queue:       queue,
		exitChan:    make(chan struct{}),
	}
}

func (tw *TimingWheel) add(t *Timer) bool {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if t.expiration < currentTime+tw.tick {
		return false
	}

	// if the expiration is in the current layer
	if t.expiration < currentTime+tw.interval {
		// get index in timingWheel
		virtualID := t.expiration / tw.tick
		b := tw.buckets[virtualID%tw.wheelSize]
		b.Add(t)
		// if is same time, return false
		if b.SetExpiration(virtualID * tw.tick) {
			tw.queue.Insert(b, b.Expiration())
		}
		return true
	}

	// the expiration is in the next layer
	overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
	if overflowWheel == nil {
		atomic.CompareAndSwapPointer(&tw.overflowWheel, nil, unsafe.Pointer(newTimingWheel(tw.interval, tw.wheelSize, currentTime, tw.queue)))
		overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
	}
	return (*TimingWheel)(overflowWheel).add(t)
}

func (tw *TimingWheel) Start() {
	tw.waitGroup.Wrap(func() {
		tw.queue.Poll(tw.exitChan, func() int64 {
			return timeToMs(time.Now().UTC())
		})
	})

	tw.waitGroup.Wrap(func() {
		for {
			select {
			case elem := <-tw.queue.C:
				b := elem.(*bucket)
				// move currentTime to expiration of bucket
				tw.advanceClock(b.Expiration())
				// Flush it
				b.Flush(tw.addOrRun)
			case <-tw.exitChan:
				return
			}
		}
	})
}

func (tw *TimingWheel) advanceClock(expiration int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= currentTime+tw.tick {
		currentTime = truncate(expiration, tw.tick)
		atomic.StoreInt64(&tw.currentTime, currentTime)

		// if it has the overflow wheel, also advance its clock
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advanceClock(currentTime)
		}
	}
}

func (tw *TimingWheel) addOrRun(t *Timer) {
	if !tw.add(t) {
		go t.task()
	}
}
