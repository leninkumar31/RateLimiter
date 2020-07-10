package limiter

import "../models"

// IRateLimiter :
type IRateLimiter interface {
	Acquire() (*models.Token, error)
}
