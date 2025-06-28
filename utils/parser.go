package utils

import (
	"errors"
	"strings"
)

// ParseCommand parses a raw command into name and arguments.
// Example: "SET key value" → name="SET", args=["key", "value"]
func ParseCommand(input string) (string, []string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil, errors.New("empty command")
	}
	parts := strings.Fields(input)
	cmd := strings.ToUpper(parts[0])
	args := parts[1:]
	return cmd, args, nil
}
