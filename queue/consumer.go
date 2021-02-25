package queue

import (
	"github.com/changsongl/delay-queue-client/client"
	"sync"
	"time"
)

type Queue interface {
	Consume() <-chan Message
	Close() bool
	IsClosed() bool
}

type queue struct {
	client client.Client
	topic  string
	cs     []chan Message
	sync.Mutex
	closed  bool
	bufSize int
	log     Log
}

func New(client client.Client, topic string) Queue {
	return &queue{
		client: client,
		topic:  topic,
		log:    newLog(),
	}
}

func (q *queue) Consume() <-chan Message {
	c := q.createChan()

	go func(q *queue, c chan Message) {
		for {
			if q.IsClosed() {
				close(c)
				break
			}

			id, body, err := q.client.PopJob(q.topic)
			if err != nil {
				_ = q.log.Write(err.Error())
				time.Sleep(2 * time.Second)
				continue
			}

			msg := NewMessage(q.topic, body, id, q.client)
			c <- msg
		}
	}(q, c)

	return c
}

func (q *queue) createChan() chan Message {
	var c chan Message
	if q.bufSize == 0 {
		c = make(chan Message)
	} else {
		c = make(chan Message, q.bufSize)
	}

	return c
}

func (q *queue) Close() bool {
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

func (q *queue) IsClosed() bool {
	q.Lock()
	defer q.Unlock()

	return q.closed
}
