package client

import (
	"github.com/changsongl/delay-queue-client/api"
	"github.com/changsongl/delay-queue-client/api/http"
)

type Client interface {
	AddJob(job *Job) error
	ReplaceJob(job *Job) error
	FinishJob(topic string, id string) error
	DeleteJob(topic string, id string) error
	PopJob(topic string) (id string, body string, err error)
}

type client struct {
	request api.Request
}

func New() Client {
	res := http.NewRequester()
	c := &client{
		request: res,
	}

	return c
}

func (c client) AddJob(job *Job) error {
	return c.request.AddDelayJob(job.extractData())
}

func (c client) ReplaceJob(job *Job) error {
	return c.request.ReplaceDelayJob(job.extractData())
}

func (c client) FinishJob(topic string, id string) error {
	return c.request.FinishDelayJob(topic, id)
}

func (c client) DeleteJob(topic string, id string) error {
	return c.request.DeleteDelayJob(topic, id)
}

func (c client) PopJob(topic string) (id string, body string, err error) {
	return c.request.PopDelayJob(topic)
}
