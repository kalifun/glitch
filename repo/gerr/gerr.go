package gerr

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
