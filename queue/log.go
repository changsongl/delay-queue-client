package queue

import "fmt"

type Log interface {
	Write(msg string) error
}

type log struct {
}

func newLog() Log {
	return &log{}
}

func (l log) Write(msg string) error {
	fmt.Println(msg)
	return nil
}
