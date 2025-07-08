package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

type Model2 struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}
