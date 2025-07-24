package gerr

import "context"

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

// Handler processes errors
type Handler interface {
	// Handle processes an error and returns the result
	Handle(ctx context.Context, err *Error) interface{}

	// CanHandle returns true if this handler can process the error
	CanHandle(err *Error) bool
}

// Formatter formats errors for output
type Formatter interface {
	// Format formats an error for output
	Format(ctx context.Context, err *Error) interface{}

	// FormatWithOptions formats an error with specific options
	FormatWithOptions(ctx context.Context, err *Error, options FormatOptions) interface{}
}

// Localizer handles message localization
type Localizer interface {
	// Localize returns a localized message for the given error
	Localize(ctx context.Context, err *Error) string

	// LocalizeWithLanguage returns a localized message for a specific language
	LocalizeWithLanguage(language string, err *Error) string

	// GetSupportedLanguages returns all supported languages
	GetSupportedLanguages() []string
}
