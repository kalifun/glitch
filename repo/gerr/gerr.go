package gerr

import (
	"fmt"
	"time"
)

// Error represents a structured error with rich metadata and multi-language support
type Error struct {
	code       string
	key        string
	args       []interface{}
	metadata   map[string]interface{}
	cause      error
	time       time.Time
	errWrapper *ErrWrapper
}

// NewError creates a new error
func NewError(err ErrWrapper) *Error {
	return &Error{
		key:        err.Key,
		code:       err.Code,
		args:       nil,
		metadata:   make(map[string]interface{}),
		time:       time.Now(),
		errWrapper: &err,
	}
}

// New creates a new Error with the specified key and arguments
func New(key string, args ...interface{}) *Error {
	return &Error{
		key:      key,
		args:     args,
		metadata: make(map[string]interface{}),
		time:     time.Now(),
	}
}

// Wrap creates a new Error that wraps an existing error
func Wrap(cause error, key string, args ...interface{}) *Error {
	return &Error{
		key:      key,
		args:     args,
		cause:    cause,
		metadata: make(map[string]interface{}),
		time:     time.Now(),
	}
}

// Code sets the error code
func (e *Error) Code(code string) *Error {
	e.code = code
	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.errWrapper != nil {
		msg := ""
		if enMsg, ok := e.errWrapper.Messages["en"]; ok {
			msg = enMsg
		} else {
			// Use first available message
			for _, m := range e.errWrapper.Messages {
				msg = m
				break
			}
		}

		// Format with args if available
		if len(e.args) > 0 {
			return fmt.Sprintf(msg, e.args...)
		}
		return fmt.Sprintf("[%s] %s", e.code, msg)
	}

	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.key, e.cause)
	}
	return fmt.Sprintf("%s: %v", e.key, e.args)
}

// Args creates a new Error instance with arguments
func (e *Error) Args(args ...interface{}) *Error {
	if e.errWrapper == nil {
		newE := *e
		newE.args = append(newE.args, args...)
		return &newE
	}

	// create a new instance
	return &Error{
		key:        e.key,
		code:       e.code,
		args:       args,
		metadata:   make(map[string]interface{}),
		time:       time.Now(),
		errWrapper: e.errWrapper,
	}
}

// With adds metadata to the error
func (e *Error) With(key string, value interface{}) *Error {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}
	e.metadata[key] = value
	return e
}

// Unwrap returns the wrapped error for error unwrapping
func (e *Error) Unwrap() error {
	return e.cause
}

// GetKey returns the error key
func (e *Error) GetKey() string {
	return e.key
}

// GetCode returns the error code
func (e *Error) GetCode() string {
	return e.code
}

// GetArgs returns the error arguments
func (e *Error) GetArgs() []interface{} {
	return e.args
}

// GetMetadata returns the error metadata
func (e *Error) GetMetadata() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range e.metadata {
		result[k] = v
	}
	return result
}

// GetCause returns the wrapped error
func (e *Error) GetCause() error {
	return e.cause
}

// GetTime returns when the error was created
func (e *Error) GetTime() time.Time {
	return e.time
}

// GetErrWrapper returns the error definition
func (e *Error) GetErrWrapper() *ErrWrapper {
	return e.errWrapper
}

// GetMessages returns all messages
func (e *Error) GetMessages() map[string]string {
	if e.errWrapper != nil {
		return e.errWrapper.Messages
	}
	return nil
}

// GetMessage returns a message in the specified language
func (e *Error) GetMessage(lang string) string {
	if e.errWrapper == nil {
		return ""
	}

	if msg, ok := e.errWrapper.Messages[lang]; ok {
		return msg
	}
	// Fallback to English
	if msg, ok := e.errWrapper.Messages["en"]; ok {
		return msg
	}
	// Fallback to first available
	for _, msg := range e.errWrapper.Messages {
		return msg
	}
	return e.key
}

// Is implements error comparison for errors.Is
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	if t, ok := target.(*Error); ok {
		return e.key == t.key && e.code == t.code
	}

	return false
}

// ErrWrapper represents an error wrapper
type ErrWrapper struct {
	Key         string                 `json:"key" yaml:"key"`
	Code        string                 `json:"code" yaml:"code"`
	Messages    map[string]string      `json:"messages" yaml:"messages"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Category    string                 `json:"category,omitempty" yaml:"category,omitempty"`
	Severity    Severity               `json:"severity,omitempty" yaml:"severity,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Severity represents error severity levels
type Severity string

// Severity constants
const (
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)
