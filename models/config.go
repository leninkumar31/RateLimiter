package models

import "time"

// Config :
type Config struct {
	// This is min time between requests for Throttle Rate Limiter
	Throttle time.Duration
	// Number of Active connections to be allowed in given time
	Limit int
	// After how much time we should forcefully release the token if it is not release
	TokenResetAfter time.Duration
	// Fixed interval is defined for fixed window rate limiter
	FixedInterval time.Duration
}
