package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Findings []Finding

// Finding repository scan finding object
type Finding struct {
	Path        string `json:"path"`
	Line        int    `json:"line"`
	Description string `json:"description"`
}

// Scan database driver scan implement
func (f *Findings) Scan(val interface{}) error {
	var bytes []byte

	switch v := val.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("scan `Finding` unsupported type: %T", v)
	}

	if err := json.Unmarshal(bytes, f); err != nil {
		return fmt.Errorf("error to unmarshal `Finding`: %w", err)
	}

	return nil
}

// Value database driver value implement
func (f Findings) Value() (driver.Value, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("error to marshal `Finding` value: %w", err)
	}

	return bytes, nil
}
