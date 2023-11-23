package tasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agent"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
	"time"
)

type Executor struct {
	registry *registry
	tomb     *tomb.Tomb
	logger   *zap.Logger
	client   *agent.Client
	db       agent.Database
}

func (e *Executor) Initialize() {
	for _, t := range e.registry.all() {
		e.maybeSchedule(t)
	}
}

func (e *Executor) schedule(task *Task, next time.Time) {
	e.tomb.Go(func() error {
		select {
		case <-e.tomb.Dying():
			e.logger.Warn("scheduled task cancelled", zap.String("task", task.Name))
			return nil
		case <-time.After(time.Until(next)):
			err := e.execute(task)
			if err != nil {
				e.logger.Error("task failed", zap.String("task", task.Name), zap.Error(err))
			}
			task.setLastRun(time.Now(), err == nil)
			e.maybeSchedule(task)
		}
		return nil
	})
}

func (e *Executor) maybeSchedule(task *Task) {
	n, ok := task.Schedule.next(task.getLastRun())
	if !ok {
		e.logger.Debug("task not scheduled to run again", zap.String("task", task.Name))
		return
	}
	e.logger.Info("scheduling task", zap.String("task", task.Name), zap.Time("next", n), zap.String("in", n.Sub(time.Now()).String()))
	e.schedule(task, n)
}

func (e *Executor) execute(task *Task) error {
	return task.Action(e.executionContext(context.Background(), task))
}

func NewExecutor(tomb *tomb.Tomb, db agent.Database, logger *zap.Logger, client *agent.Client) *Executor {
	e := Executor{
		registry: &DefaultRegistry,
		tomb:     tomb,
		logger:   logger.Named("executor"),
		client:   client,
		db:       db,
	}
	return &e
}
