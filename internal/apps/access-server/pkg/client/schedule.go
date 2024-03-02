package client

import "time"

type Scheduler struct {
	Interval time.Duration
}

func (s *Scheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}
