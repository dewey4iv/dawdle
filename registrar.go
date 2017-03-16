package dawdle

// Registrar acts as a map for the name of a function and the actual function
// Tihs is used when pulling Tasks out of a database and converting them to Performers
type Registrar interface {
	Register(name string, converter Converter) error
	Fetch(name string) (Converter, error)
}

// Converter takes a task and returns a performer
type Converter func(t Task) (Performer, error)
