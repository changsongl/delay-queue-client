package consumer

import (
	"github.com/changsongl/delay-queue-client/client"
	"github.com/changsongl/delay-queue-client/common"
	"sync"
	"time"
)

type consumeStatus int

const (
	// consumer has error
	consumerError consumeStatus = iota
)

// Consumer interface
type Consumer interface {
	Consume() <-chan Message
	Close() bool
	IsClosed() bool
}

// consumer object
type consumer struct {
	client client.Client
	topic  string
	cs     []chan Message
	sync.Mutex
	closed bool
	conf   *config
}

// timer for sleep
type timer struct {
	errorBackOff time.Duration
}

// New a consumer
func New(client client.Client, topic string, options ...Option) Consumer {
	conf := newConfig()
	conf.apply(options...)

	return &consumer{
		client: client,
		topic:  topic,
		conf:   conf,
	}
}

// Consume the delay queue for the topic
func (q *consumer) Consume() <-chan Message {
	c := q.createChan()

	producerNum := q.conf.workerNum
	if producerNum == 0 {
		producerNum = 1
	}

	for i := 0; i < int(producerNum); i++ {
		go func(q *consumer, c chan Message) {
			timer := &timer{
				errorBackOff: q.conf.errorBackOff,
			}

			for {
				if q.IsClosed() {
					close(c)
					break
				}

				topic, id, body, delay, ttr, err := q.client.PopJob(q.topic, q.conf.popTimeout)
				if err == common.ErrorNoAvailableJob {
					continue
				} else if err != nil {
					_ = q.log(err.Error())
					q.Wait(timer, consumerError)
					continue
				}

				msg := NewMessage(topic, body, id, delay, ttr, q.client)
				c <- msg
			}
		}(q, c)
	}

	return c
}

// create a chan for consumer
func (q *consumer) createChan() chan Message {
	var c chan Message
	if q.conf.workerNum == 0 {
		c = make(chan Message)
	} else {
		c = make(chan Message, q.conf.workerNum)
	}

	return c
}

// Close close consumer
func (q *consumer) Close() bool {
	if q.closed {
		return false
	}

	q.Lock()
	defer q.Unlock()
	if q.closed {
		return false
	}

	q.closed = true
	return true
}

// IsClosed check it is closed
func (q *consumer) IsClosed() bool {
	q.Lock()
	defer q.Unlock()

	return q.closed
}

// log msg
func (q *consumer) log(msg string) error {
	return q.conf.l.Write(msg)
}

// Wait timer sleep function
func (q *consumer) Wait(t *timer, status consumeStatus) {
	switch status {
	case consumerError:
		time.Sleep(t.errorBackOff)
		t.errorBackOff = time.Duration(float64(t.errorBackOff+q.conf.errorBackOffInc) * q.conf.errorBackOffFactor)
		if t.errorBackOff > q.conf.errorBackOffMax {
			t.errorBackOff = q.conf.errorBackOffMax
		}
	}
}
