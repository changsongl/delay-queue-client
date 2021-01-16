package client

import (
	"errors"
	"time"
)

type (
	jobTopic string
	jobId    string
	jobDelay time.Duration
	jobTTR   time.Duration
	jobBody  string
)

type JobField interface {
	IsValid() error
}

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

func (id jobId) IsValid() error {
	if id == "" {
		return errors.New("jobId is empty")
	}

	l := len(id)
	if l > 50 {
		return errors.New("jobId contains over 50 characters")
	}

	return nil
}
