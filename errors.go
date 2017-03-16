package dawdle

// ErrNoPendingTasks should be returned when
// a store doesn't currently have any ready tasks
type ErrNoPendingTasks struct{}

func (err ErrNoPendingTasks) Error() string {
	return "no ready tasks"
}
