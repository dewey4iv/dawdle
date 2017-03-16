package dawdle

import (
	"time"

	"github.com/dewey4iv/timestamps"
	uuid "github.com/satori/go.uuid"
)

// EmptyTask defines an empty task.
var EmptyTask = Task{}

// NewTask takes the func name, args, and variadic list of options to create a new Task
func NewTask(fn string, args []byte, opts ...TaskOpt) Task {
	t := Task{
		ID:       uuid.NewV4().String(),
		FuncName: fn,
		Args:     args,
		Status:   PendingTaskStatus,
	}

	t.Mark(timestamps.Created)

	for _, opt := range opts {
		opt.Apply(&t)
	}

	return t
}

// Task is the basic storage
type Task struct {
	ID       string     `json:"id"`
	FuncName string     `json:"func_name"`
	Args     []byte     `json:"args"`
	Status   TaskStatus `json:"status"`
	RunAt    *time.Time `json:"run_at,omitempty"`

	timestamps.Timestamps
}

// TaskStatus is an enum type
type TaskStatus string

// Enum values for TaskStatus
const (
	PendingTaskStatus TaskStatus = "Pending"
	RunningTaskStatus TaskStatus = "Running"
	PassedTaskStatus  TaskStatus = "Passed"
	FailedTaskStatus  TaskStatus = "Failed"
)

// TaskOpt helps build a Task
type TaskOpt interface {
	Apply(*Task)
}

// WithDelay adds a time delay to the task being run
func WithDelay(d time.Duration) TaskOpt {
	return &withDelay{d, time.Now()}
}

type withDelay struct {
	d time.Duration
	t time.Time
}

func (opt *withDelay) Apply(t *Task) {
	runAt := opt.t.Add(opt.d)

	t.RunAt = &runAt
}

// IsReady determines if a task is ready to be performed
func IsReady(t Task) bool {
	if t.Status != PendingTaskStatus {
		return false
	}

	if time.Since(*t.RunAt) < 0 {
		return false
	}

	return true
}
