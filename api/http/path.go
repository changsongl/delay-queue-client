package http

import (
	"fmt"
	"time"
)

const (
	addPathFormat = "/topic/%s/job"
	endPathFormat = "/topic/%s/job/%s"
	popPathFormat = "/topic/%s/job?timeout=%d"
)

// add job path
func addJobPath(topic string) string {
	return fmt.Sprintf(addPathFormat, topic)
}

// finish job path
func finishJobPath(topic, id string) string {
	return fmt.Sprintf(endPathFormat, topic, id)
}

// pop job path
func popJobPath(topic string, timeout time.Duration) string {
	return fmt.Sprintf(popPathFormat, topic, timeout/time.Second)
}
