package asynq

import (
	"context"
	"errors"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/hibiken/asynq"
)

var queue *Asynq
var qonce once.Once

func initialize() {
	opt, err := RedisOptFromCache()
	if err != nil {
		log.Error(err)
		return
	}
	queue, err = New(opt)
	if err != nil {
		log.Error(err)
		return
	}
}

func Default() *Asynq {
	qonce.Do(initialize)
	return queue
}

// ------- producer -------

func Send(task *Task, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return queue.Send(task, options...)
}

func SendBy(typeName string, payload []byte, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return q.SendBy(typeName, payload, options...)
}

func SendJSON(typeName string, payload interface{}, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return q.SendJSON(typeName, payload, options...)
}

// ------- consumer -------

func WorkerHandleFunc(pattern string, handler func(context.Context, *asynq.Task) error) error {
	q := Default()
	if q == nil {
		return errors.New(`failed to initialize queue`)
	}

	q.HandleFunc(pattern, handler)

	//q.StartWorker()
	return nil
}

func WorkerHandle(pattern string, handler asynq.Handler) error {
	q := Default()
	if q == nil {
		return errors.New(`failed to initialize queue`)
	}

	q.Handle(pattern, handler)

	//q.StartWorker()
	return nil
}

func StartWorker(configs ...*asynq.Config) error {
	q := Default()
	if q == nil {
		return errors.New(`failed to initialize queue`)
	}

	return q.StartWorker(configs...)
}
