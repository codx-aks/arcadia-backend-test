package controller

import "time"

type AuctionItem struct {
	ID             uint         `json:"id"`
	CreatedAt      time.Time    `json:"created_at"`
	SellerID       uint         `json:"seller_id"`
	OwnedMiniconID uint         `json:"owned_minicon_id"`
	BasePrice      uint         `json:"base_price"`
	CurrentPrice   uint         `json:"current_price"`
	CurrentBuyerID uint         `json:"current_buyer_id"`
	Status         string       `json:"status"`
	OwnedMinicon   OwnedMinicon `json:"owned_minicon"`
}
type OwnedMinicon struct {
	ID        uint    `json:"id"`
	Health    uint    `json:"health"`
	Attack    uint    `json:"attack"`
	XP        uint    `json:"xp"`
	Level     uint    `json:"level"`
	MiniconID uint    `json:"minicon_id"`
	Minicon   Minicon `json:"minicon"`
	IsOwned   bool    `json:"is_owned"`
}
type Minicon struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageLink   string `json:"image_link"`
	TypeID      uint   `json:"type_id"`
}

type GetAuctionBuyerResponse struct {
	AuctionItems []AuctionItem `json:"auctionItems"`
}
