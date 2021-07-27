package common

import "errors"

var (
	// ErrorNoAvailableJob there is no job in the topic
	ErrorNoAvailableJob = errors.New("no available job")
)
