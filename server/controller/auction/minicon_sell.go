package controller

import (
	"math"
	"net/http"
	"time"

	"github.com/delta/arcadia-backend/database"
	helper "github.com/delta/arcadia-backend/server/helper/general"
	"github.com/delta/arcadia-backend/server/model"
	"github.com/delta/arcadia-backend/utils"
	"github.com/gin-gonic/gin"
)

type CreateAuctionRequest struct {
	SellerID       uint `json:"seller_id" binding:"required"`
	OwnedMiniconID uint `json:"owned_minicon_id" binding:"required"`
}
type LFactorRequest struct {
	BasePrice      uint               `json:"base_price"`
	CurrentPrice   uint               `json:"current_price"`
	OwnedMiniconID uint               `json:"owned_minicon_id"`
	OwnedMinicon   model.OwnedMinicon `json:"owned_minicon"`
}

func calculateBasePrice(ownedMinicon model.OwnedMinicon) float64 {
	health_norm := float64(float64(ownedMinicon.Health-940) / float64(1200-940)) //min-max-normalisation
	attack_norm := float64(float64(ownedMinicon.Attack-400) / float64(600-400))
	base_price := float64((0.5*float64(health_norm) + 0.5*float64(attack_norm)) * 5000)
	return base_price
}

// highest price is 150% of the final base_price
// bid_threshold will be 2.5% of the final base_price
func SellPOST(c *gin.Context) {

	var request CreateAuctionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	var ownedMinicon model.OwnedMinicon
	db := database.GetDB()
	if err := db.Model(model.OwnedMinicon{}).
		Where("id = ?", request.OwnedMiniconID).
		First(&ownedMinicon).Error; err != nil {
		log := utils.GetControllerLogger("/api/auction/sell [POST]")
		log.Errorln(err)
		helper.SendError(c, http.StatusNotFound, "OwnedMinicon not found")
		return
	}
	// min_id := ownedMinicon.MiniconID
	var LAllItems []LFactorRequest
	var timeCap = time.Now().Add(-6 * time.Hour)
	if err := db.Model(model.Auction{}).
		Preload("OwnedMinicon").
		Preload("OwnedMinicon.Minicon").
		Where("created_at >= ?", timeCap.UTC()).
		Find(&LAllItems).Error; err != nil {
		log := utils.GetControllerLogger("/api/auction/sell [POST]")
		log.Errorln(err)

		helper.SendError(c, http.StatusInternalServerError, "Unable to fetch item for Lfactor")
		return
	}
	var minicon_demand float64 = 0.0
	var minicon_count float64 = 0.0
	var overall_demand float64 = 0.0
	var overall_count float64 = 0.0

	for index, value := range LAllItems {
		if value.OwnedMinicon.MiniconID == ownedMinicon.MiniconID {
			num := float64(2.0 * float64(float64(value.CurrentPrice-value.BasePrice)/float64(value.BasePrice)))
			if value.CurrentPrice >= value.BasePrice {
				minicon_demand += num
			}
			minicon_count++
		}
		nums := float64(2.0 * float64(float64(value.CurrentPrice-value.BasePrice)/float64(value.BasePrice)))
		if value.CurrentPrice >= value.BasePrice {
			overall_demand += nums
		}
		overall_count = float64(index)

	}
	overall_count++
	if overall_count != 0 {
		if minicon_count != 0 {
			minicon_demand /= minicon_count
		}
		overall_demand /= overall_count
	}

	Lfactor := 1.0 + minicon_demand - overall_demand
	basePrice := uint(calculateBasePrice(ownedMinicon) * (Lfactor))
	basePrice = uint(math.Round(float64(basePrice)/100.0)) * 100
	currentPrice := basePrice - uint(0.025*float64(basePrice))
	newAuction := model.Auction{
		SellerID:       request.SellerID,
		OwnedMiniconID: request.OwnedMiniconID,
		BasePrice:      basePrice,
		CurrentPrice:   currentPrice,
		Status:         "ongoing",
	}

	if err := db.Create(&newAuction).Error; err != nil {
		log := utils.GetControllerLogger("/api/auction/sell [POST]")
		log.Errorln(err)
		helper.SendError(c, http.StatusInternalServerError, "Unable to create auction item")
		return
	}

	var auctionItems []AuctionItem
	if err := db.Model(model.Auction{}).
		Preload("OwnedMinicon").
		Preload("OwnedMinicon.Minicon").
		Where("ID = ?", newAuction.ID).
		First(&auctionItems).Error; err != nil {
		log := utils.GetControllerLogger("/api/auction/sell [POST]")
		log.Errorln(err)
		helper.SendError(c, http.StatusInternalServerError, "Unable to fetch created auction item")
		return
	}

	res := GetAuctionBuyerResponse{
		AuctionItems: auctionItems,
	}

	helper.SendResponse(c, http.StatusCreated, res)
}
