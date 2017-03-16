package dawdle

// Performer defines something that can be run.
// Most of these should embed a Task
type Performer interface {
	Perform() error
}
