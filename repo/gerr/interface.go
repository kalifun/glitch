package gerr

// Registry manages error definitions
type Registry interface {
	// Register registers an error definition
	Register(def ErrWrapper) error

	// Get retrieves an error definition by key
	Get(key string) (ErrWrapper, bool)

	// List returns all registered error definitions
	List() map[string]ErrWrapper

	// Remove removes an error definition
	Remove(key string) bool

	// Clear removes all error definitions
	Clear()
}
