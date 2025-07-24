package gerr

import (
	"fmt"
	"sync"
)

// CacheRegistry is an in-memory implementation of Registry
type CacheRegistry struct {
	mu      sync.RWMutex
	wrapper map[string]ErrWrapper
}

// NewCacheRegistry creates a new memory registry
func NewCacheRegistry() *CacheRegistry {
	return &CacheRegistry{
		wrapper: make(map[string]ErrWrapper),
	}
}

// Register registers an error definition
func (r *CacheRegistry) Register(def ErrWrapper) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if def.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if def.Code == "" {
		return fmt.Errorf("code cannot be empty for key: %s", def.Key)
	}

	// Check for duplicate keys
	if _, exists := r.wrapper[def.Key]; exists {
		return fmt.Errorf("key '%s' already exists", def.Key)
	}

	// Set defaults
	if def.Category == "" {
		def.Category = "general"
	}
	if def.Severity == "" {
		def.Severity = SeverityError
	}
	if def.Messages == nil {
		def.Messages = make(map[string]string)
	}

	r.wrapper[def.Key] = def
	return nil
}

// Get retrieves an error definition by key
func (r *CacheRegistry) Get(key string) (ErrWrapper, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, exists := r.wrapper[key]
	return def, exists
}

// List returns all registered error definitions
func (r *CacheRegistry) List() map[string]ErrWrapper {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]ErrWrapper)
	for k, v := range r.wrapper {
		result[k] = v
	}
	return result
}

// Remove removes an error definition
func (r *CacheRegistry) Remove(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.wrapper[key]; exists {
		delete(r.wrapper, key)
		return true
	}
	return false
}

// Clear removes all error definitions
func (r *CacheRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.wrapper = make(map[string]ErrWrapper)
}

// Global registry instance
var globalRegistry Registry = NewCacheRegistry()

// Register registers an error wrapper globally
func Register(def ErrWrapper) error {
	return globalRegistry.Register(def)
}
