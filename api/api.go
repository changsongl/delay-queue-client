package api

type Request interface {
	AddDelayJob(topic, id, body string, delay, ttr uint) error
	ReplaceDelayJob(topic, id, body string, delay, ttr uint) error
	FinishDelayJob(topic string, id string) error
	DeleteDelayJob(topic string, id string) error
	PopDelayJob(topic string) (id string, body string, err error)
}
