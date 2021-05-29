package consumer

import "time"

const (
	defaultBackOff       = 10 * time.Millisecond
	defaultBackOffInc    = 0 * time.Second
	defaultBackOffFactor = 2.0
	defaultBackOffMax    = 2 * time.Second

	defaultErrorBackOff       = 100 * time.Millisecond
	defaultErrorBackOffInc    = 0 * time.Second
	defaultErrorBackOffFactor = 2.0
	defaultErrorBackOffMax    = 10 * time.Second

	defaultTTR = 0 * time.Second

	defaultWorkNum = 0
)

type config struct {
	backOff       time.Duration
	backOffInc    time.Duration
	backOffFactor float64
	backOffMax    time.Duration

	errorBackOff       time.Duration
	errorBackOffInc    time.Duration
	errorBackOffFactor float64
	errorBackOffMax    time.Duration

	ttr time.Duration

	workerNum int64

	l Log
}

type Option func(config *config)

func newConfig() *config {
	return &config{
		backOff:            defaultBackOff,
		backOffInc:         defaultBackOffInc,
		backOffFactor:      defaultBackOffFactor,
		backOffMax:         defaultBackOffMax,
		errorBackOff:       defaultErrorBackOff,
		errorBackOffInc:    defaultErrorBackOffInc,
		errorBackOffFactor: defaultErrorBackOffFactor,
		errorBackOffMax:    defaultErrorBackOffMax,
		ttr:                defaultTTR,
		l:                  newLog(),
		workerNum:          defaultWorkNum,
	}
}

func (c *config) apply(options ...Option) {
	for _, opt := range options {
		opt(c)
	}
}

func BackOffOption(t time.Duration) Option {
	return func(config *config) {
		config.backOff = t
	}
}

func BackOffIncOption(t time.Duration) Option {
	return func(config *config) {
		config.backOffInc = t
	}
}

func BackOffMaxOption(t time.Duration) Option {
	return func(config *config) {
		config.backOffMax = t
	}
}

func BackOffFactor(f float64) Option {
	return func(config *config) {
		config.backOffFactor = f
	}
}

func ErrorBackOffOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOff = t
	}
}

func ErrorBackOffIncOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOffInc = t
	}
}

func ErrorBackOffMaxOption(t time.Duration) Option {
	return func(config *config) {
		config.errorBackOffMax = t
	}
}

func ErrorBackOffFactorOption(f float64) Option {
	return func(config *config) {
		config.errorBackOffFactor = f
	}
}

func TTROption(ttr time.Duration) Option {
	return func(config *config) {
		config.ttr = ttr
	}
}

func LoggerOption(l Log) Option {
	return func(config *config) {
		config.l = l
	}
}

func WorkerNumOption(num int64) Option {
	return func(config *config) {
		config.workerNum = num
	}
}
