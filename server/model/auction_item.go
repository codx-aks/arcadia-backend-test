package model

import (
	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	SellerID       uint         `gorm:"not null;"`
	OwnedMiniconID uint         `gorm:"not null;"`
	OwnedMinicon   OwnedMinicon `gorm:"foreignKey:OwnedMiniconID;"`

	BasePrice      uint `gorm:"not null;"`
	CurrentPrice   uint
	CurrentBuyerID uint
	Status         string `gorm:"not null;"`
}
