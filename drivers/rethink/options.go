package rethink

// Option defines something that can build a Store
type Option interface {
	Apply(*Store) error
}
