package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSONReadOptions controls behavior for missing or empty files.
type JSONReadOptions struct {
	AllowMissing bool
	AllowEmpty   bool
}

// ReadJSON reads a JSON file at path into T and annotates errors with label context.
func ReadJSON[T any](path string, label string, opts JSONReadOptions) (T, error) {
	var zero T

	target := strings.TrimSpace(label)
	if target == "" {
		target = "file"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && opts.AllowMissing {
			return zero, nil
		}
		return zero, fmt.Errorf("failed to read %s %q: %w", target, path, err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		if opts.AllowEmpty {
			return zero, nil
		}
		return zero, fmt.Errorf("failed to parse %s %q: file is empty", target, path)
	}

	var decoded T
	if err := json.Unmarshal(data, &decoded); err != nil {
		return zero, fmt.Errorf("failed to parse %s %q: %w", target, path, err)
	}

	return decoded, nil
}
