package job

import "time"

// Option job option
type Option interface {
	apply(j *Job)
}

type jobOptionFunc func(*Job)

func (f jobOptionFunc) apply(j *Job) {
	f(j)
}

// Job job object
type Job struct {
	Topic jobTopic `json:"jobTopic"`
	ID    jobID    `json:"jobID"`
	Delay jobDelay `json:"jobDelay"`
	TTR   jobTTR   `json:"jobTTR"`
	Body  jobBody  `json:"jobBody"`
}

// New a job
func New(topic, id string, opts ...Option) (*Job, error) {
	j := &Job{
		ID:    jobID(id),
		Topic: jobTopic(topic),
	}

	for _, opt := range opts {
		opt.apply(j)
	}

	return j, nil
}

// DelayOption job delay time
func DelayOption(delay time.Duration) Option {
	return jobOptionFunc(
		func(j *Job) {
			j.Delay = jobDelay(delay)
		},
	)
}

// TTROption job time to run
func TTROption(ttr time.Duration) Option {
	return jobOptionFunc(
		func(j *Job) {
			j.TTR = jobTTR(ttr)
		},
	)
}

// BodyOption job body
func BodyOption(body string) Option {
	return jobOptionFunc(
		func(j *Job) {
			j.Body = jobBody(body)
		},
	)
}

// check job is valid
func (j *Job) checkValid() error {
	fields := j.getCheckFields()
	for _, field := range fields {
		err := field.IsValid()
		if err != nil {
			return err
		}
	}

	return nil
}

// check field
func (j *Job) getCheckFields() []Field {
	return []Field{j.Topic, j.ID}
}

// ExtractData extract job data
func (j *Job) ExtractData() (topic, id, body string, delay, ttr uint) {
	return string(j.Topic), string(j.ID), string(j.Body),
		uint(time.Duration(j.Delay) / time.Second), uint(time.Duration(j.TTR) / time.Second)
}
