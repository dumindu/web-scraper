package model

import (
	"time"

	"github.com/google/uuid"
)

type UserAuth struct {
	UserID    uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Password  string     `json:"-"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

func NewUserAuth(userID uuid.UUID, hashedPassword string) *UserAuth {
	return &UserAuth{
		UserID:   userID,
		Password: hashedPassword,
	}
}
