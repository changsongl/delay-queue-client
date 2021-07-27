package consumer

import "fmt"

// Log interface for consumer
type Log interface {
	Write(msg string) error
}

// log default object
type log struct {
}

// create a Log
func newLog() Log {
	return &log{}
}

// implement Write
func (l log) Write(msg string) error {
	fmt.Println(msg)
	return nil
}
