package job

import (
	"errors"
	"time"
)

type (
	jobTopic string
	jobID    string
	jobDelay time.Duration
	jobTTR   time.Duration
	jobBody  string
)

// Field job field interface
type Field interface {
	IsValid() error
}

// IsValid check it is valid
func (t jobTopic) IsValid() error {
	if t == "" {
		return errors.New("jobTopic is empty")
	}

	l := len(t)
	if l > 50 {
		return errors.New("jobTopic contains over 50 characters")
	}

	return nil
}

// IsValid check is valid
func (id jobID) IsValid() error {
	if id == "" {
		return errors.New("jobID is empty")
	}

	l := len(id)
	if l > 50 {
		return errors.New("jobID contains over 50 characters")
	}

	return nil
}
