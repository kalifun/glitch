package generator

import (
	"fmt"
	"strings"

	"github.com/kalifun/glitch/utils"
)

type ErrorDesc struct {
	Error []struct {
		Key         string            `yaml:"key"`
		Code        string            `yaml:"code"`
		Category    string            `yaml:"category,omitempty"`
		Severity    string            `yaml:"severity,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Message     map[string]string `yaml:"message,omitempty"`
	} `yaml:"error"`
}

func (e ErrorDesc) ToString() string {
	var definitions []string
	var variables []string
	var inits []string
	var helpers []string

	for i, v := range e.Error {
		if i == 0 {
			inits = append(inits, "func init() {\n")
		}

		// Handle messages
		messages := make(map[string]string)

		// First, copy all messages from the new format
		if v.Message != nil {
			for lang, msg := range v.Message {
				messages[lang] = msg
			}
		}

		// Build messages string
		var messageLines []string
		for lang, msg := range messages {
			escapedMsg := utils.EscapeString(msg)
			messageLines = append(messageLines, fmt.Sprintf(`		"%s": "%s",`, lang, escapedMsg))
		}
		messagesStr := strings.Join(messageLines, "\n")

		// Set defaults
		category := v.Category
		if category == "" {
			category = "general"
		}

		severity := v.Severity
		if severity == "" {
			severity = "SeverityError"
		} else {
			severity = "Severity" + utils.FirstUpper(severity)
		}

		description := v.Description
		if description == "" {
			description = utils.FirstUpper(strings.ReplaceAll(v.Key, "_", " "))
		}

		low := utils.FirstLower(v.Key)
		upper := utils.FirstUpper(utils.ToCamelCase(v.Key))

		// Generate error definition (private)
		defStr := fmt.Sprintf(DeclareErr, low+"Err", v.Key, v.Code, category, severity, messagesStr, description)
		definitions = append(definitions, defStr+"\n")

		// Generate public error variable that can be used directly
		varStr := fmt.Sprintf("// %s represents %s\nvar %s = gerr.NewError(%sErr)",
			upper, description, upper, low)
		variables = append(variables, varStr+"\n")

		// Generate registration call
		inits = append(inits, fmt.Sprintf(RegisterCall, low+"Err")+"\n")

		// Generate F suffix variable and helpers for errors that need formatting
		hasArgs := e.hasFormatArgs(messages)
		if hasArgs {
			// Generate F suffix variable with clear documentation
			formatVar := fmt.Sprintf("// %sF indicates this error requires format arguments\n// Usage: %sF.Args(\"%s\")\nvar %sF = %s",
				upper, upper, "ErrMessage", upper, upper)
			helpers = append(helpers, formatVar+"\n")

			// Generate helper function
			helperComment := fmt.Sprintf("// New%sWithArgs creates a %s error with arguments", upper, v.Key)
			helperFunc := fmt.Sprintf("func New%sWithArgs(args ...interface{}) *gerr.Error {\n\treturn %s.Args(args...)\n}",
				upper, upper)
			helpers = append(helpers, helperComment+"\n"+helperFunc+"\n")
		}

		// Generate helper function for creating instances with metadata
		helperComment2 := fmt.Sprintf("// New%sWithMetadata creates a %s error with metadata", upper, v.Key)
		helperFunc2 := fmt.Sprintf("func New%sWithMetadata(key string, value interface{}) *gerr.Error {\n\treturn %s.With(key, value)\n}",
			upper, upper)
		helpers = append(helpers, helperComment2+"\n"+helperFunc2+"\n")

		if i == len(e.Error)-1 {
			inits = append(inits, "}")
		}
	}

	result := strings.Join(definitions, "\n") + "\n\n" + strings.Join(variables, "\n") + "\n\n" + strings.Join(inits, "") + "\n\n" + strings.Join(helpers, "\n")
	return result
}

// hasFormatArgs checks if any message contains format arguments
func (e ErrorDesc) hasFormatArgs(messages map[string]string) bool {
	for _, msg := range messages {
		if ok := needsFormat(msg); ok {
			return ok
		}
	}
	return false
}
