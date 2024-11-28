package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;generated;default:gen_random_uuid()"`
	Name       string    `gorm:"type:varchar(50);not null"`
	Url        string    `gorm:"type:varchar(255);not null"`
	UploadedAt time.Time `gorm:"not null"`
}
