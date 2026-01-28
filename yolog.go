// Package yolog is a logging package where You Only LOG.... once. Inspired by https://loggingsucks.com.
package yolog

import (
	"io"
	"log/slog"
	"time"
)

const (
	// AttrStream used to filter logs in log aggregation systems.
	AttrStream = "stream"
	// AttrElapsedTime tracking. Used for performance measurement.
	AttrElapsedTime = "elapsed_time_ms"
	// AttrErrorCode for programmatically analysing log events.
	AttrErrorCode = "error_code"
	// AttrErrorMessage for deeper understanding of the error.
	AttrErrorMessage = "error_message"
	// AttrErrors for capturing multiple errors in a single log entry.
	AttrErrors = "errors"
)

// Logger represents a logging logger for yolog.
type Logger struct {
	stream     string
	sampleRate float64
	startTime  time.Time
	attributes map[string]interface{}
	errors     []string
}

// OptFunc defines a function type for manipulating the Logger thoughout the lifecycle of an event.
// This function is used to set options when creating a new Logger as well as when logging.
type OptFunc func(*Logger)

// NewLogger creates a new yolog Logger with default settings.
func NewLogger(stream string, optFuncs ...OptFunc) *Logger {
	logger := &Logger{
		stream:     stream,
		startTime:  time.Now(),
		sampleRate: sampleRateFromEnv(),
		attributes: make(map[string]interface{}),
	}

	for _, fn := range optFuncs {
		fn(logger)
	}

	return logger
}

// SetAttr sets a custom attribute for the log entry eg. client_id, request_id etc
func (c *Logger) SetAttr(key string, value interface{}) *Logger {
	c.attributes[key] = value
	return c
}

// SetAttrs sets custom attributes for the log entry (e.g. client_id, request_id).
// Arguments are expected as key, value pairs. Invalid pairs are ignored.
func (c *Logger) SetAttrs(kv ...any) *Logger {
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue
		}
		c.attributes[key] = kv[i+1]
	}

	return c
}

// SetErrorCode sets the error code attribute for the log entry. If set the log will always be sampled.
func (c *Logger) SetErrorCode(code interface{}) *Logger {
	c.attributes[AttrErrorCode] = code
	return c
}

// WrapError appends an error message to the log entry and returns the original error.
// Handy when handling "if err != nil" statements.
func (c *Logger) WrapError(err error) error {
	if err == nil {
		return nil
	}

	c.AddError(err)
	return err
}

// AddError appends an error message to the log entry. If any errors are added the log will always be sampled.
func (c *Logger) AddError(err error) *Logger {
	c.errors = append(c.errors, err.Error())
	return c
}

// Returns true if the log entry has errors associated with it.
func (c *Logger) hasErrors() bool {
	if _, ok := c.attributes[AttrErrorCode]; ok {
		return true
	}

	if len(c.errors) > 0 {
		return true
	}

	return false
}

// Log writes the log entry to the provided writer if it is sampled.
// The log entry is encoded as JSON and typically deferred by the caller.
func (c *Logger) Log(w io.Writer) {
	if !c.shouldSample() {
		return
	}

	logger := slog.New(NewNoMsgJSONHandler(w, nil)).With(AttrStream, c.stream)

	// Calculate elapsed time so we can determine performance characteristics.
	logger = logger.With(AttrElapsedTime, time.Since(c.startTime).Milliseconds())

	for k, v := range c.attributes {
		logger = logger.With(k, v)
	}

	if len(c.errors) > 0 {
		logger = logger.With(AttrErrors, c.errors)
	}

	if c.hasErrors() {
		logger.Error("") // Empty message as we don't use it.
	} else {
		logger.Info("") // Empty message as we don't use it.
	}
}
