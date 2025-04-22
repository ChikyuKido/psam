package models

import (
	"time"
)

type APIKey struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}
