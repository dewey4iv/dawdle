package scheduler

import "github.com/dewey4iv/dawdle"

// New returns a new Scheduler
func New(opts ...Option) (*Scheduler, error) {
	var s Scheduler

	for _, opt := range opts {
		if err := opt.Apply(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// Scheduler is the implementation of a Scheduler
type Scheduler struct {
	store dawdle.Store
}

// Schedule saves a task to the underlying storage
func (s *Scheduler) Schedule(t dawdle.Task) error {
	return s.store.Tasks().Save(t)
}

// Option build a Scheduler
type Option interface {
	Apply(*Scheduler) error
}

// WithStore sets the store
func WithStore(store dawdle.Store) Option {
	return &withStore{store}
}

type withStore struct {
	store dawdle.Store
}

func (opt *withStore) Apply(s *Scheduler) error {
	s.store = opt.store

	return nil
}
