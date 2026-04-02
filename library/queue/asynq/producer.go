package asynq

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Send Create a task with task type and payload
// task := asynq.NewTask("send_welcome_email", map[string]interface{}{"user_id": 42})
// options:
// - asynq.MaxRetry
// - asynq.Queue - 指定队列名称
// - asynq.Timeout
// - asynq.Deadline
// - asynq.Unique - errors.Is(err, asynq.ErrDuplicateTask)
// - asynq.ProcessAt - 指定处理时间
// - asynq.ProcessIn - 指定延后时长
func (a *Asynq) Send(ctx context.Context, task *asynq.Task, options ...asynq.Option) (*asynq.TaskInfo, error) {
	return a.Client().EnqueueContext(ctx, task, options...)
}

func (a *Asynq) SendBy(ctx context.Context, typeName string, payload []byte, options ...asynq.Option) (*asynq.TaskInfo, error) {
	return a.Send(ctx, NewTask(typeName, payload), options...)
}

func (a *Asynq) SendJSON(ctx context.Context, typeName string, payload interface{}, options ...asynq.Option) (*asynq.TaskInfo, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return a.Send(ctx, NewTask(typeName, b), options...)
}
