package http

import "fmt"

const (
	addPathFormat = "/topic/%s/job"
	endPathFormat = "/topic/%s/job/%s"
	popPathFormat = "/topic/%s/job"
)

func addJobPath(topic string) string {
	return fmt.Sprintf(addPathFormat, topic)
}

func finishJobPath(topic, id string) string {
	return fmt.Sprintf(endPathFormat, topic, id)
}

func popJobPath(topic string) string {
	return fmt.Sprintf(popPathFormat, topic)
}
