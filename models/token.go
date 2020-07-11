package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// TokenFactory :
type TokenFactory func() *Token

// Token :
type Token struct {
	ID        string
	CreatedAt time.Time
}

// NewToken :
func NewToken() *Token {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return &Token{
		ID:        id.String(),
		CreatedAt: time.Now(),
	}
}

// NeedReset :
func (token *Token) NeedReset(resetAfter time.Duration) bool {
	if time.Since(token.CreatedAt) >= resetAfter {
		return true
	}
	return false
}
