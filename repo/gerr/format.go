package gerr

import (
	"context"
	"encoding/json"
	"fmt"
)

// FormatOptions controls error formatting
type FormatOptions struct {
	IncludeMetadata bool       `json:"include_metadata"`
	IncludeCause    bool       `json:"include_cause"`
	Language        string     `json:"language,omitempty"`
	Format          FormatType `json:"format,omitempty"` // "json", "text", "structured"
}

type FormatType string

const (
	FormatTypeJSON       FormatType = "json"
	FormatTypeText       FormatType = "text"
	FormatTypeStructured FormatType = "structured"
)

// DefaultFormatter is the default implementation of Formatter
type DefaultFormatter struct {
	localizer Localizer
}

// NewDefaultFormatter creates a new default formatter
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		localizer: NewDefaultLocalizer(),
	}
}

// Format formats an error for output
func (f *DefaultFormatter) Format(ctx context.Context, err *Error) interface{} {
	return f.FormatWithOptions(ctx, err, FormatOptions{
		IncludeMetadata: true,
		IncludeCause:    true,
		Language:        "en",
		Format:          FormatTypeStructured,
	})
}

// FormatWithOptions formats an error with specific options
func (f *DefaultFormatter) FormatWithOptions(ctx context.Context, err *Error, options FormatOptions) interface{} {
	switch options.Format {
	case FormatTypeJSON:
		return f.formatJSON(ctx, err, options)
	case FormatTypeText:
		return f.formatText(ctx, err, options)
	default:
		return f.formatStructured(ctx, err, options)
	}
}

// formatStructured formats error as structured data (map)
func (f *DefaultFormatter) formatStructured(ctx context.Context, err *Error, options FormatOptions) map[string]interface{} {
	result := map[string]interface{}{
		"key":     err.key,
		"code":    err.code,
		"message": f.localizer.LocalizeWithLanguage(options.Language, err),
		"time":    err.time,
	}

	if options.IncludeMetadata && len(err.metadata) > 0 {
		result["metadata"] = err.metadata
	}

	if options.IncludeCause && err.cause != nil {
		if causeErr, ok := err.cause.(*Error); ok {
			result["cause"] = f.formatStructured(ctx, causeErr, options)
		} else {
			result["cause"] = err.cause.Error()
		}
	}

	if err.errWrapper != nil {
		result["category"] = err.errWrapper.Category
		result["severity"] = string(err.errWrapper.Severity)
	}

	return result
}

// formatJSON formats error as JSON string
func (f *DefaultFormatter) formatJSON(ctx context.Context, err *Error, options FormatOptions) string {
	structured := f.formatStructured(ctx, err, options)
	data, _ := json.MarshalIndent(structured, "", "  ")
	return string(data)
}

// formatText formats error as plain text
func (f *DefaultFormatter) formatText(ctx context.Context, err *Error, options FormatOptions) string {
	var parts []string

	// Basic error info
	if f.localizer != nil {
		if options.Language != "" {
			parts = append(parts, f.localizer.LocalizeWithLanguage(options.Language, err))
		} else {
			parts = append(parts, f.localizer.Localize(ctx, err))
		}
	} else {
		parts = append(parts, err.Error())
	}

	// Add code if available
	if err.code != "" {
		parts = append(parts, fmt.Sprintf("Code: %s", err.code))
	}

	// Add metadata
	if options.IncludeMetadata && len(err.metadata) > 0 {
		parts = append(parts, fmt.Sprintf("Metadata: %v", err.metadata))
	}

	// Add cause
	if options.IncludeCause && err.cause != nil {
		if causeErr, ok := err.cause.(*Error); ok {
			parts = append(parts, fmt.Sprintf("Cause: %s", f.formatText(ctx, causeErr, options)))
		} else {
			parts = append(parts, fmt.Sprintf("Cause: %s", err.cause.Error()))
		}
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "\n"
		}
		result += part
	}

	return result
}
