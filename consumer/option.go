package consumer

import "time"

const (
	// popping timeout
	defaultPopTimeout = 4 * time.Second

	// error backoff times
	defaultErrorBackOff       = 100 * time.Millisecond
	defaultErrorBackOffInc    = 0 * time.Second
	defaultErrorBackOffFactor = 2.0
	defaultErrorBackOffMax    = 10 * time.Second

	// default ttr time
	defaultTTR = 0 * time.Second

	// default work num
	defaultWorkNum = 0
)

// config of consumer
type config struct {
	popTimeout         time.Duration
	errorBackOff       time.Duration
	errorBackOffInc    time.Duration
	errorBackOffFactor float64
	errorBackOffMax    time.Duration

	workerNum int64

	l Log
}

// Option func
type Option func(config *config)

// create a config
func newConfig() *config {
	return &config{
		popTimeout:         defaultPopTimeout,
		errorBackOff:       defaultErrorBackOff,
		errorBackOffInc:    defaultErrorBackOffInc,
		errorBackOffFactor: defaultErrorBackOffFactor,
		errorBackOffMax:    defaultErrorBackOffMax,
		l:                  newLog(),
		workerNum:          defaultWorkNum,
	}
}

func (c *config) apply(options ...Option) {
	for _, opt := range options {
		opt(c)
	}
}

// PopTimeoutOption pop job timeout
func PopTimeoutOption(t time.Duration) Option {
	return func(config *config) {
		config.popTimeout = t
	}
}

// ErrorBackOffOption delay queue error backoff time
func ErrorBackOffOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOff = t
	}
}

// ErrorBackOffIncOption delay queue error backoff time increment
func ErrorBackOffIncOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOffInc = t
	}
}

// ErrorBackOffMaxOption max error backoff time
func ErrorBackOffMaxOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOffMax = t
	}
}

// ErrorBackOffFactorOption error backoff time increasing factor
func ErrorBackOffFactorOption(f float64) Option {
	return func(config *config) {
		config.errorBackOffFactor = f
	}
}

// LoggerOption set consumer logger
func LoggerOption(l Log) Option {
	return func(config *config) {
		config.l = l
	}
}

// WorkerNumOption worker num to consume
func WorkerNumOption(num int64) Option {
	return func(config *config) {
		config.workerNum = num
	}
}
