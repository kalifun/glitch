package gerr

import (
	"context"
	"fmt"
)

// DefaultLocalizer is the default implementation of Localizer.
type DefaultLocalizer struct {
	registry Registry
}

// NewDefaultLocalizer creates a new default localizer
func NewDefaultLocalizer() *DefaultLocalizer {
	return &DefaultLocalizer{
		registry: globalRegistry,
	}
}

// SetRegistry sets the registry for the localizer
func (l *DefaultLocalizer) SetRegistry(registry Registry) {
	l.registry = registry
}

// Localize returns a localized message for the given error
func (l *DefaultLocalizer) Localize(ctx context.Context, err *Error) string {
	// Try to get language from context
	lang := "en" // default
	if v := ctx.Value("language"); v != nil {
		if langStr, ok := v.(string); ok {
			lang = langStr
		}
	}

	return l.LocalizeWithLanguage(lang, err)
}

// LocalizeWithLanguage returns a localized message for a specific language
func (l *DefaultLocalizer) LocalizeWithLanguage(language string, err *Error) string {
	if err.errWrapper != nil {
		// Use predefined error definition
		if msg, ok := err.errWrapper.Messages[language]; ok {
			if len(err.args) > 0 {
				return fmt.Sprintf(msg, err.args...)
			}
			return msg
		}

		// Fallback to English
		if msg, ok := err.errWrapper.Messages["en"]; ok {
			if len(err.args) > 0 {
				return fmt.Sprintf(msg, err.args...)
			}
			return msg
		}

		// Fallback to first available
		for _, msg := range err.errWrapper.Messages {
			if len(err.args) > 0 {
				return fmt.Sprintf(msg, err.args...)
			}
			return msg
		}
	}

	// Fallback to error string
	return err.Error()
}

// GetSupportedLanguages returns all supported languages
func (l *DefaultLocalizer) GetSupportedLanguages() []string {
	languages := make(map[string]bool)

	for _, def := range l.registry.List() {
		for lang := range def.Messages {
			languages[lang] = true
		}
	}

	result := make([]string, 0, len(languages))
	for lang := range languages {
		result = append(result, lang)
	}

	return result
}
