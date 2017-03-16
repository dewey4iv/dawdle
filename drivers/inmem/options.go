package inmem

// Option defines an interface that is used to build a Store
type Option interface {
	Apply(*Store) error
}

// WithDefaults is automatically called and initializes the in-memory maps
func WithDefaults() Option {
	return &withDefaults{}
}

type withDefaults struct{}

func (opt *withDefaults) Apply(s *Store) error {
	s.tasks = &TaskStore{}

	s.invocations = &InvocationStore{
		byTaskID: make(map[string][]int),
	}

	return nil
}
