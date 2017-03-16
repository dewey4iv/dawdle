package dawdle

// Scheduler is anything that can schedule a task
type Scheduler interface {
	Schedule(t Task) error
}
