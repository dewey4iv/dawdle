package processor

import (
	"fmt"

	"github.com/dewey4iv/dawdle"
)

// NewRegistrar returns the default Registrar
func NewRegistrar() *Registrar {
	return &Registrar{
		make(map[string]dawdle.Converter),
	}
}

// Registrar holds a map of Converters by their given name
type Registrar struct {
	converters map[string]dawdle.Converter
}

// Register adds a converter
func (r *Registrar) Register(name string, c dawdle.Converter) error {
	r.converters[name] = c

	return nil
}

// Fetch gets a converter
func (r *Registrar) Fetch(name string) (dawdle.Converter, error) {
	if converter, exists := r.converters[name]; exists && converter != nil {
		return converter, nil
	}

	return nil, fmt.Errorf("no converter with name %s found", name)
}
