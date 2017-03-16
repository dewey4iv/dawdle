package postgres

import "github.com/dewey4iv/sqlgo"

// Option defines and interface that builds a store
type Option interface {
	Apply(*Store) error
}

func WithDefaults() Option {
	return &withDefaults{}
}

type withDefaults struct{}

func (opt *withDefaults) Apply(s *Store) error {
	db, err := sqlgo.New(
		sqlgo.WithDriver("postgres"),
		sqlgo.WithHostPort("127.0.0.1", "5432"),
		sqlgo.WithCreds("user", "password", "dawdle"),
	)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}
