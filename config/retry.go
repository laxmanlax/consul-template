package config

import (
	"fmt"
	"math"
	"time"
)

const (
	// DefaultRetryAttempts is the default number of maximum retry attempts.
	DefaultRetryAttempts = 5

	// DefaultRetryBackoff is the default base for the exponential backoff
	// algorithm.
	DefaultRetryBackoff = 250 * time.Millisecond
)

// RetryFunc is the signature of a function that supports retries.
type RetryFunc func(int) (bool, time.Duration)

// RetryConfig is a shared configuration for upstreams that support retires on
// failure.
type RetryConfig struct {
	// Attempts is the total number of maximum attempts to retry before letting
	// the error fall through.
	Attempts *int

	// Backoff is the base of the exponentialbackoff. This number will be multipled
	// by the next power of 2 on each iteration.
	Backoff *time.Duration

	// Enabled signals if this retry is enabled.
	Enabled *bool
}

// DefaultRetryConfig returns a configuration that is populated with the
// default values.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *RetryConfig) Copy() *RetryConfig {
	if c == nil {
		return nil
	}

	var o RetryConfig

	o.Attempts = c.Attempts

	o.Backoff = c.Backoff

	o.Enabled = c.Enabled

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *RetryConfig) Merge(o *RetryConfig) *RetryConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Attempts != nil {
		r.Attempts = o.Attempts
	}

	if o.Backoff != nil {
		r.Backoff = o.Backoff
	}

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	return r
}

// RetryFunc returns the retry function associated with this configuration.
func (c *RetryConfig) RetryFunc() RetryFunc {
	return func(retry int) (bool, time.Duration) {
		if !BoolVal(c.Enabled) {
			return false, 0
		}

		if IntVal(c.Attempts) > 0 && retry > IntVal(c.Attempts)-1 {
			return false, 0
		}

		base := math.Pow(2, float64(retry))
		sleep := time.Duration(base) * TimeDurationVal(c.Backoff)

		return true, sleep
	}
}

// Finalize ensures there no nil pointers.
func (c *RetryConfig) Finalize() {
	if c.Attempts == nil {
		c.Attempts = Int(DefaultRetryAttempts)
	}

	if c.Backoff == nil {
		c.Backoff = TimeDuration(DefaultRetryBackoff)
	}

	if c.Enabled == nil {
		c.Enabled = Bool(true)
	}
}

// GoString defines the printable version of this struct.
func (c *RetryConfig) GoString() string {
	if c == nil {
		return "(*RetryConfig)(nil)"
	}

	return fmt.Sprintf("&RetryConfig{"+
		"Attempts:%s, "+
		"Backoff:%s, "+
		"Enabled:%s"+
		"}",
		IntGoString(c.Attempts),
		TimeDurationGoString(c.Backoff),
		BoolGoString(c.Enabled),
	)
}
