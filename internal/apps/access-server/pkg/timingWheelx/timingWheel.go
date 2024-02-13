package timingWheelx

//import (
//	"math"
//	"sync"
//	"time"
//)
//
//var Executor = func(f func()) {
//	go f()
//}
//
//type Task struct {
//	offset int
//	s      *slot
//	at     time.Time
//
//	fn func()
//	C  chan struct{}
//}
//
//func (s *Task) TTL() int64 {
//	now := float64(time.Now().UnixNano())
//	at := float64(s.at.UnixNano())
//	return int64(math.Floor((at-now)/float64(time.Millisecond) + 1.0/2.0))
//}
//
//func (s *Task) call() {
//	if s.s == nil {
//		return
//	}
//	Executor(func() {
//		s.Cancel()
//		if s.fn != nil {
//			s.fn()
//		}
//		select {
//		case s.C <- struct{}{}:
//		default:
//
//		}
//	})
//}
//
//func (s *Task) Callback(f func()) {
//	s.fn = f
//}
//
//func (s *Task) Cancel() {
//	if s.s != nil {
//		s.s.remove(s)
//		s.s = nil
//	}
//}
//
//type slot struct {
//	index  int
//	next   *slot
//	len    int
//	values map[*Task]interface{}
//
//	sync.Mutex
//	circulate bool
//}
//
//func newSlot(circulate bool, len int) *slot {
//	var head *slot
//	var s *slot
//	for i := 0; i < len; i++ {
//		n := &slot{
//			index:     i,
//			len:       len,
//			values:    make(map[*Task]interface{}),
//			circulate: circulate,
//		}
//		if i == 0 {
//			head = n
//		} else {
//			s.next = n
//		}
//		s = n
//	}
//	s.next = head
//	return s
//}
//
//func (s *slot) put(offset int, v *Task) int {
//	if offset < 0 {
//		offset = 1
//	}
//	if !s.circulate && s.index == s.len && offset > 0 {
//		return offset
//	}
//	if offset == 0 {
//		s.Lock()
//		s.values[v] = nil
//		v.s = s
//		s.Unlock()
//		return 0
//	}
//	if offset >= s.len {
//		return offset - s.len
//	}
//	return s.next.put(offset-1, v)
//}
