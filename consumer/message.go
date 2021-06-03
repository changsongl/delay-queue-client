package consumer

import "github.com/changsongl/delay-queue-client/client"

type Message interface {
	GetId() string
	GetBody() string
	GetTopic() string
	GetDelay() uint64
	GetTTR() uint64

	Finish() error
	Delete() error
}

type message struct {
	topic string
	body  string
	id    string
	delay uint64
	ttr   uint64

	client client.Client
}

func NewMessage(topic, body, id string, delay, ttr uint64, client client.Client) Message {
	return &message{
		topic:  topic,
		body:   body,
		id:     id,
		delay:  delay,
		ttr:    ttr,
		client: client,
	}
}

func (m *message) GetId() string {
	return m.id
}

func (m *message) GetBody() string {
	return m.body
}

func (m *message) GetTopic() string {
	return m.topic
}

func (m *message) GetDelay() uint64 {
	return m.delay
}

func (m *message) GetTTR() uint64 {
	return m.ttr
}

func (m *message) Finish() error {
	return m.client.FinishJob(m.topic, m.id)
}

func (m *message) Delete() error {
	return m.client.DeleteJob(m.topic, m.id)
}
