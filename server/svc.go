// Package utils /*
package utils

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

// Create variable signal.Notify function so we can mock it in tests
var signalNotify = signal.Notify

var ErrStop = errors.New("stopping service")

// Service interface contains Start and Stop methods which are called
// when the service is started and stopped. The Init method is called
// before the service is started, and after it's determined if the program
// is running as a Windows Service.
//
// The Start and Init methods must be non-blocking.
//
// Implement this interface and pass it to the Run function to start your program.
type Service interface {
	// Init is called before the program/service is started and after it's
	// determined if the program is running as a Windows Service. This method must
	// be non-blocking.
	Init(Environment) error

	// Start is called after Init. This method must be non-blocking.
	Start() error

	// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
	// Windows Service is stopped.
	Stop() error
}

// Context interface contains an optional Context function which a Service can implement.
// When implemented the context.Done() channel will be used in addition to signal handling
// to exit a process.
type Context interface {
	Context() context.Context
}

// Environment contains information about the environment
// your application is running in.
type Environment interface {
	// IsWindowsService reports whether the program is running as a Windows Service.
	IsWindowsService() bool
}

// Handler is an optional interface a Service can implement.
// When implemented, Handle() is called when a signal is received.
// Returning ErrStop from this method will result in Service.Stop() being called.
type Handler interface {
	Handle(os.Signal) error
}
