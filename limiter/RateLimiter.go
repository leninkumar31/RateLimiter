package limiter

import (
	"fmt"
	"sync/atomic"
	"time"

	"../models"
)

// RateLimiter :
type RateLimiter struct {
	Incoming      chan struct{}
	Outgoing      chan *models.Token
	ReleaseChan   chan *models.Token
	Err           chan error
	limit         int
	needToken     int64
	activeTokens  map[string]*models.Token
	GenerateToken models.TokenFactory
}

// NewRateLimiter :
func NewRateLimiter(config *models.Config) *RateLimiter {
	r := &RateLimiter{
		Incoming:      make(chan struct{}),
		Outgoing:      make(chan *models.Token),
		ReleaseChan:   make(chan *models.Token),
		Err:           make(chan error),
		limit:         config.Limit,
		needToken:     0,
		activeTokens:  make(map[string]*models.Token),
		GenerateToken: models.NewToken,
	}
	if config.TokenResetAfter > 0 {
		go r.runResetTokenTask(config.TokenResetAfter)
	}
	return r
}

// Acquire :
func (ratelimiter *RateLimiter) Acquire() (*models.Token, error) {
	go func() {
		ratelimiter.Incoming <- struct{}{}
	}()
	select {
	case t := <-ratelimiter.Outgoing:
		return t, nil
	case err := <-ratelimiter.Err:
		return nil, err
	}
}

// Release :
func (ratelimiter *RateLimiter) Release(token *models.Token) {
	if token.IsExpired() {
		go func() {
			ratelimiter.ReleaseChan <- token
		}()
	}
}

// CreateToken :
func (ratelimiter *RateLimiter) CreateToken() {
	if ratelimiter.GenerateToken == nil {
		panic("Error Token Factory Not Defined")
	}
	if ratelimiter.isLimitExceeded() {
		ratelimiter.incNeedToken()
		return
	}
	t := ratelimiter.GenerateToken()
	ratelimiter.activeTokens[t.ID] = t
	go func() {
		ratelimiter.Outgoing <- t
	}()
}

func (ratelimiter *RateLimiter) incNeedToken() {
	atomic.AddInt64(&ratelimiter.needToken, 1)
}

func (ratelimiter *RateLimiter) decNeedToken() {
	atomic.AddInt64(&ratelimiter.needToken, -1)
}

func (ratelimiter *RateLimiter) awaitingToken() bool {
	return atomic.LoadInt64(&ratelimiter.needToken) > 0
}

func (ratelimiter *RateLimiter) isLimitExceeded() bool {
	if len(ratelimiter.activeTokens) >= ratelimiter.limit {
		return true
	}
	return false
}

// ReleaseToken :
func (ratelimiter *RateLimiter) ReleaseToken(token *models.Token) {
	if token == nil {
		fmt.Println("Unable to release nil token")
		return
	}
	if _, ok := ratelimiter.activeTokens[token.ID]; !ok {
		fmt.Printf("Unable to release token %s not in use\n", token)
		return
	}
	delete(ratelimiter.activeTokens, token.ID)
	if ratelimiter.awaitingToken() {
		ratelimiter.decNeedToken()
		go ratelimiter.CreateToken()
	}
}

func (ratelimiter *RateLimiter) runResetTokenTask(resetAfter time.Duration) {
	go func() {
		ticker := time.NewTicker(resetAfter)
		for range ticker.C {
			for _, token := range ratelimiter.activeTokens {
				if token.NeedReset(resetAfter) {
					ratelimiter.ReleaseChan <- token
				}
			}
		}
	}()
}

// RunReleaseExpiredTokens :
func (ratelimiter *RateLimiter) RunReleaseExpiredTokens() {
	for _, token := range ratelimiter.activeTokens {
		if token.IsExpired() {
			go func(token *models.Token) {
				ratelimiter.ReleaseChan <- token
			}(token)
		}
	}
}
