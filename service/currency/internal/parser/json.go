package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ParseErr = errors.New("parse-err")

// Reads and unmarshals content from buffer.
func ParseJSON(read io.Reader, target any) error {
	body, err := io.ReadAll(read)
	if err != nil {
		return fmt.Errorf("%w: %w", ParseErr, err)
	}
	if err = json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("%w: %w", ParseErr, err)
	}

	return nil
}
