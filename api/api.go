package api

type Request interface {
	AddDelayJob(topic, id, body string, delay, ttr uint) error
	ReplaceDelayJob(topic, id, body string, delay, ttr uint) error
	FinishDelayJob(topic string, id string) error
	DeleteDelayJob(topic string, id string) error
	PopDelayJob(topic string) (topicName, id, body string, delay, ttr uint64, err error)
}
