package limiter

import (
	"../models"
)

// RateLimiter :
type RateLimiter struct {
	Incoming      chan int
	Outgoing      chan *models.Token
	Err           chan error
	GenerateToken models.TokenFactory
}

// NewRateLimiter :
func NewRateLimiter(config *models.Config) *RateLimiter {
	return &RateLimiter{
		Incoming:      make(chan int),
		Outgoing:      make(chan *models.Token),
		Err:           make(chan error),
		GenerateToken: models.NewToken,
	}
}

// Acquire :
func (ratelimiter *RateLimiter) Acquire() (*models.Token, error) {
	go func() {
		ratelimiter.Incoming <- 0
	}()
	select {
	case t := <-ratelimiter.Outgoing:
		return t, nil
	case err := <-ratelimiter.Err:
		return nil, err
	}
}

// CreateToken :
func (ratelimiter *RateLimiter) CreateToken() {
	if ratelimiter.GenerateToken == nil {
		panic("Error Token Factory Not Defined")
	}
	t := ratelimiter.GenerateToken()
	go func() {
		ratelimiter.Outgoing <- t
	}()
}
