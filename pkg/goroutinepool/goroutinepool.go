package goroutinepool

import (
	"errors"
	"runtime/debug"
	"sync"
	"time"
)

var (
	NoMoreWorkerErr = errors.New("no more worker,the popl is full")
	TimeoutErr      = errors.New("time out")
)

// GoroutinePool
type GoroutinePool struct {
	allWorkers []*worker
	// cap == len(allWorkers)
	cap int
	// wait all workers stopped
	sync.WaitGroup
	// protect 'running' and 'workers'
	sync.RWMutex
	workers chan *worker
	running bool
}

// worker
type worker struct {
	pool *GoroutinePool
	// function chan
	job chan func()
	// stop chan
	stop chan struct{}
}

// New returns a GoroutinePool With Cap
func New(cap int) *GoroutinePool {
	pool := &GoroutinePool{
		workers: make(chan *worker, cap),
		cap:     cap,
		running: false,
	}
	return pool
}

// Start if the `p` is running,return
// otherwise create the worker node and run it with a goroutine
func (p *GoroutinePool) Start() {
	p.Lock()
	defer p.Unlock()

	// if p has been started,return
	if p.running {
		return
	}

	if p.allWorkers == nil {
		for i := 0; i < p.cap; i++ {
			w := &worker{pool: p, job: make(chan func()), stop: make(chan struct{}, 1)}
			p.allWorkers = append(p.allWorkers, w)
			w.pool.Add(1)
			go w.run()
		}
	} else {
		p.workers = make(chan *worker, p.cap)
		for _, w := range p.allWorkers {
			w.pool.Add(1)
			go w.run()
		}
	}
	p.running = true
}

// Stop Stop the GoroutinePool
// use channel to send the exit signal to the worker
func (p *GoroutinePool) Stop() {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return
	}
	for _, w := range p.allWorkers {
		w.stop <- struct{}{}
	}

	// block the step util all the worker stopped
	p.Wait()

	close(p.workers)
	p.running = false
}

// RunningStatus Get the goroutinePool Status
func (p *GoroutinePool) RunningStatus() bool {
	p.Lock()
	defer p.Unlock()

	return p.running
}

// Go Select a worker to exec method
// if the worker is not exist currently
// return an error
func (p *GoroutinePool) Go(f func()) error {
	p.RLock()
	defer p.RUnlock()

	if !p.running {
		panic("not running ")
	}

	select {
	case worker := <-p.workers:
		worker.job <- f
	default:
		return NoMoreWorkerErr
	}
	return nil
}

// GoWithTimeout Select a worker to exec method with Timeout
func (p *GoroutinePool) GoWithTimeout(f func(), timeout time.Duration) error {
	p.RLock()
	defer p.RUnlock()

	if !p.running {
		panic("not running")
	}

	select {
	case worker := <-p.workers:
		worker.job <- f
	case <-time.After(timeout):
		return TimeoutErr
	}

	return nil
}

// MustGo select a worker until the p.workers is not empty
// and set the method into the worker to exec
func (p *GoroutinePool) MustGo(f func()) {
	p.RLock()
	defer p.RUnlock()

	if !p.running {
		panic("not running")
	}

	select {
	case worker := <-p.workers:
		worker.job <- f
	}
}

// Statistics return [<IDLE-worker-num>, <total-worker-num>]
func (p *GoroutinePool) Statistics() [2]int {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return [2]int{0, 0}
	}
	return [2]int{len(p.workers), len(p.allWorkers)}
}

// addIdleWorker add the worker into p.workers
func (p *GoroutinePool) addIdleWorker(w *worker) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	p.workers <- w
}

// run worker run method
// add the worker in p.workers(is a channel)
// and use the p.workers to control the cap of the goroutinePool
// for - select to wait the signal
// if is a job,exec it
// otherwise,if get a stop signal,stop the worker
func (w *worker) run() {
	for {
		w.pool.addIdleWorker(w)
		select {
		case job := <-w.job:
			func() {
				defer func() {
					if err := recover(); err != nil {
						debug.PrintStack()
					}
				}()
				job()
			}()
		case <-w.stop:
			w.pool.Done()
			return
		}
	}
}
