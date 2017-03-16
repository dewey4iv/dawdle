package postgres

import (
	"fmt"
	"sync"

	"github.com/dewey4iv/dawdle"
	"github.com/dewey4iv/sqlgo"
)

func New(opts ...Option) (*Store, error) {
	var s Store

	for _, opt := range opts {
		if err := opt.Apply(&s); err != nil {
			return nil, err
		}
	}

	s.tasks = &TaskStore{
		db: s.db,
	}
	s.invocations = &InvocationStore{
		db: s.db,
	}

	return &s, nil
}

type Store struct {
	db          sqlgo.DB
	tasks       dawdle.TaskStore
	invocations dawdle.InvocationStore
}

func (s *Store) Tasks() dawdle.TaskStore {
	return s.tasks
}

func (s *Store) Invocations() dawdle.InvocationStore {
	return s.invocations
}

type TaskStore struct {
	db  sqlgo.DB
	mux sync.Mutex
}

// GetOne
func (store *TaskStore) GetOne() (*dawdle.Task, error) {
	store.mux.Lock()
	defer store.mux.Unlock()

	rows, err := store.db.Query(`
		SELECT * FROM tasks
		WHERE status = 'pending'
		ORDER BY run_at
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var task dawdle.Task
	for rows.Next() {
		if err := rows.Scan(&task); err != nil {
			return nil, err
		}

		if dawdle.IsReady(task) {
			return &task, nil
		}

		task.Status = dawdle.RunningTaskStatus

		if err := store.Save(task); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("no tasks ready to run")
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
	var exists int

	if err := store.db.QueryRow(`
		SELECT COUNT(true) FROM tasks
		WHERE id=$1
	`, t.ID).Scan(&exists); err != nil {
		return err
	}

	if exists > 0 {
		if _, err := store.db.Exec(`
			UPDATE tasks SET
				func_name = $2,
				status = $4,
				run_at = $5,
				created = $6,
				updated = $7
			WHERE id = $1;
		`, t.ID, t.FuncName, t.Status, t.Args, t.RunAt, t.Created, t.Updated); err != nil {
			return err
		}
	} else {
		if _, err := store.db.Exec(`
			INSERT INTO tasks (
				id,
				func_name,
				args,
				status,
				run_at,
				created,
				updated
			) VALUES ($1, $2, $3, $4, $5, $6, $7);
		`, t.ID, t.FuncName, t.Status, t.Args, t.RunAt, t.Created, t.Updated); err != nil {
			return err
		}
	}

	return nil
}

type InvocationStore struct {
	db sqlgo.DB
}

// AllByTask
func (store *InvocationStore) AllByTask(taskID string) ([]dawdle.Invocation, error) {
	rows, err := store.db.Query(`
		SELECT * FROM invocations
		WHERE task_id = $1;
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invocations []dawdle.Invocation
	var invocation dawdle.Invocation
	for rows.Next() {
		if err := rows.Scan(&invocation); err != nil {
			return nil, err
		}

		invocations = append(invocations, invocation)
	}

	return invocations, nil
}

// Save
func (store *InvocationStore) Save(inv dawdle.Invocation) error {
	var exists int

	if err := store.db.QueryRow(`
		SELECT COUNT(true) FROM invocations
		WHERE id=$1
	`, inv.ID).Scan(&exists); err != nil {
		return err
	}

	if exists > 0 {
		if _, err := store.db.Exec(`
			UPDATE invocations SET
				task_id = $2,
				result = $3,
				error = $4,
				created = $5,
				updated = $6
			WHERE id = $1;
		`, inv.ID, inv.TaskID, inv.Result, inv.Err.Error(), inv.Created, inv.Updated); err != nil {
			return err
		}
	} else {
		if _, err := store.db.Exec(`
			INSERT INTO invocations (
				id,
				task_id,
				result,
				error,
				created,
				updated
			) VALUES ($1, $2, $3, $4, $5, $6, $7);
		`, inv.ID, inv.TaskID, inv.Result, inv.Err.Error(), inv.Created, inv.Updated); err != nil {
			return err
		}
	}

	return nil
}
