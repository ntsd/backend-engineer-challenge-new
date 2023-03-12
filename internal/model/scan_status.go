package model

import (
	"database/sql/driver"
	"fmt"
)

// ScanStatus scan status
type ScanStatus string

const (
	// ScanStatusQueue scan status in queue
	ScanStatusQueue ScanStatus = "Queued"
	// ScanStatusInProgress scan status in progress
	ScanStatusInProgress ScanStatus = "In Progress"
	// ScanStatusSuccess scan status success
	ScanStatusSuccess ScanStatus = "Success"
	// ScanStatusFailure scan status failure
	ScanStatusFailure ScanStatus = "Failure"
)

// Scan database driver scan implement
func (p *ScanStatus) Scan(val interface{}) error {
	var str string
	switch v := val.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("scan `Finding` unsupported type: %T", v)
	}
	*p = ScanStatus(str)

	return nil
}

// Value database driver value implement
func (p ScanStatus) Value() (driver.Value, error) {
	return string(p), nil
}
