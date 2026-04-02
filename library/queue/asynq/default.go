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

func Send(ctx context.Context, task *Task, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return q.Send(ctx, task, options...)
}

func SendBy(ctx context.Context, typeName string, payload []byte, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return q.SendBy(ctx, typeName, payload, options...)
}

func SendJSON(ctx context.Context, typeName string, payload interface{}, options ...asynq.Option) (*asynq.TaskInfo, error) {
	q := Default()
	if q == nil {
		return nil, errors.New(`failed to initialize queue`)
	}
	return q.SendJSON(ctx, typeName, payload, options...)
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
