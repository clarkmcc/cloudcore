package tasks

import "time"

// schedule returns the next time a Task should be run, if at all
type schedule interface {
	next(last time.Time) (time.Time, bool)
}

var _ schedule = RunOnceNow{}

// RunOnceNow is a schedule that runs once, immediately when next is called
// and then never runs again until the process is restarted.
type RunOnceNow struct{}

func (r RunOnceNow) next(last time.Time) (time.Time, bool) {
	if last.IsZero() {
		return time.Now(), true
	}
	return time.Time{}, false
}

// Interval is a schedule that runs at a fixed Interval
type Interval time.Duration

func (i Interval) next(last time.Time) (time.Time, bool) {
	if last.IsZero() {
		return time.Now(), true
	}
	return last.Add(time.Duration(i)), true
}
