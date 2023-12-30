package model

import (
	"gorm.io/gorm"
)

type MiniconType struct {
	gorm.Model
	Name        string `gorm:"not null; unique;"`
	Description string `gorm:"not null;"`
}
