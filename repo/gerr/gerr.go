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
