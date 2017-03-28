package inmem

import (
	"fmt"
	"sync"

	"github.com/dewey4iv/dawdle"
	"github.com/dewey4iv/timestamps"
	uuid "github.com/satori/go.uuid"
)

// New takes a set of options and returns a store
func New(opts ...Option) (*Store, error) {
	var s Store

	if err := WithDefaults().Apply(&s); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt.Apply(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// Store is an in-memory implementation of a procrastinate.Store
type Store struct {
	tasks       dawdle.TaskStore
	invocations dawdle.InvocationStore
}

// Tasks returns the task store
func (s *Store) Tasks() dawdle.TaskStore {
	return s.tasks
}

// Invocations returns the invocation store
func (s *Store) Invocations() dawdle.InvocationStore {
	return s.invocations
}

// TaskStore holds tasks
type TaskStore struct {
	mux   sync.Mutex
	tasks []dawdle.Task
}

// ByID returns the task by ID
func (store *TaskStore) ByID(taskID string) (*dawdle.Task, error) {
	var t dawdle.Task

	for i := range store.tasks {
		if store.tasks[i].ID == taskID {
			return &store.tasks[i], nil
		}
	}

	return &t, fmt.Errorf("no task with id: %s", taskID)
}

// GetOne returns a single Task with status pending
func (store *TaskStore) GetOne() (*dawdle.Task, error) {
	for i := range store.tasks {
		if dawdle.IsReady(store.tasks[i]) {
			store.tasks[i].Status = dawdle.RunningTaskStatus
			store.tasks[i].Mark(timestamps.Updated)
			return &store.tasks[i], nil
		}
	}

	return &dawdle.Task{}, dawdle.ErrNoPendingTasks{}
}

// Pending pulls all pending tasks
func (store *TaskStore) Pending() ([]dawdle.Task, error) {
	return store.byStatus(dawdle.PendingTaskStatus)
}

func (store *TaskStore) byStatus(status dawdle.TaskStatus) ([]dawdle.Task, error) {
	var tasks []dawdle.Task

	for i := range store.tasks {
		if store.tasks[i].Status == status {
			tasks = append(tasks, store.tasks[i])
		}
	}

	return tasks, nil
}

// Save saves a task
func (store *TaskStore) Save(t dawdle.Task) error {
	if t.ID == "" {
		t.ID = uuid.NewV4().String()
	}

	var exists bool

	for i := range store.tasks {
		if store.tasks[i].ID == t.ID {
			store.tasks[i] = t
			exists = true
		}
	}

	if !exists {
		store.tasks = append(store.tasks, t)
	}

	return nil
}

// InvocationStore holds invocations
type InvocationStore struct {
	invocations []dawdle.Invocation
	mux         sync.Mutex
	byTaskID    map[string][]int
}

// AllByTask gets all invocations by taskID
func (store *InvocationStore) AllByTask(taskID string) ([]dawdle.Invocation, error) {
	var invocations []dawdle.Invocation

	indexes, exists := store.byTaskID[taskID]
	if !exists {
		return nil, fmt.Errorf("invalid task id: %s", taskID)
	}

	for _, i := range indexes {
		invocations = append(invocations, store.invocations[i])
	}

	return invocations, nil
}

// Save adds the invocation to the list
func (store *InvocationStore) Save(inv dawdle.Invocation) error {
	var exists bool

	inv.Mark(timestamps.Created)

	for i := range store.invocations {
		if inv.ID == store.invocations[i].ID {
			store.invocations[i] = inv
			exists = true
		}
	}

	if !exists {
		store.mux.Lock()
		store.invocations = append(store.invocations, inv)
		store.byTaskID[inv.TaskID] = append(store.byTaskID[inv.TaskID], len(store.invocations)-1)
		store.mux.Unlock()
	}

	return nil
}
