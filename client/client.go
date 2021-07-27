package client

import (
	"github.com/changsongl/delay-queue-client/api"
	"github.com/changsongl/delay-queue-client/api/http"
	"github.com/changsongl/delay-queue-client/job"
	"time"
)

// Client client interface
type Client interface {
	AddJob(job *job.Job) error
	ReplaceJob(job *job.Job) error
	FinishJob(topic string, id string) error
	DeleteJob(topic string, id string) error
	PopJob(topic string, timeout time.Duration) (topicName, id, body string, delay, ttr uint64, err error)
}

// client object
type client struct {
	request api.Request
}

// NewClient create a new delay queue client
func NewClient(host string) Client {
	res := http.NewRequester(host)

	c := &client{
		request: res,
	}

	return c
}

// AddJob add job to delay queue, if job is already exist, it will return error
func (c client) AddJob(job *job.Job) error {
	return c.request.AddDelayJob(job.ExtractData())
}

// ReplaceJob replace job in delay queue, it will replace or add job to delay queue
func (c client) ReplaceJob(job *job.Job) error {
	return c.request.ReplaceDelayJob(job.ExtractData())
}

// FinishJob finish job, when have processed job
func (c client) FinishJob(topic string, id string) error {
	return c.request.FinishDelayJob(topic, id)
}

// DeleteJob delete job from delay queue
func (c client) DeleteJob(topic string, id string) error {
	return c.request.DeleteDelayJob(topic, id)
}

// PopJob pop job from delay queue, and wait at least timeout seconds
func (c client) PopJob(topic string, timeout time.Duration) (topicName, id, body string, delay, ttr uint64, err error) {
	return c.request.PopDelayJob(topic, timeout)
}
