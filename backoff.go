// backoff provides a backoff that helps to retry some work by slowing clients down.
package backoff

import (
	"context"
	"math/rand"
	"time"
)

const (
	// MaxRetries is a default number of retries.
	MaxRetries = 10
	// Multiplier is a default delay used as a multiplier that increases the wait interval.
	Multiplier = 25 * time.Millisecond
	// MaxWait is a maximum delay by default.
	MaxWait = 20 * time.Second
)

// Option sets up a Config.
type Option func(*Config)

// Config configures a retrier.
type Config struct {
	random *rand.Rand
	// maxRetries indicates how many times the work should be retried.
	maxRetries float64
	// multiplierSeconds is a multiplier that increases the wait interval, e.g., 5 seconds.
	multiplierSeconds float64
	// maxWaitSeconds is a max waiting time between attempts in seconds, e.g., 300 seconds.
	maxWaitSeconds float64
}

// WithRand sets a pseudo-random number generator.
// It's primarily used to facilitate testing.
func WithRand(r *rand.Rand) Option {
	return func(c *Config) {
		c.random = r
	}
}

// WithMaxRetries sets an upper limit on retries.
// It should be >= 0 or else the default value MaxRetries is used.
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		if n < 0 {
			c.maxRetries = MaxRetries
			return
		}
		c.maxRetries = float64(n)
	}
}

// WithMultiplier sets a multiplier that increases the wait interval.
// It should be > 0ms or else the default value Multiplier is used.
func WithMultiplier(d time.Duration) Option {
	return func(c *Config) {
		if d <= 0 {
			c.multiplierSeconds = Multiplier.Seconds()
			return
		}
		c.multiplierSeconds = d.Seconds()
	}
}

// WithMaxWait sets an upper limit of waiting time between attempts.
// It should be > 0ms or else the default value MaxWait is used.
func WithMaxWait(d time.Duration) Option {
	return func(c *Config) {
		if d <= 0 {
			c.maxWaitSeconds = MaxWait.Seconds()
			return
		}
		c.maxWaitSeconds = d.Seconds()
	}
}

// Retryer represents a retryer whose Next func returns true
// if an attempt should be made.
// Delay func returns a wait duration between attempts unless it's zero.
type Retryer interface {
	Next() bool
	Delay() time.Duration
	Reset()
}

// Run calls f and if it returns an error,
// new f calls will be made with backoff until
// retryer r finishes or a context is cancelled.
//
// An error from the last attempt is returned.
func Run(ctx context.Context, r Retryer, f func(attempt int) error) error {
	var (
		err error
		d   time.Duration
	)

Loop:
	for i := 1; r.Next(); i++ {
		if err = f(i); err == nil {
			break Loop
		}

		d = r.Delay()
		if d > 0 {
			select {
			case <-time.After(d):
			case <-ctx.Done():
				break Loop
			}
		}
	}

	return err
}
