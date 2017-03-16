package rethink

import (
	"github.com/dancannon/gorethink"
	"github.com/dewey4iv/dawdle"
)

// New takes a variadic list of options and returns an initialized Store
func New(opts ...Option) (*Store, error) {
	var s Store

	for _, opt := range opts {
		if err := opt.Apply(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// Store implments dawdle.Store
type Store struct {
	db          gorethink.Term
	session     *gorethink.Session
	tasks       dawdle.TaskStore
	invocations dawdle.InvocationStore
}

func (s *Store) Tasks() dawdle.TaskStore {
	return s.tasks
}

func (s *Store) Invocations() dawdle.InvocationStore {
	return s.invocations
}

// TaskStore
type TaskStore struct {
}

// GetOne
func (store *TaskStore) GetOne() (*dawdle.Task, error) {
	return nil, nil
}

// Pending
func (store *TaskStore) Pending() ([]dawdle.Task, error) {
	return nil, nil
}

// ByID
func (store *TaskStore) ByID(taskID string) (*dawdle.Task, error) {
	return nil, nil
}

// Save
func (store *TaskStore) Save(t dawdle.Task) error {
	return nil
}

// InvocationStore
type InvocationStore struct {
}

// AllByTask
func (store *InvocationStore) AllByTask(taskID string) ([]dawdle.Invocation, error) {
	return nil, nil
}

// Save
func (store *InvocationStore) Save(inv dawdle.Invocation) error {
	return nil
}
