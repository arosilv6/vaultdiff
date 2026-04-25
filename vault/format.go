package vault

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OutputFormat represents the output format for vault data.
type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatYAML  OutputFormat = "yaml"
)

// Formatter formats vault secret data into different output representations.
type Formatter struct {
	format OutputFormat
}

// NewFormatter creates a new Formatter with the given output format.
func NewFormatter(format OutputFormat) *Formatter {
	return &Formatter{format: format}
}

// FormatData renders the provided key-value data according to the formatter's output format.
func (f *Formatter) FormatData(data map[string]interface{}) (string, error) {
	if data == nil {
		return "", fmt.Errorf("data is nil")
	}
	switch f.format {
	case FormatJSON:
		return f.toJSON(data)
	case FormatTable:
		return f.toTable(data), nil
	case FormatYAML:
		return f.toYAML(data), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", f.format)
	}
}

func (f *Formatter) toJSON(data map[string]interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}

func (f *Formatter) toTable(data map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %s\n", "KEY", "VALUE"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for k, v := range data {
		sb.WriteString(fmt.Sprintf("%-30s %v\n", k, v))
	}
	return sb.String()
}

func (f *Formatter) toYAML(data map[string]interface{}) string {
	var sb strings.Builder
	for k, v := range data {
		sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}
	return sb.String()
}
