package dawdle

// Store embodies
type Store interface {
	Tasks() TaskStore
	Invocations() InvocationStore
}

// TaskStore defines a generic store for handling Tasks
type TaskStore interface {
	GetOne() (*Task, error)
	Pending() ([]Task, error)
	ByID(taskID string) (*Task, error)
	Save(t Task) error
}

// InvocationStore defines a generic store for handling Invocationss
type InvocationStore interface {
	AllByTask(taskID string) ([]Invocation, error)
	Save(inv Invocation) error
}
