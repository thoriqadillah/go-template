// Code generated by BobGen psql v0.30.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql/driver"
	"fmt"
)

// Enum values for RiverJobState
const (
	RiverJobStateAvailable RiverJobState = "available"
	RiverJobStateCancelled RiverJobState = "cancelled"
	RiverJobStateCompleted RiverJobState = "completed"
	RiverJobStateDiscarded RiverJobState = "discarded"
	RiverJobStatePending   RiverJobState = "pending"
	RiverJobStateRetryable RiverJobState = "retryable"
	RiverJobStateRunning   RiverJobState = "running"
	RiverJobStateScheduled RiverJobState = "scheduled"
)

func AllRiverJobState() []RiverJobState {
	return []RiverJobState{
		RiverJobStateAvailable,
		RiverJobStateCancelled,
		RiverJobStateCompleted,
		RiverJobStateDiscarded,
		RiverJobStatePending,
		RiverJobStateRetryable,
		RiverJobStateRunning,
		RiverJobStateScheduled,
	}
}

type RiverJobState string

func (e RiverJobState) String() string {
	return string(e)
}

func (e RiverJobState) Valid() bool {
	switch e {
	case RiverJobStateAvailable,
		RiverJobStateCancelled,
		RiverJobStateCompleted,
		RiverJobStateDiscarded,
		RiverJobStatePending,
		RiverJobStateRetryable,
		RiverJobStateRunning,
		RiverJobStateScheduled:
		return true
	default:
		return false
	}
}

func (e RiverJobState) MarshalText() ([]byte, error) {
	return []byte(e), nil
}

func (e *RiverJobState) UnmarshalText(text []byte) error {
	return e.Scan(text)
}

func (e RiverJobState) MarshalBinary() ([]byte, error) {
	return []byte(e), nil
}

func (e *RiverJobState) UnmarshalBinary(data []byte) error {
	return e.Scan(data)
}

func (e RiverJobState) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *RiverJobState) Scan(value any) error {
	switch x := value.(type) {
	case string:
		*e = RiverJobState(x)
	case []byte:
		*e = RiverJobState(x)
	case nil:
		return fmt.Errorf("cannot nil into RiverJobState")
	default:
		return fmt.Errorf("cannot scan type %T: %v", value, value)
	}

	if !e.Valid() {
		return fmt.Errorf("invalid RiverJobState value: %s", *e)
	}

	return nil
}
