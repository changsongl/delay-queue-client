package client

import "time"

type JobOption interface {
	apply(j *Job)
}

type jobOptionFunc func(*Job)

func (f jobOptionFunc) apply(j *Job) {
	f(j)
}

type Job struct {
	Topic jobTopic `json:"jobTopic"`
	ID    jobId    `json:"jobId"`
	Delay jobDelay `json:"jobDelay"`
	TTR   jobTTR   `json:"jobTTR"`
	Body  jobBody  `json:"jobBody"`
}

func NewJob(topic, id string, opts ...JobOption) (*Job, error) {
	j := &Job{
		ID:    jobId(id),
		Topic: jobTopic(topic),
	}

	for _, opt := range opts {
		opt.apply(j)
	}

	return j, nil
}

func JobDelayOption(delay time.Duration) JobOption {
	return jobOptionFunc(
		func(j *Job) {
			j.Delay = jobDelay(delay)
		},
	)
}

func JobTTROption(ttr time.Duration) JobOption {
	return jobOptionFunc(
		func(j *Job) {
			j.TTR = jobTTR(ttr)
		},
	)
}

func JobBodyOption(body string) JobOption {
	return jobOptionFunc(
		func(j *Job) {
			j.Body = jobBody(body)
		},
	)
}

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

func (j *Job) getCheckFields() []JobField {
	return []JobField{j.Topic, j.ID}
}

func (j *Job) extractData() (topic, id, body string, delay, ttr uint) {
	return string(j.Topic), string(j.ID), string(j.Body),
		uint(time.Duration(j.Delay) / time.Second), uint(time.Duration(j.TTR) / time.Second)
}
