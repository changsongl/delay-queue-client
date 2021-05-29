package consumer

import (
	"github.com/changsongl/delay-queue-client/client"
	"sync"
	"time"
)

type consumeStatus int

const (
	consumerLoadSuccess consumeStatus = iota
	consumerLoadFailed
	consumerError
)

type Consumer interface {
	Consume() <-chan Message
	Close() bool
	IsClosed() bool
}

type consumer struct {
	client client.Client
	topic  string
	cs     []chan Message
	sync.Mutex
	closed bool
	conf   *config
}

type timer struct {
	backOff      time.Duration
	errorBackOff time.Duration
}

func New(client client.Client, topic string, options ...Option) Consumer {
	conf := newConfig()
	conf.apply(options...)

	return &consumer{
		client: client,
		topic:  topic,
		conf:   conf,
	}
}

func (q *consumer) Consume() <-chan Message {
	c := q.createChan()

	producerNum := q.conf.workerNum
	if producerNum == 0 {
		producerNum = 1
	}

	for i := 0; i < int(producerNum); i++ {
		go func(q *consumer, c chan Message) {
			timer := &timer{
				backOff:      q.conf.backOff,
				errorBackOff: q.conf.errorBackOff,
			}

			for {
				if q.IsClosed() {
					close(c)
					break
				}

				id, body, err := q.client.PopJob(q.topic)
				if err != nil {
					_ = q.log(err.Error())
					q.Wait(timer, consumerError)
					continue
				}

				if id == "" {
					q.Wait(timer, consumerLoadFailed)
				} else {
					msg := NewMessage(q.topic, body, id, q.client)
					c <- msg
					q.Wait(timer, consumerLoadSuccess)
				}
			}
		}(q, c)
	}

	return c
}

func (q *consumer) createChan() chan Message {
	var c chan Message
	if q.conf.workerNum == 0 {
		c = make(chan Message)
	} else {
		c = make(chan Message, q.conf.workerNum)
	}

	return c
}

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

func (q *consumer) IsClosed() bool {
	q.Lock()
	defer q.Unlock()

	return q.closed
}

func (q *consumer) log(msg string) error {
	return q.conf.l.Write(msg)
}

func (q *consumer) Wait(t *timer, status consumeStatus) {
	switch status {
	case consumerLoadSuccess:
		t.backOff = q.conf.backOff
		t.errorBackOff = q.conf.errorBackOff
	case consumerLoadFailed:
		time.Sleep(t.backOff)
		t.backOff = time.Duration(float64(t.backOff+q.conf.backOffInc) * q.conf.backOffFactor)
		if t.backOff > q.conf.backOffMax {
			t.backOff = q.conf.backOffMax
		}
		t.errorBackOff = q.conf.errorBackOff
	case consumerError:
		time.Sleep(t.errorBackOff)
		t.backOff = q.conf.backOff
		t.errorBackOff = time.Duration(float64(t.errorBackOff+q.conf.errorBackOffInc) * q.conf.errorBackOffFactor)
		if t.errorBackOff > q.conf.errorBackOffMax {
			t.errorBackOff = q.conf.errorBackOffMax
		}
	}
}
