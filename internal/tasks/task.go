package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var DefaultRegistry = registry{}

type Task struct {
	Name     string
	Schedule schedule
	Action   func(ctx context.Context) error

	lock           sync.Mutex
	lastSuccessRun time.Time
	lastRun        time.Time
}

func (t *Task) setLastRun(tim time.Time, success bool) {
	t.lock.Lock()
	t.lastRun = tim
	if success {
		t.lastSuccessRun = tim
	}
	t.lock.Unlock()
}

func (t *Task) getLastRun() time.Time {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.lastRun
}

func (t *Task) getLastSuccessRun() time.Time {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.lastSuccessRun
}

type registry map[string]*Task

func (r *registry) Register(t *Task) {
	if _, ok := (*r)[t.Name]; ok {
		panic(fmt.Sprintf("Task %q already registered", t.Name))
	}
	(*r)[t.Name] = t
}

func (r *registry) get(name string) (*Task, bool) {
	v, ok := (*r)[name]
	return v, ok
}

func (r *registry) all() []*Task {
	var out []*Task
	for _, v := range *r {
		out = append(out, v)
	}
	return out
}
