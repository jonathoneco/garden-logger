package state

type StateError struct {
	Op   string
	Path string
	Err  error
}

func (e *StateError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

func (e *StateError) Unwrap() error { return e.Err }

// Timeout reports whether this error represents a timeout.
func (e *StateError) Timeout() bool {
	t, ok := e.Err.(interface{ Timeout() bool })
	return ok && t.Timeout()
}
