package model

import (
	"time"

	"github.com/google/uuid"
)

type UserActivationToken struct {
	UserID         uuid.UUID  `gorm:"primaryKey"`
	Token          string     `json:"-"`
	TokenExpiredAt *time.Time `json:"-"`
}

func NewUserActivationToken(userID uuid.UUID) *UserActivationToken {
	token, expiry := NewToken()

	return &UserActivationToken{
		UserID:         userID,
		Token:          token,
		TokenExpiredAt: expiry,
	}
}
