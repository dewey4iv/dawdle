package dawdle

import (
	"github.com/dewey4iv/timestamps"
	uuid "github.com/satori/go.uuid"
)

// NewInvocation returns a new invocation
func NewInvocation(taskID string) Invocation {
	inv := Invocation{
		ID:     uuid.NewV4().String(),
		TaskID: taskID,
	}

	inv.Mark(timestamps.Created)

	return inv
}

// Invocation is created when a task is performed.
// It stores a resulting error and when the perform was attempted.
type Invocation struct {
	ID     string `json:"id"`
	TaskID string `json:"task_id"`
	Result bool   `json:"result"`
	Err    error  `json:"error"`

	timestamps.Timestamps
}
