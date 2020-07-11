package client

import (
	"errors"
	"time"

	"../limiter"
	"../models"
)

// NewThrottleRateLimiter :
func NewThrottleRateLimiter(config *models.Config) (limiter.IRateLimiter, error) {
	if config.Throttle <= 0 {
		return nil, errors.New("Throttle must be greater than 0")
	}
	ratelimiter := limiter.NewRateLimiter(config)
	await := func(throttle time.Duration) {
		ticker := time.NewTicker(throttle)
		go func() {
			for ; true; <-ticker.C {
				<-ratelimiter.Incoming
				ratelimiter.CreateToken()
			}
		}()
	}
	await(config.Throttle)
	return ratelimiter, nil
}

// NewMaxConcurrencyLimiter :
func NewMaxConcurrencyLimiter(config *models.Config) (limiter.IRateLimiter, error) {
	if config.Limit <= 0 {
		return nil, errors.New("Limit must be greater than 0")
	}
	ratelimiter := limiter.NewRateLimiter(config)
	await := func() {
		go func() {
			for {
				select {
				case <-ratelimiter.Incoming:
					ratelimiter.CreateToken()
				case t := <-ratelimiter.ReleaseChan:
					ratelimiter.ReleaseToken(t)
				}
			}
		}()
	}
	await()
	return ratelimiter, nil
}

// NewFixedWindowRateLimiter :
func NewFixedWindowRateLimiter(config *models.Config) (limiter.IRateLimiter, error) {
	if config.FixedInterval == 0 {
		return nil, errors.New("Fixed Interval should be greater than zero")
	}
	if config.Limit == 0 {
		return nil, errors.New("Limit should be greater than zero")
	}
	ratelimiter := limiter.NewRateLimiter(config)
	fixedWindow := &models.FixedWindowInterval{Interval: config.FixedInterval}
	makeToken := func() *models.Token {
		token := models.NewToken()
		token.ExpiredAt = fixedWindow.EndTime
		return token
	}
	ratelimiter.GenerateToken = makeToken
	await := func() {
		go func() {
			for {
				select {
				case <-ratelimiter.Incoming:
					ratelimiter.CreateToken()
				case t := <-ratelimiter.ReleaseChan:
					ratelimiter.ReleaseToken(t)
				}
			}
		}()
	}
	fixedWindow.Run(ratelimiter.RunReleaseExpiredTokens)
	await()
	return ratelimiter, nil
}
