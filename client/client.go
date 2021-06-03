package client

import (
	"github.com/changsongl/delay-queue-client/api"
	"github.com/changsongl/delay-queue-client/api/http"
	"github.com/changsongl/delay-queue-client/job"
)

type Client interface {
	AddJob(job *job.Job) error
	ReplaceJob(job *job.Job) error
	FinishJob(topic string, id string) error
	DeleteJob(topic string, id string) error
	PopJob(topic string) (topicName, id, body string, delay, ttr uint64, err error)
}

type client struct {
	request api.Request
}

func NewClient(host string) Client {
	res := http.NewRequester(host)

	c := &client{
		request: res,
	}

	return c
}

func (c client) AddJob(job *job.Job) error {
	return c.request.AddDelayJob(job.ExtractData())
}

func (c client) ReplaceJob(job *job.Job) error {
	return c.request.ReplaceDelayJob(job.ExtractData())
}

func (c client) FinishJob(topic string, id string) error {
	return c.request.FinishDelayJob(topic, id)
}

func (c client) DeleteJob(topic string, id string) error {
	return c.request.DeleteDelayJob(topic, id)
}

func (c client) PopJob(topic string) (topicName, id, body string, delay, ttr uint64, err error) {
	return c.request.PopDelayJob(topic)
}
