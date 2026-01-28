package yolog

import (
	"log"
	"math/rand/v2"
	"os"
	"strconv"
)

const (
	// Default sample rate if not specified in environment.
	// We want to log everything by default and then sample down when needed.
	defaultSampleRate = 1.0
	// Environment variable for sample rate configuration.
	envSampleRate = "YOLOG_SAMPLE_RATE"
)

// WithSampleRate configures the sample rate for the Logger.
// The rate must be between 0.0 and 1.0.
func WithSampleRate(rate float64) OptFunc {
	return func(c *Logger) {
		c.sampleRate = rate
	}
}

// Determines if the log entry should be sampled based on the sample rate and error code presence.
func (c *Logger) shouldSample() bool {
	// Always log errors
	if c.hasErrors() {
		return true
	}

	// Always log error messages
	if _, ok := c.attributes[AttrErrorMessage]; ok {
		return true
	}

	return rand.Float64() < c.sampleRate
}

// Returns the sample rate from the environment variable, with validation and fallback to default.
func sampleRateFromEnv() float64 {
	v := os.Getenv(envSampleRate)
	if v == "" {
		return defaultSampleRate
	}

	rate, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Printf("yolog: invalid %s=%q (parse error: %v); using default %.2f", envSampleRate, v, err, defaultSampleRate)
		return defaultSampleRate
	}

	if rate < 0.0 || rate > 1.0 {
		log.Printf("yolog: invalid %s=%q (out of range); using default %.2f", envSampleRate, v, defaultSampleRate)
		return defaultSampleRate
	}

	return rate
}
