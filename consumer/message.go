package consumer

import "github.com/changsongl/delay-queue-client/client"

// Message from delay queue
type Message interface {
	GetID() string
	GetBody() string
	GetTopic() string
	GetDelay() uint64
	GetTTR() uint64

	Finish() error
	Delete() error
}

// message object
type message struct {
	topic string
	body  string
	id    string
	delay uint64
	ttr   uint64

	client client.Client
}

// NewMessage create a new message
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

// GetID get id of the job
func (m *message) GetID() string {
	return m.id
}

// GetBody get body of the job
func (m *message) GetBody() string {
	return m.body
}

// GetTopic get topic of the job
func (m *message) GetTopic() string {
	return m.topic
}

// GetDelay get delay time of the job
func (m *message) GetDelay() uint64 {
	return m.delay
}

// GetTTR get time to run of the job
func (m *message) GetTTR() uint64 {
	return m.ttr
}

// Finish the job
func (m *message) Finish() error {
	return m.client.FinishJob(m.topic, m.id)
}

// Delete the job
func (m *message) Delete() error {
	return m.client.DeleteJob(m.topic, m.id)
}
