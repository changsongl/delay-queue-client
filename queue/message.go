package queue

import "github.com/changsongl/delay-queue-client/client"

type Message interface {
	GetId() string
	GetBody() string
	GetTopic() string

	Finish() error
	Delete() error
}

type message struct {
	topic string
	body  string
	id    string

	client client.Client
}

func NewMessage(topic, body, id string, client client.Client) Message {
	return &message{
		topic:  topic,
		body:   body,
		id:     id,
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

func (m *message) Finish() error {
	return m.client.FinishJob(m.topic, m.id)
}

func (m *message) Delete() error {
	return m.client.DeleteJob(m.topic, m.id)
}
