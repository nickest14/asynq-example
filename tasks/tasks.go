package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	// TypeHeartBeat is a name of the task type
	TypeHeartBeat = "heartbeat"
)

// HeartBeatTask payload.
func HeartBeatTask(id int) *asynq.Task {
	// Specify task payload.
	payload := map[string]interface{}{
		"user_id": id, // set user ID
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeHeartBeat, payload)
}

// HandlHeartBeatTask handler.
func HandlHeartBeatTask(c context.Context, t *asynq.Task) error {
	// Get user ID from given task.
	id, err := t.Payload.GetInt("user_id")
	if err != nil {
		return err
	}
	time := time.Now()
	fmt.Printf("heart beat user_id %d, now is %v\n", id, time)
	return nil
}
